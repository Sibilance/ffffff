#include "lauxlib.h"

#include "executor.h"

static int error_handler(lua_State *L)
{
    const char *msg = lua_tostring(L, 1);
    luaL_traceback(L, L, msg, 1);
    lua_remove(L, 1); // Remove original error message string.
    return 1;
}

static int execute_lua(lua_State *L, const char *buf, size_t len)
{
    int base = lua_gettop(L);
    lua_pushcfunction(L, error_handler);

    int status = luaL_loadbufferx(L, buf, len, buf, "t");
    if (status != LUA_OK) {
        lua_settop(L, base); // Restore the stack.
        return status;
    }

    status = lua_pcall(L, 0, 1, base + 1);

    lua_remove(L, base + 1); // Remove the error_handler.
    return status;
}

int yl_execute_stream(yl_execution_context_t *ctx)
{
    yl_event_t next_event = {0};

    bool done = false;
    while (!done) {
        if (!yl_parser_parse(&ctx->parser, &next_event, &ctx->err))
            return 0;

        if (!ctx->handler(ctx->data, &next_event, &ctx->err))
            goto error;

        switch (next_event.type) {
        case YAML_STREAM_START_EVENT:
            // Ignore the start of the stream, as there is no reason the caller
            // already would have consumed this before calling yl_execute_stream().
            break;
        case YAML_DOCUMENT_START_EVENT:
            if (!yl_execute_document(ctx, &next_event))
                goto error;
            break;
        case YAML_STREAM_END_EVENT:
            done = true;
            break;
        default:
            ctx->err.type = YL_EXECUTION_ERROR;
            ctx->err.line = next_event.line;
            ctx->err.column = next_event.column;
            ctx->err.context = "While executing a stream, got unexpected event";
            ctx->err.message = yl_event_name(next_event.type);
            goto error;
        }

        yl_event_delete(&next_event);
    }

    return 1;

error:
    yl_event_delete(&next_event);
    return 0;
}

int yl_execute_document(yl_execution_context_t *ctx, yl_event_t *event)
{
    yl_event_t next_event = {0};

    bool done = false;
    while (!done) {
        if (!yl_parser_parse(&ctx->parser, &next_event, &ctx->err))
            return 0;

        if (!ctx->handler(ctx->data, &next_event, &ctx->err))
            goto error;

        switch (next_event.type) {
        case YAML_SCALAR_EVENT:
            if (0)
                execute_lua(ctx->lua, "0", 1);
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
            done = true;
            break;
        default:
            ctx->err.type = YL_EXECUTION_ERROR;
            ctx->err.line = next_event.line;
            ctx->err.column = next_event.column;
            ctx->err.context = "While executing a document, got unexpected event";
            ctx->err.message = yl_event_name(next_event.type);
            goto error;
        }

        yl_event_delete(&next_event);
    }
    return 1;

error:
    yl_event_delete(&next_event);
    return 0;
}

int yl_execute_sequence(yl_execution_context_t *ctx, yl_event_t *event)
{
    yl_event_t next_event = {0};

    bool done = false;
    while (!done) {
        if (!yl_parser_parse(&ctx->parser, &next_event, &ctx->err))
            return 0;

        if (!ctx->handler(ctx->data, &next_event, &ctx->err))
            goto error;

        switch (next_event.type) {
        case YAML_SCALAR_EVENT:
            if (0)
                execute_lua(ctx->lua, "0", 1);
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
            done = true;
            break;
        default:
            ctx->err.type = YL_EXECUTION_ERROR;
            ctx->err.line = next_event.line;
            ctx->err.column = next_event.column;
            ctx->err.context = "While executing a sequence, got unexpected event";
            ctx->err.message = yl_event_name(next_event.type);
            goto error;
        }

        yl_event_delete(&next_event);
    }
    return 1;

error:
    yl_event_delete(&next_event);
    return 0;
}

int yl_execute_mapping(yl_execution_context_t *ctx, yl_event_t *event)
{
    yl_event_t next_event = {0};

    bool done = false;
    while (!done) {
        if (!yl_parser_parse(&ctx->parser, &next_event, &ctx->err))
            return 0;

        if (!ctx->handler(ctx->data, &next_event, &ctx->err))
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
            done = true;
            break;
        default:
            ctx->err.type = YL_EXECUTION_ERROR;
            ctx->err.line = next_event.line;
            ctx->err.column = next_event.column;
            ctx->err.context = "While executing a mapping, got unexpected event";
            ctx->err.message = yl_event_name(next_event.type);
            goto error;
        }

        yl_event_delete(&next_event);
    }
    return 1;

error:
    yl_event_delete(&next_event);
    return 0;
}

int yl_execute_scalar(yl_execution_context_t *ctx, yl_event_t *event)
{
    if (0)
        execute_lua(ctx->lua, "0", 1);

    return 1;
}
