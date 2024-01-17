#include <stdbool.h>
#include <string.h>

#include "lauxlib.h"
#include "lua.h"
#include "yaml.h"

#include "event.h"
#include "render.h"
#include "test.h"

int yl_test_stream(yl_execution_context_t *ctx)
{
    yaml_event_t next_event = {0};

    yl_event_record_t actual_events = {0};
    yl_event_record_t expected_events = {0};
    char *expected_rendering = NULL, *actual_rendering = NULL;

    if (!ctx->producer.callback(ctx->producer.data, &next_event, &ctx->err))
        goto error;

    if (next_event.type != YAML_STREAM_START_EVENT) {
        ctx->err.type = YL_PARSER_ERROR;
        ctx->err.line = next_event.start_mark.line;
        ctx->err.column = next_event.start_mark.column;
        ctx->err.context = "While trying to read a stream, got unexpected event";
        ctx->err.message = yl_event_name(next_event.type);
        goto error;
    }

    if (!ctx->consumer.callback(ctx->consumer.data, &next_event, NULL, &ctx->err))
        goto error;

    while (true) {
        if (!ctx->producer.callback(ctx->producer.data, &next_event, &ctx->err))
            goto error;

        size_t line = next_event.start_mark.line;
        size_t column = next_event.start_mark.column;

        if (!lua_checkstack(ctx->lua, 10)) {
            ctx->err.type = YL_MEMORY_ERROR;
            ctx->err.line = line;
            ctx->err.column = column;
            ctx->err.context = "While trying to run a testcase, encountered an error";
            ctx->err.message = "could not expand Lua stack space";
            goto error;
        }

        bool is_parameterized = false;
        lua_pushlightuserdata(ctx->lua, &is_parameterized);
        lua_pushnil(ctx->lua); // Second upvalue is nil, replaced by Lua table on first call.
        lua_pushcclosure(ctx->lua, yl_test_testcase, 2);
        lua_pushvalue(ctx->lua, -1); // Duplicate the closure so it is not consumed when setting global.
        lua_setglobal(ctx->lua, "testcases");
        // -1: yl_test_testcase closure

        switch (next_event.type) {
        case YAML_DOCUMENT_START_EVENT:
            if (!yl_test_record_document(ctx, &next_event, &actual_events, true))
                goto error;
            break;

        case YAML_STREAM_END_EVENT:
            if (!ctx->consumer.callback(ctx->consumer.data, &next_event, NULL, &ctx->err))
                goto error;
            goto done;

        default:
            ctx->err.type = YL_PARSER_ERROR;
            ctx->err.line = next_event.start_mark.line;
            ctx->err.column = next_event.start_mark.column;
            ctx->err.context = "While trying to read a document, got unexpected event";
            ctx->err.message = yl_event_name(next_event.type);
            goto error;
        }

        if (is_parameterized) {
            yl_event_record_delete(&actual_events);

            if (!ctx->producer.callback(ctx->producer.data, &next_event, &ctx->err))
                goto error;

            if (next_event.type != YAML_DOCUMENT_START_EVENT) {
                ctx->err.type = YL_PARSER_ERROR;
                ctx->err.line = next_event.start_mark.line;
                ctx->err.column = next_event.start_mark.column;
                ctx->err.context = "While trying to read a document, got unexpected event";
                ctx->err.message = yl_event_name(next_event.type);
                goto error;
            }

            if (!yl_test_record_document(ctx, &next_event, &actual_events, false))
                goto error;
            fprintf(stderr, "IT'S PARAMETERIZED!!!\n");
        }

        actual_rendering = yl_event_record_to_string(&actual_events, &ctx->err);
        if (actual_rendering == NULL)
            goto error;
        for (size_t i = 0; i < actual_events.length; ++i) {
            if (!ctx->consumer.callback(ctx->consumer.data, &actual_events.events[i], NULL, &ctx->err))
                goto error;
        }

        if (!ctx->producer.callback(ctx->producer.data, &next_event, &ctx->err))
            goto error;

        switch (next_event.type) {
        case YAML_DOCUMENT_START_EVENT:
            if (!yl_test_record_document(ctx, &next_event, &expected_events, true))
                goto error;
            break;

        default:
            ctx->err.type = YL_PARSER_ERROR;
            ctx->err.line = next_event.start_mark.line;
            ctx->err.column = next_event.start_mark.column;
            ctx->err.context = "While trying to read a document, got unexpected event";
            ctx->err.message = yl_event_name(next_event.type);
            goto error;
        }

        expected_rendering = yl_event_record_to_string(&expected_events, &ctx->err);
        if (expected_rendering == NULL)
            goto error;
        for (size_t i = 0; i < expected_events.length; ++i) {
            if (!ctx->consumer.callback(ctx->consumer.data, &expected_events.events[i], NULL, &ctx->err))
                goto error;
        }

        if (strcmp(actual_rendering, expected_rendering) != 0) {
            ctx->err.type = YL_ASSERTION_ERROR;
            ctx->err.line = line;
            ctx->err.column = column;
            ctx->err.context = "While comparing rendered documents";
            ctx->err.message = "actual document differs from expected document";
            goto error;
        }

        free(actual_rendering);
        actual_rendering = NULL;
        free(expected_rendering);
        expected_rendering = NULL;
        yl_event_record_delete(&actual_events);
        yl_event_record_delete(&expected_events);
    }
done:
    return 1;

error:
    if (actual_rendering != NULL)
        free(actual_rendering);
    if (expected_rendering != NULL)
        free(expected_rendering);
    yl_event_record_delete(&actual_events);
    yl_event_record_delete(&expected_events);
    yaml_event_delete(&next_event);
    return 0;
}

int yl_test_record_document(yl_execution_context_t *ctx, yaml_event_t *next_event, yl_event_record_t *event_record, bool execute)
{
    yl_event_consumer_t saved_consumer = ctx->consumer;

    ctx->consumer.callback = (yl_event_consumer_callback_t *)yl_record_event;
    ctx->consumer.data = event_record;

    if (execute) {
        if (!yl_execute_document(ctx, next_event))
            goto error;
    } else {
        if (!yl_test_passthrough_document(ctx, next_event))
            goto error;
    }

    ctx->consumer = saved_consumer;
    yaml_event_delete(next_event);

    return 1;

error:
    ctx->consumer = saved_consumer;
    yaml_event_delete(next_event);

    return 0;
}

int yl_test_passthrough_document(yl_execution_context_t *ctx, yaml_event_t *next_event)
{
    if (!ctx->consumer.callback(ctx->consumer.data, next_event, NULL, &ctx->err))
        goto error;

    bool done = false;
    while (!done) {
        if (!ctx->producer.callback(ctx->producer.data, next_event, &ctx->err))
            goto error;

        if (next_event->type == YAML_DOCUMENT_END_EVENT)
            done = true;

        if (!ctx->consumer.callback(ctx->consumer.data, next_event, NULL, &ctx->err))
            goto error;
    }

    return 1;

error:
    return 0;
}

int yl_test_testcase(lua_State *L)
{
    fprintf(stderr, "yl_test_testcase\n");
    luaL_checktype(L, lua_upvalueindex(1), LUA_TLIGHTUSERDATA);
    bool *is_parameterized = lua_touserdata(L, lua_upvalueindex(1));

    if (!*is_parameterized) {
        luaL_checktype(L, 1, LUA_TTABLE);

        *is_parameterized = true;
        lua_replace(L, lua_upvalueindex(2)); // Store table as upvalue.

        lua_pushnil(L);
        return 1;
    }
    return 0;
}
