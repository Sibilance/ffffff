#include <ctype.h>
#include <stdlib.h>
#include <string.h>

#include "lauxlib.h"

#include "error.h"
#include "executor.h"
#include "lua_helpers.h"
#include "render.h"

int yl_execute_stream(yl_execution_context_t *ctx)
{
    yaml_event_t next_event = {0};

    bool done = false;
    while (!done) {
        if (!ctx->producer.callback(ctx->producer.data, &next_event, &ctx->err))
            goto error;

        switch (next_event.type) {
        case YAML_STREAM_START_EVENT:
            if (!ctx->consumer.callback(ctx->consumer.data, &next_event, NULL, &ctx->err))
                goto error;
            break;
        case YAML_DOCUMENT_START_EVENT:
            if (!yl_execute_document(ctx, &next_event))
                goto error;
            break;
        case YAML_STREAM_END_EVENT:
            if (!ctx->consumer.callback(ctx->consumer.data, &next_event, NULL, &ctx->err))
                goto error;
            done = true;
            break;
        default:
            ctx->err.type = YL_EXECUTION_ERROR;
            ctx->err.line = next_event.start_mark.line;
            ctx->err.column = next_event.start_mark.column;
            ctx->err.context = "While executing a stream, got unexpected event";
            ctx->err.message = yl_event_name(next_event.type);
            goto error;
        }

        yaml_event_delete(&next_event);
    }

    return 1;

error:
    yaml_event_delete(&next_event);
    return 0;
}

int yl_execute_document(yl_execution_context_t *ctx, yaml_event_t *event)
{
    yaml_event_t next_event = {0};

    yl_event_consumer_t wrapped_consumer = ctx->consumer;
    ctx->consumer.callback = (yl_event_consumer_callback_t *)yl_render_event;
    ctx->consumer.data = &wrapped_consumer;

    if (!ctx->consumer.callback(ctx->consumer.data, event, NULL, &ctx->err))
        goto error;

    bool done = false;
    while (!done) {
        if (!ctx->producer.callback(ctx->producer.data, &next_event, &ctx->err))
            goto error;

        switch (next_event.type) {
        case YAML_SCALAR_EVENT:
            if (!yl_execute_scalar(ctx, &next_event))
                goto error;
            break;
        case YAML_SEQUENCE_START_EVENT:
            if (!yl_execute_sequence(ctx, &next_event))
                goto error;
            break;
        case YAML_MAPPING_START_EVENT:
            if (!yl_execute_mapping(ctx, &next_event))
                goto error;
            break;
        case YAML_DOCUMENT_END_EVENT:
            if (!ctx->consumer.callback(ctx->consumer.data, &next_event, NULL, &ctx->err))
                goto error;
            done = true;
            break;
        default:
            ctx->err.type = YL_EXECUTION_ERROR;
            ctx->err.line = next_event.start_mark.line;
            ctx->err.column = next_event.start_mark.column;
            ctx->err.context = "While executing a document, got unexpected event";
            ctx->err.message = yl_event_name(next_event.type);
            goto error;
        }

        yaml_event_delete(&next_event);
    }

    ctx->consumer = wrapped_consumer;
    return 1;

error:
    ctx->consumer = wrapped_consumer;
    yaml_event_delete(&next_event);
    return 0;
}

int yl_execute_sequence(yl_execution_context_t *ctx, yaml_event_t *event)
{
    yaml_event_t next_event = {0};

    size_t line = event->start_mark.line;
    size_t column = event->start_mark.column;
    yl_event_consumer_t saved_consumer = {0};
    yl_lua_table_builder_t table_builder = {0};
    char *tag = NULL;
    if (event->data.sequence_start.tag)
        if ((tag = strdup((char *)event->data.sequence_start.tag)) == NULL)
            goto memory_error;

    if (tag && tag[0] == '!' && tag[1] != '!') {
        saved_consumer = ctx->consumer;
        table_builder.L = ctx->lua;
        ctx->consumer.callback = (yl_event_consumer_callback_t *)yl_lua_table_builder;
        ctx->consumer.data = &table_builder;
    }

    // Ensure room for executing lua functions.
    if (!lua_checkstack(ctx->lua, 10))
        goto memory_error;

    if (!ctx->consumer.callback(ctx->consumer.data, event, NULL, &ctx->err))
        goto error;

    bool done = false;
    while (!done) {
        if (!ctx->producer.callback(ctx->producer.data, &next_event, &ctx->err))
            goto error;

        switch (next_event.type) {
        case YAML_SCALAR_EVENT:
            if (!yl_execute_scalar(ctx, &next_event))
                goto error;
            break;
        case YAML_SEQUENCE_START_EVENT:
            if (!yl_execute_sequence(ctx, &next_event))
                goto error;
            break;
        case YAML_MAPPING_START_EVENT:
            if (!yl_execute_mapping(ctx, &next_event))
                goto error;
            break;
        case YAML_SEQUENCE_END_EVENT:
            if (!ctx->consumer.callback(ctx->consumer.data, &next_event, NULL, &ctx->err))
                goto error;
            done = true;
            break;
        default:
            ctx->err.type = YL_EXECUTION_ERROR;
            ctx->err.line = next_event.start_mark.line;
            ctx->err.column = next_event.start_mark.column;
            ctx->err.context = "While executing a sequence, got unexpected event";
            ctx->err.message = yl_event_name(next_event.type);
            goto error;
        }

        yaml_event_delete(&next_event);
    }

    if (saved_consumer.callback != NULL)
        ctx->consumer = saved_consumer;

    if (tag != NULL) {
        int status = LUA_OK;
        if (strcmp(tag, "!") != 0)
            status = yl_lua_execute_lua_function(ctx->lua, tag + 1, 1);

        free(tag);
        tag = NULL;

        if (status != LUA_OK) {
            ctx->err.type = yl_error_from_lua_error(status);
            ctx->err.line = line;
            ctx->err.column = column;
            ctx->err.context = "While executing a sequence, encountered an error";
            ctx->err.message = lua_tostring(ctx->lua, 1);
            goto error;
        }

        if (!ctx->consumer.callback(ctx->consumer.data, event, ctx->lua, &ctx->err))
            goto error;
    }

    return 1;

memory_error:
    ctx->err.type = YL_MEMORY_ERROR;
    ctx->err.line = line;
    ctx->err.column = column;
    ctx->err.context = "While executing a mapping, encountered an error";
    ctx->err.message = "could not expand Lua stack space";
    goto error;

error:
    if (saved_consumer.callback != NULL)
        ctx->consumer = saved_consumer;
    if (tag != NULL) {
        free(tag);
        tag = NULL;
    }
    yaml_event_delete(&next_event);
    yaml_event_delete(event);
    return 0;
}

int yl_execute_mapping(yl_execution_context_t *ctx, yaml_event_t *event)
{
    yaml_event_t next_event = {0};

    size_t line = event->start_mark.line;
    size_t column = event->start_mark.column;
    yl_event_consumer_t saved_consumer = {0};
    yl_lua_table_builder_t table_builder = {0};
    char *tag = NULL;
    if (event->data.mapping_start.tag)
        if ((tag = strdup((char *)event->data.mapping_start.tag)) == NULL)
            goto memory_error;

    if (tag && tag[0] == '!' && tag[1] != '!') {
        saved_consumer = ctx->consumer;
        table_builder.L = ctx->lua;
        ctx->consumer.callback = (yl_event_consumer_callback_t *)yl_lua_table_builder;
        ctx->consumer.data = &table_builder;
    }

    // Ensure room for executing lua functions.
    if (!lua_checkstack(ctx->lua, 10))
        goto memory_error;

    if (!ctx->consumer.callback(ctx->consumer.data, event, NULL, &ctx->err))
        goto error;

    bool done = false;
    while (!done) {
        if (!ctx->producer.callback(ctx->producer.data, &next_event, &ctx->err))
            goto error;

        switch (next_event.type) {
        case YAML_SCALAR_EVENT:
            if (!yl_execute_scalar(ctx, &next_event))
                goto error;
            break;
        case YAML_SEQUENCE_START_EVENT:
            if (!yl_execute_sequence(ctx, &next_event))
                goto error;
            break;
        case YAML_MAPPING_START_EVENT:
            if (!yl_execute_mapping(ctx, &next_event))
                goto error;
            break;
        case YAML_MAPPING_END_EVENT:
            if (!ctx->consumer.callback(ctx->consumer.data, &next_event, NULL, &ctx->err))
                goto error;
            done = true;
            break;
        default:
            ctx->err.type = YL_EXECUTION_ERROR;
            ctx->err.line = next_event.start_mark.line;
            ctx->err.column = next_event.start_mark.column;
            ctx->err.context = "While executing a mapping, got unexpected event";
            ctx->err.message = yl_event_name(next_event.type);
            goto error;
        }

        yaml_event_delete(&next_event);
    }

    if (saved_consumer.callback != NULL)
        ctx->consumer = saved_consumer;

    if (tag != NULL) {
        int status = LUA_OK;
        if (strcmp(tag, "!") != 0)
            status = yl_lua_execute_lua_function(ctx->lua, tag + 1, 1);

        free(tag);
        tag = NULL;

        if (status != LUA_OK) {
            ctx->err.type = yl_error_from_lua_error(status);
            ctx->err.line = line;
            ctx->err.column = column;
            ctx->err.context = "While executing a mapping, encountered an error";
            ctx->err.message = lua_tostring(ctx->lua, 1);
            goto error;
        }

        if (!ctx->consumer.callback(ctx->consumer.data, event, ctx->lua, &ctx->err))
            goto error;
    }

    return 1;

memory_error:
    ctx->err.type = YL_MEMORY_ERROR;
    ctx->err.line = line;
    ctx->err.column = column;
    ctx->err.context = "While executing a mapping, encountered an error";
    ctx->err.message = "could not expand Lua stack space";
    goto error;

error:
    if (saved_consumer.callback != NULL)
        ctx->consumer = saved_consumer;
    if (tag != NULL) {
        free(tag);
        tag = NULL;
    }
    yaml_event_delete(&next_event);
    yaml_event_delete(event);
    return 0;
}

int yl_execute_scalar(yl_execution_context_t *ctx, yaml_event_t *event)
{
    int base = lua_gettop(ctx->lua);

    size_t line = event->start_mark.line;
    size_t column = event->start_mark.column;
    yaml_scalar_style_t style = event->data.scalar.style;
    char *tag = NULL;
    if (event->data.scalar.tag)
        tag = (char *)event->data.scalar.tag;

    if (!tag || tag[0] != '!' || tag[1] == '!') {
        if (!ctx->consumer.callback(ctx->consumer.data, event, NULL, &ctx->err))
            goto error;

        return 1;
    }

    // Ensure room for the scalar value, and executing lua functions.
    if (!lua_checkstack(ctx->lua, 10))
        goto memory_error;

    int status = LUA_OK;
    char *value = (char *)event->data.scalar.value;
    size_t length = event->data.scalar.length;
    if (strcmp(tag, "!") == 0) {
        if (style != YAML_DOUBLE_QUOTED_SCALAR_STYLE &&
            style != YAML_SINGLE_QUOTED_SCALAR_STYLE)
            status = yl_lua_execute_lua(ctx->lua, value);
        else
            lua_pushlstring(ctx->lua, value, length);
    } else {
        yl_lua_value_from_scalar(ctx->lua, style, length, value);
        status = yl_lua_execute_lua_function(ctx->lua, (char *)tag + 1, 1);
    }

    if (status != LUA_OK) {
        ctx->err.type = yl_error_from_lua_error(status);
        ctx->err.line = event->start_mark.line;
        ctx->err.column = event->start_mark.column;
        ctx->err.context = "While executing a scalar, encountered an error";
        ctx->err.message = lua_tostring(ctx->lua, 1);
        goto error;
    }

    if (!ctx->consumer.callback(ctx->consumer.data, event, ctx->lua, &ctx->err))
        goto error;

    lua_settop(ctx->lua, base);
    return 1;

memory_error:
    ctx->err.type = YL_MEMORY_ERROR;
    ctx->err.line = line;
    ctx->err.column = column;
    ctx->err.context = "While executing a scalar, encountered an error";
    ctx->err.message = "could not expand Lua stack space";
    goto error;

error:
    // Don't reset the stack on error, as it may contain an error message.
    yaml_event_delete(event);
    return 0;
}
