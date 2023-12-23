#include <ctype.h>

#include "producer.h"

// -2^63 is 20 characters, plus NULL = 21.
// Also plenty for 17 digit precision floats.
#define NUMBUFSIZE 32

int yl_produce_scalar(yl_execution_context_t *ctx, yaml_event_t *event)
{
    size_t line = event->start_mark.line;
    size_t column = event->start_mark.column;
    char *anchor = NULL;

    switch (event->type) {
    case YAML_SCALAR_EVENT:
        if (event->data.scalar.anchor != NULL)
            anchor = strdup((char *)event->data.scalar.anchor);
        break;
    case YAML_SEQUENCE_START_EVENT:
        if (event->data.sequence_start.anchor)
            anchor = strdup((char *)event->data.sequence_start.anchor);
        break;
    case YAML_MAPPING_START_EVENT:
        if (event->data.mapping_start.anchor)
            anchor = strdup((char *)event->data.mapping_start.anchor);
        break;
    default:
        break;
    }

    yaml_event_delete(event);

    char *buf = NULL;

    const char *value = NULL;
    size_t length = 0;
    yaml_scalar_style_t style = YAML_PLAIN_SCALAR_STYLE;

    int type = lua_type(ctx->lua, -1);
    switch (type) {
    case LUA_TNUMBER: {
        buf = malloc(NUMBUFSIZE);
        int len;
        if (lua_isinteger(ctx->lua, -1)) {
            len = snprintf(buf, NUMBUFSIZE, "%lld", lua_tointeger(ctx->lua, -1));
        } else {
            len = snprintf(buf, NUMBUFSIZE, "%.17g", lua_tonumber(ctx->lua, -1));
            if (strchr(buf, '.') == NULL && strchr(buf, 'e') == NULL) {
                strcpy(buf + len, ".0");
                len += 2;
            }
        }
        if (len < 0) {
            free(buf);
            ctx->err.type = YL_RUNTIME_ERROR;
            ctx->err.line = line;
            ctx->err.column = column;
            ctx->err.context = "While executing a scalar, got error formatting integer";
            ctx->err.message = "sprintf failed";
            goto error;
        }
        value = buf;
        length = len;
    } break;
    case LUA_TBOOLEAN:
        if (lua_toboolean(ctx->lua, -1)) {
            value = "true";
            length = 4;
        } else {
            value = "false";
            length = 5;
        }
        break;
    case LUA_TSTRING:
        value = lua_tolstring(ctx->lua, -1, &length);
        if (strchr(value, '\n'))
            style = YAML_LITERAL_SCALAR_STYLE;
        else if (length == 4 && strcmp(value, "true") == 0)
            style = YAML_DOUBLE_QUOTED_SCALAR_STYLE;
        else if (length == 5 && strcmp(value, "false") == 0)
            style = YAML_DOUBLE_QUOTED_SCALAR_STYLE;
        else if (length > 100)
            style = YAML_FOLDED_SCALAR_STYLE;
        else if (length > 0 && isdigit(value[0]))
            style = YAML_DOUBLE_QUOTED_SCALAR_STYLE;
        else if (length > 1 && value[0] == '.' && isdigit(value[1]))
            style = YAML_DOUBLE_QUOTED_SCALAR_STYLE;
        break;
    case LUA_TTABLE:
        // TODO, it's not a scalar, whoops
        break;
    case LUA_TNIL:
        value = "~";
        length = 1;
        break;
    default:
        ctx->err.type = YL_TYPE_ERROR;
        ctx->err.line = line;
        ctx->err.column = column;
        ctx->err.context = "While executing a scalar, got unexpected return type";
        ctx->err.message = lua_typename(ctx->lua, type);
        goto error;
    }

    if (!yaml_scalar_event_initialize(event,
                                      (yaml_char_t *)anchor,
                                      NULL,
                                      (yaml_char_t *)value,
                                      length,
                                      1, 1,
                                      style))
        goto error;

    if (buf != NULL)
        free(buf);
    if (anchor != NULL)
        free(anchor);

    return 1;

error:
    if (buf != NULL)
        free(buf);
    if (anchor != NULL)
        free(anchor);

    return 0;
}
