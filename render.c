#include <ctype.h>

#include "lauxlib.h"

#include "event.h"
#include "render.h"

// -2^63 is 20 characters, plus NULL = 21.
// Also plenty for 17 digit precision floats.
#define NUMBUFSIZE 32

int yl_render_scalar(lua_State *L, yaml_event_t *event, yl_error_t *err)
{
    size_t line = event->start_mark.line;
    size_t column = event->start_mark.column;
    char *anchor = yl_copy_anchor(event);

    yaml_event_delete(event);

    char *buf = NULL;

    const char *value = NULL;
    size_t length = 0;
    yaml_scalar_style_t style = YAML_PLAIN_SCALAR_STYLE;

    int type = lua_type(L, -1);
    switch (type) {
    case LUA_TNUMBER: {
        buf = malloc(NUMBUFSIZE);
        int len;
        if (lua_isinteger(L, -1)) {
            len = snprintf(buf, NUMBUFSIZE, "%lld", lua_tointeger(L, -1));
        } else {
            len = snprintf(buf, NUMBUFSIZE, "%.17g", lua_tonumber(L, -1));
            if (strchr(buf, '.') == NULL && strchr(buf, 'e') == NULL) {
                strcpy(buf + len, ".0");
                len += 2;
            }
        }
        if (len < 0) {
            free(buf);
            err->type = YL_RUNTIME_ERROR;
            err->line = line;
            err->column = column;
            err->context = "While executing a scalar, got error formatting integer";
            err->message = "sprintf failed";
            goto error;
        }
        value = buf;
        length = len;
    } break;
    case LUA_TBOOLEAN:
        if (lua_toboolean(L, -1)) {
            value = "true";
            length = 4;
        } else {
            value = "false";
            length = 5;
        }
        break;
    case LUA_TSTRING:
        value = lua_tolstring(L, -1, &length);
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
        err->type = YL_TYPE_ERROR;
        err->line = line;
        err->column = column;
        err->context = "While rendering a scalar, got unexpected Lua type";
        err->message = lua_typename(L, type);
        goto error;
    }

    if (!yaml_scalar_event_initialize(event,
                                      (yaml_char_t *)anchor,
                                      NULL,
                                      (yaml_char_t *)value,
                                      length,
                                      1, 1,
                                      style)) {
        err->type = YL_RENDER_ERROR;
        err->line = line;
        err->column = column;
        err->context = "While rendering a scalar, got unexpected error";
        err->message = "could not initialize scalar event";
        goto error;
    }

    if (buf != NULL)
        free(buf);
    if (anchor != NULL)
        free(anchor);
    lua_pop(L, 1); // Remove the argument from the stack.

    return 1;

error:
    if (buf != NULL)
        free(buf);
    if (anchor != NULL)
        free(anchor);
    lua_pop(L, 1); // Remove the argument from the stack.

    return 0;
}

int yl_render_sequence(lua_State *L, yaml_event_t *event, yl_event_record_t *event_record, yl_error_t *err)
{
    size_t line = event->start_mark.line;
    size_t column = event->start_mark.column;
    char *anchor = yl_copy_anchor(event);

    yaml_event_delete(event);

    int type = lua_type(L, -1);
    if (type != LUA_TTABLE) {
        err->type = YL_TYPE_ERROR;
        err->line = line;
        err->column = column;
        err->context = "While rendering a sequence, got unexpected Lua type";
        err->message = lua_typename(L, type);
        goto error;
    }

    if (!lua_checkstack(L, 10)) {
        err->type = YL_MEMORY_ERROR;
        err->line = line;
        err->column = column;
        err->context = "While rendering a sequence, got memory error";
        err->message = "could not expand Lua stack";
        goto error;
    }

    long int length = -1;
    if (lua_getmetatable(L, -1)) {
        // If there's a metatable, use the ordinary length operator.
        lua_pop(L, 1); // Remove the metatable.
        lua_len(L, -1);
        int isnum;
        length = lua_tointegerx(L, -1, &isnum);
        if (!isnum) {
            err->type = YL_TYPE_ERROR;
            err->line = line;
            err->column = column;
            err->context = "While rendering a sequence, got invalid length type";
            err->message = luaL_typename(L, -1);
            lua_pop(L, 1); // Remove the length.
            goto error;
        }
        lua_pop(L, 1); // Remove the length.
    } else {
        // If there's no metatable, get the "n" field if it exists,
        // or use the raw length operator.
        lua_getfield(L, -1, "n");
        if (lua_isinteger(L, -1)) {
            length = lua_tointeger(L, -1);
        } else {
            length = lua_rawlen(L, -2);
        }
        lua_pop(L, 1); // Remove the "n" field value.
    }

    if (!yaml_sequence_start_event_initialize(event, (yaml_char_t *)anchor, NULL, 1, YAML_ANY_SEQUENCE_STYLE)) {
        err->type = YL_RENDER_ERROR;
        err->line = line;
        err->column = column;
        err->context = "While rendering a sequence, got unexpected error";
        err->message = "could not initialize sequence start event";
        goto error;
    }

    if (!yl_record_event(event_record, event, err))
        goto error;

    for (long int i = 1; i <= length; ++i) {
        int type = lua_geti(L, -1, i);
        (void)type;
        // TODO: infer the type and render nested sequences and mappings
        if (!yl_render_scalar(L, event, err)) {
            goto error;
        }

        if (!yl_record_event(event_record, event, err))
            goto error;
    }

    if (!yaml_sequence_end_event_initialize(event)) {
        err->type = YL_RENDER_ERROR;
        err->line = line;
        err->column = column;
        err->context = "While rendering a sequence, got unexpected error";
        err->message = "could not initialize sequence end event";
        goto error;
    }

    if (!yl_record_event(event_record, event, err))
        goto error;

    if (anchor != NULL)
        free(anchor);
    lua_pop(L, 1); // Remove the argument from the stack.

    return 1;

error:
    if (anchor != NULL)
        free(anchor);
    lua_pop(L, 1); // Remove the argument from the stack.

    return 0;
}
