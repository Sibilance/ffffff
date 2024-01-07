#include <ctype.h>

#include "lauxlib.h"
#include "lualib.h"

#include "event.h"
#include "render.h"

// -2^63 is 20 characters, plus NULL = 21.
// Also plenty for 17 digit precision floats.
#define NUMBUFSIZE 32

/*
** Get the length of the table at the top of the stack.
** Does not remove the table from the stack.
** On return, pushes either an error message, the length, or nil to the stack.
**/
static yl_error_type_t get_length(lua_State *L)
{
    if (lua_type(L, -1) != LUA_TTABLE) {
        lua_pushfstring(L, luaL_typename(L, -1));
        return YL_TYPE_ERROR;
    }

    if (lua_getmetatable(L, -1)) {
        // If there's a metatable, use it as the length operator.
        int type = lua_getfield(L, -1, "__len");
        lua_pop(L, 2); // Remove the function and metatable.
        if (type == LUA_TFUNCTION) {
            lua_len(L, -1); // Push the length onto the stack.

            int isnum;
            lua_tointegerx(L, -1, &isnum);
            if (!isnum) {
                lua_pushfstring(L, luaL_typename(L, -1));
                lua_remove(L, -2); // Remove the length.
                return YL_TYPE_ERROR;
            }

            return YL_NO_ERROR;
        }

        lua_pushnil(L); // No length, this is a mapping.
        return YL_NO_ERROR;
    } else {
        // If there's no metatable, get the "n" field if it exists,
        // or use the raw length operator.
        lua_getfield(L, -1, "n");
        if (lua_isinteger(L, -1)) {
            return YL_NO_ERROR;
        } else {
            lua_pop(L, 1); // Remove the nil "n" field value.
            int type = lua_rawgeti(L, -1, 1);
            // If there is a 1st element, it's a sequence. Get the length.
            if (type != LUA_TNIL) {
                lua_pop(L, 1); // Remove the value from rawgeti().
                lua_pushinteger(L, lua_rawlen(L, -1));
                return YL_NO_ERROR;
            }
            // -1: nil; -2: table
            // Assume a fully empty table is a sequence of zero length.
            if (lua_next(L, -2) == 0) {
                lua_pushinteger(L, 0);
                return YL_NO_ERROR;
            }
            // -1: value; -2: key; -3: table
            lua_pop(L, 2);
        }
    }

    lua_pushnil(L); // No length; this is a mapping.
    return YL_NO_ERROR;
}

int yl_render_event(yl_event_consumer_t *consumer, yaml_event_t *event, lua_State *L, yl_error_t *err)
{
    size_t line = event->start_mark.line;
    size_t column = event->start_mark.column;

    if (L != NULL) {
        int type = lua_type(L, -1);
        if (type == LUA_TTABLE) {
            yl_error_type_t errtype = get_length(L);
            if (errtype != YL_NO_ERROR) {
                err->type = errtype;
                err->line = line;
                err->column = column;
                err->context = "While rendering an event, got unexpected problem checking table length";
                err->message = lua_tostring(L, -1);
                goto error;
            }

            int isnum = lua_isinteger(L, -1);
            lua_pop(L, 1); // Pop the length.

            if (isnum) {
                if (!yl_render_sequence(consumer, event, L, err))
                    goto error;
            } else {
                if (!yl_render_mapping(consumer, event, L, err))
                    goto error;
            }
        } else {
            if (!yl_render_scalar(consumer, event, L, err))
                goto error;
        }
    } else if (!consumer->callback(consumer->data, event, NULL, err)) {
        goto error;
    }

    return 1;

error:
    yaml_event_delete(event);

    return 0;
}

int yl_render_scalar(yl_event_consumer_t *consumer, yaml_event_t *event, lua_State *L, yl_error_t *err)
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

    if (!consumer->callback(consumer->data, event, NULL, err))
        goto error;

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
    yaml_event_delete(event);
    lua_pop(L, 1); // Remove the argument from the stack.

    return 0;
}

int yl_render_sequence(yl_event_consumer_t *consumer, yaml_event_t *event, lua_State *L, yl_error_t *err)
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

    yl_error_type_t errtype = get_length(L);
    if (errtype != YL_NO_ERROR) {
        err->type = errtype;
        err->line = line;
        err->column = column;
        err->context = "While rendering a sequence, got invalid length";
        err->message = lua_tostring(L, -1);
        goto error;
    }
    int isnum;
    long int length = lua_tointegerx(L, -1, &isnum);
    if (!isnum) {
        err->type = YL_TYPE_ERROR;
        err->line = line;
        err->column = column;
        err->context = "While rendering a sequence, could not get length";
        err->message = luaL_typename(L, -1);
        lua_pop(L, 1);
        goto error;
    }
    lua_pop(L, 1);

    if (!yaml_sequence_start_event_initialize(event, (yaml_char_t *)anchor, NULL, 1, YAML_ANY_SEQUENCE_STYLE)) {
        err->type = YL_RENDER_ERROR;
        err->line = line;
        err->column = column;
        err->context = "While rendering a sequence, got unexpected error";
        err->message = "could not initialize sequence start event";
        goto error;
    }

    if (!consumer->callback(consumer->data, event, NULL, err))
        goto error;

    for (long int i = 1; i <= length; ++i) {
        lua_geti(L, -1, i);
        if (!yl_render_event(consumer, event, L, err))
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

    if (!consumer->callback(consumer->data, event, NULL, err))
        goto error;

    if (anchor != NULL)
        free(anchor);
    lua_pop(L, 1); // Remove the argument from the stack.

    return 1;

error:
    if (anchor != NULL)
        free(anchor);
    yaml_event_delete(event);
    lua_pop(L, 1); // Remove the argument from the stack.

    return 0;
}

int yl_render_mapping(yl_event_consumer_t *consumer, yaml_event_t *event, lua_State *L, yl_error_t *err)
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
        err->context = "While rendering a mapping, got unexpected Lua type";
        err->message = lua_typename(L, type);
        goto error;
    }

    if (!lua_checkstack(L, 10)) {
        err->type = YL_MEMORY_ERROR;
        err->line = line;
        err->column = column;
        err->context = "While rendering a mapping, got memory error";
        err->message = "could not expand Lua stack";
        goto error;
    }

    if (!yaml_mapping_start_event_initialize(event, (yaml_char_t *)anchor, NULL, 1, YAML_ANY_MAPPING_STYLE)) {
        err->type = YL_RENDER_ERROR;
        err->line = line;
        err->column = column;
        err->context = "While rendering a mapping, got unexpected error";
        err->message = "could not initialize mapping start event";
        goto error;
    }

    if (!consumer->callback(consumer->data, event, NULL, err))
        goto error;

    lua_newtable(L); // Add table for collecting keys.
    lua_insert(L, -2);
    lua_pushnil(L); // -1: nil; -2: table; -3: list of keys
    long int i = 1;
    while (lua_next(L, -2) != 0) {
        // -1: value; -2: key; -3: table; -4: list of keys
        lua_pop(L, 1);        // Pop value, we don't need it yet.
        lua_pushvalue(L, -1); // Duplicate key, it gets consumed by rawseti().
        // -1: key; -2: key; -3: table; -4: list of keys
        lua_rawseti(L, -4, i++);
        // -1: key; -2: table; -3: list of keys
    }
    // -1: table; -2: list of keys
    lua_pushcfunction(L, yl_lua_error_handler);
    luaL_requiref(L, LUA_TABLIBNAME, luaopen_table, false);
    type = lua_getfield(L, -1, "sort");
    lua_remove(L, -2); // Delete table module reference.
    // -1: table.sort; -2: error handler; -3; table; -4: list of keys
    lua_pushvalue(L, -4); // Duplicate list of keys, it gets consumed by sort().
    // -1: list of keys; -2: table.sort; -3: error handler; -4; table; -5: list of keys
    int status = lua_pcall(L, 1, 0, -3);
    if (status != LUA_OK) {
        // -1: error message; -2: error handler; -3: table; -4: list of keys
        lua_remove(L, -4); // Remove list of keys.
        lua_insert(L, -3); // Insert error message below table.
        lua_pop(L, 1);     // Pop error handler.
        // -1: table; -2: list of keys
        switch (status) {
        case LUA_ERRRUN:
            err->type = YL_RUNTIME_ERROR;
            break;
        case LUA_ERRMEM:
            err->type = YL_MEMORY_ERROR;
            break;
        case LUA_ERRERR:
            err->type = YL_ERROR_HANDLER_ERROR;
            break;
        default:
            err->type = YL_RUNTIME_ERROR;
            break;
        }
        err->type = YL_RUNTIME_ERROR;
        err->line = line;
        err->column = column;
        err->context = "While rendering a mapping, got error sorting keys";
        err->message = lua_tostring(L, -2);
        goto error;
    }
    // -1: error handler; -2 table; -3: list of keys (sorted)
    lua_pop(L, 1);     // Remove the error handler.
    lua_insert(L, -2); // Move table below list of keys.

    long int length = lua_rawlen(L, -1);
    for (long int i = 1; i <= length; ++i) {
        lua_rawgeti(L, -1, i);
        lua_pushvalue(L, -1); // Duplicate key, it gets consumed by yl_render_event().
        if (!yl_render_event(consumer, event, L, err)) {
            lua_pop(L, 2); // Pop the key and list of keys.
            goto error;
        }
        // -1: key; -2 list of keys (sorted); -3: table
        lua_gettable(L, -3); // Consumes key.
        if (!yl_render_event(consumer, event, L, err)) {
            lua_pop(L, 2); // Pop the value and list of keys.
            goto error;
        }
    }
    lua_pop(L, 1); // Pop the list of keys.

    if (!yaml_mapping_end_event_initialize(event)) {
        err->type = YL_RENDER_ERROR;
        err->line = line;
        err->column = column;
        err->context = "While rendering a mapping, got unexpected error";
        err->message = "could not initialize mapping end event";
        goto error;
    }

    if (!consumer->callback(consumer->data, event, NULL, err))
        goto error;

    if (anchor != NULL)
        free(anchor);
    lua_pop(L, 1); // Remove the argument from the stack.

    return 1;

error:
    if (anchor != NULL)
        free(anchor);
    yaml_event_delete(event);
    lua_pop(L, 1); // Remove the argument from the stack.

    return 0;
}
