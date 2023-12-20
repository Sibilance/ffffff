#include <ctype.h>
#include <stdlib.h>
#include <string.h>

#include "lauxlib.h"

#include "executor.h"
#include "producer.h"

static int lua_error_handler(lua_State *L)
{
    const char *msg = lua_tostring(L, 1);
    luaL_traceback(L, L, msg, 1);
    return 1; // Just return the traceback, discard the message.
}

/*
** Execute a buffer in the Lua interpreter.
** Returns one of LUA_OK, LUA_ERRSYNTAX, LUA_ERRRUN, LUA_ERRMEM, or LUA_ERRERR.
** On return, leaves any return value or error message on the Lua stack.
*/
static int execute_lua(lua_State *L, const char *buf)
{
    if (!lua_checkstack(L, 3)) {
        return LUA_ERRMEM;
    }

    int base = lua_gettop(L);
    lua_pushcfunction(L, lua_error_handler);

    const char *retline = lua_pushfstring(L, "return %s;", buf);
    int status = luaL_loadbufferx(L, retline, strlen(retline), buf, "t");
    lua_remove(L, base + 2); // Remove retline.
    if (status == LUA_OK)
        status = lua_pcall(L, 0, 1, base + 1);

    lua_remove(L, base + 1); // Remove the error_handler.
    return status;
}

/*
** Execute a function in the Lua interpreter.
** Returns one of LUA_OK, LUA_ERRRUN, LUA_ERRMEM, or LUA_ERRERR.
** On return, leaves any return value or error message on the Lua stack.
*/
static int execute_lua_function(lua_State *L, const char *fnname, int nargs)
{
    int base = lua_gettop(L) - nargs;

    if (!lua_checkstack(L, 3)) {
        lua_settop(L, base);
        return LUA_ERRMEM;
    }

    int status = execute_lua(L, fnname);
    if (status != LUA_OK)
        return status;

    int type = lua_type(L, -1);

    if (type != LUA_TFUNCTION) {
        // Clear the stack and return an error message.
        lua_settop(L, base);
        lua_pushfstring(L, "expected `%s` to be a function, but instead got %s", fnname,
                        lua_typename(L, type));
        return LUA_ERRRUN;
    }

    lua_insert(L, base + 1); // Move the function below its argument(s).

    lua_pushcfunction(L, lua_error_handler);
    lua_insert(L, base + 1); // Move the error handler to the bottom.

    status = lua_pcall(L, nargs, 1, base + 1);

    lua_remove(L, base + 1); // Remove the error_handler.
    return status;
}

int yl_execute_stream(yl_execution_context_t *ctx)
{
    yaml_event_t next_event = {0};

    bool done = false;
    while (!done) {
        if (!yl_parser_parse(&ctx->parser, &next_event, &ctx->err))
            return 0;

        switch (next_event.type) {
        case YAML_STREAM_START_EVENT:
            if (!ctx->handler(ctx->data, &next_event, &ctx->err))
                goto error;
            break;
        case YAML_DOCUMENT_START_EVENT:
            if (!yl_execute_document(ctx, &next_event))
                goto error;
            break;
        case YAML_STREAM_END_EVENT:
            if (!ctx->handler(ctx->data, &next_event, &ctx->err))
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

    if (!ctx->handler(ctx->data, event, &ctx->err))
        goto error;

    bool done = false;
    while (!done) {
        if (!yl_parser_parse(&ctx->parser, &next_event, &ctx->err))
            return 0;

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
            if (!ctx->handler(ctx->data, &next_event, &ctx->err))
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
    return 1;

error:
    yaml_event_delete(&next_event);
    return 0;
}

int yl_execute_sequence(yl_execution_context_t *ctx, yaml_event_t *event)
{
    yaml_event_t next_event = {0};

    if (!ctx->handler(ctx->data, event, &ctx->err))
        goto error;

    bool done = false;
    while (!done) {
        if (!yl_parser_parse(&ctx->parser, &next_event, &ctx->err))
            return 0;

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
            if (!ctx->handler(ctx->data, &next_event, &ctx->err))
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
    return 1;

error:
    yaml_event_delete(&next_event);
    return 0;
}

int yl_execute_mapping(yl_execution_context_t *ctx, yaml_event_t *event)
{
    yaml_event_t next_event = {0};

    if (!ctx->handler(ctx->data, event, &ctx->err))
        goto error;

    bool done = false;
    while (!done) {
        if (!yl_parser_parse(&ctx->parser, &next_event, &ctx->err))
            return 0;

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
            if (!ctx->handler(ctx->data, &next_event, &ctx->err))
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
    return 1;

error:
    yaml_event_delete(&next_event);
    return 0;
}

int yl_execute_scalar(yl_execution_context_t *ctx, yaml_event_t *event)
{
    yaml_scalar_style_t style = event->data.scalar.style;
    if (!event->data.scalar.tag ||
        event->data.scalar.tag[0] != '!' ||
        event->data.scalar.tag[1] == '!') {

        free(event->data.scalar.tag);
        event->data.scalar.tag = NULL;
        event->data.scalar.plain_implicit = 1;
        event->data.scalar.quoted_implicit = 1;

        if (!ctx->handler(ctx->data, event, &ctx->err))
            goto error;

        return 1;
    }

    lua_settop(ctx->lua, 0); // Clear the stack.
    int status = LUA_OK;
    char *tag = (char *)event->data.scalar.tag;
    char *value = (char *)event->data.scalar.value;
    size_t length = event->data.scalar.length;
    if (strcmp(tag, "!") == 0) {
        if (style != YAML_DOUBLE_QUOTED_SCALAR_STYLE &&
            style != YAML_SINGLE_QUOTED_SCALAR_STYLE)
            status = execute_lua(ctx->lua, value);
        else
            lua_pushlstring(ctx->lua, value, length);
    } else {
        if (style == YAML_PLAIN_SCALAR_STYLE) {
            char *end;
            lua_Integer intvalue;
            lua_Number doublevalue;

            if (length == 0 || (length == 1 && value[0] == '~') || (length == 4 && strcmp(value, "null") == 0)) {
                lua_pushnil(ctx->lua);
            } else if (length == 4 && strcmp(value, "true") == 0) {
                lua_pushboolean(ctx->lua, true);
            } else if (length == 5 && strcmp(value, "false") == 0) {
                lua_pushboolean(ctx->lua, false);
            } else if (intvalue = strtoll(value, &end, 0), value + length == end) {
                lua_pushinteger(ctx->lua, intvalue);
            } else if (doublevalue = strtod(value, &end), value + length == end) {
                lua_pushnumber(ctx->lua, doublevalue);
            } else {
                lua_pushlstring(ctx->lua, value, length);
            }
        } else {
            lua_pushlstring(ctx->lua, value, length);
        }
        status = execute_lua_function(ctx->lua, (char *)event->data.scalar.tag + 1, 1);
    }

    if (status == LUA_OK) {
        if (!yl_produce_scalar(ctx, event))
            goto error;
    } else {
        ctx->err.type = YL_EXECUTION_ERROR;
        switch (status) {
        case LUA_ERRSYNTAX:
            ctx->err.type = YL_SYNTAX_ERROR;
            break;
        case LUA_ERRRUN:
            ctx->err.type = YL_RUNTIME_ERROR;
            break;
        case LUA_ERRMEM:
            ctx->err.type = YL_MEMORY_ERROR;
            break;
        case LUA_ERRERR:
            ctx->err.type = YL_ERROR_HANDLER_ERROR;
            break;
        default:
            ctx->err.type = YL_EXECUTION_ERROR;
            break;
        }
        ctx->err.line = event->start_mark.line;
        ctx->err.column = event->start_mark.column;
        ctx->err.context = "While executing a scalar, encountered an error";
        ctx->err.message = lua_tostring(ctx->lua, 1);
        goto error;
    }

    free(event->data.scalar.tag);
    event->data.scalar.tag = NULL;
    event->data.scalar.plain_implicit = 1;
    event->data.scalar.quoted_implicit = 1;

    if (!ctx->handler(ctx->data, event, &ctx->err))
        goto error;

    return 1;

error:
    return 0;
}
