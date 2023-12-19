#include <string.h>

#include "lauxlib.h"

#include "executor.h"

// -2^63 is 20 characters, plus NULL = 21.
// Also plenty for 17 digit precision floats.
#define NUMBUFSIZE 32

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
    lua_settop(L, 0); // Clear the stack.
    lua_pushcfunction(L, lua_error_handler);

    const char *retline = lua_pushfstring(L, "return (%s);", buf);
    int status = luaL_loadbufferx(L, retline, strlen(retline), buf, "t");
    lua_remove(L, 2); // Remove retline.
    if (status == LUA_OK)
        status = lua_pcall(L, 0, 1, 1);

    lua_remove(L, 1); // Remove the error_handler.
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
    if (style == YAML_DOUBLE_QUOTED_SCALAR_STYLE ||
        style == YAML_SINGLE_QUOTED_SCALAR_STYLE ||
        !event->data.scalar.tag ||
        strcmp((char *)event->data.scalar.tag, "!") != 0) {

        free(event->data.scalar.tag);
        event->data.scalar.tag = NULL;
        event->data.scalar.plain_implicit = 1;
        event->data.scalar.quoted_implicit = 1;

        if (!ctx->handler(ctx->data, event, &ctx->err))
            goto error;

        return 1;
    }

    int status = execute_lua(ctx->lua, (char *)event->data.scalar.value);

    if (status == LUA_OK) {
        int type = lua_type(ctx->lua, 1);
        switch (type) {
        case LUA_TNUMBER: {
            // -2^63 is 20 characters, plus NULL.
            // Also plenty for 17 digit precision floats.
            // We trim it down to the actual size after formatting.
            char *buf = malloc(NUMBUFSIZE);
            int len;
            if (lua_isinteger(ctx->lua, 1)) {
                len = snprintf(buf, NUMBUFSIZE, "%lld", lua_tointeger(ctx->lua, 1));
            } else {
                len = snprintf(buf, NUMBUFSIZE, "%.17g", lua_tonumber(ctx->lua, 1));
                if (strchr(buf, '.') == NULL && strchr(buf, 'e') == NULL) {
                    strcpy(buf + len, ".0");
                    len += 2;
                }
            }
            if (len < 0) {
                free(buf);
                ctx->err.type = YL_RUNTIME_ERROR;
                ctx->err.line = event->start_mark.line;
                ctx->err.column = event->start_mark.column;
                ctx->err.context = "While executing a scalar, got error formatting integer";
                ctx->err.message = "sprintf failed";
                goto error;
            }
            free(event->data.scalar.value);
            event->data.scalar.value = (yaml_char_t *)strndup(buf, len);
            event->data.scalar.length = len;
            free(buf);
            event->data.scalar.style = YAML_PLAIN_SCALAR_STYLE;
        } break;
        case LUA_TBOOLEAN:
            free(event->data.scalar.value);
            event->data.scalar.value = (yaml_char_t *)strdup(lua_toboolean(ctx->lua, 1) ? "true" : "false");
            event->data.scalar.length = strlen((char *)event->data.scalar.value);
            event->data.scalar.style = YAML_PLAIN_SCALAR_STYLE;
            break;
        case LUA_TSTRING:
            free(event->data.scalar.value);
            const char *lua_string = lua_tolstring(ctx->lua, 1, &event->data.scalar.length);
            event->data.scalar.value = (yaml_char_t *)strndup(lua_string, event->data.scalar.length);
            if (!strchr(lua_string, '\n'))
                event->data.scalar.style = YAML_ANY_SCALAR_STYLE;
            break;
        case LUA_TTABLE:
            break;
        case LUA_TNIL:
            free(event->data.scalar.value);
            event->data.scalar.value = (yaml_char_t *)strdup("~");
            event->data.scalar.length = 1;
            event->data.scalar.style = YAML_PLAIN_SCALAR_STYLE;
            break;
        default:
            ctx->err.type = YL_TYPE_ERROR;
            ctx->err.line = event->start_mark.line;
            ctx->err.column = event->start_mark.column;
            ctx->err.context = "While executing a scalar, got unexpected return type";
            ctx->err.message = lua_typename(ctx->lua, type);
            goto error;
        }
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
