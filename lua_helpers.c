#include <stdbool.h>

#include "lauxlib.h"
#include "lualib.h"

#include "event.h"
#include "lua_helpers.h"

/**
 * Compare the top two values on the Lua stack, allowing values of different types
 * to be compared in a consistent way.
 *
 * @param[in,out]   L           A pointer to the Lua state.
 *
 * @returns On success, returns the number of results leaves the boolean result on
 * the stack.
 */
static int yl_lua_compare(lua_State *L)
{
    int left_type = lua_type(L, -2);
    int right_type = lua_type(L, -1);
    bool result;
    if (left_type == right_type) {
        result = lua_compare(L, -2, -1, LUA_OPLT);
    } else {
        result = left_type < right_type;
    }
    lua_pop(L, 2); // Remove arguments.
    lua_pushboolean(L, result);
    return 1;
}

yl_error_type_t yl_lua_sort_keys(lua_State *L, int index)
{
    index = lua_absindex(L, index);

    if (!lua_checkstack(L, 10)) {
        lua_pushfstring(L, "could not expand Lua stack");
        return YL_MEMORY_ERROR;
    }

    if (lua_type(L, index) != LUA_TTABLE) {
        lua_pushfstring(L, "Expected a table, instead got %s", luaL_typename(L, index));
        return YL_TYPE_ERROR;
    }

    lua_newtable(L); // Add table for collecting keys.
    lua_pushnil(L);  // -1: nil; -2: list of keys; index: table
    long int i = 1;
    while (lua_next(L, index) != 0) {
        // -1: value; -2: key; -3: list of keys; index: table
        lua_pop(L, 1);        // Discard value, we don't need it.
        lua_pushvalue(L, -1); // Duplicate key, it gets consumed by rawseti().
        // -1: key; -2: key; -3: list of keys; index: table
        lua_rawseti(L, -3, i++);
        // -1: key; -2: list of keys; index: table
    }
    // -1: list of keys; index: table

    return yl_lua_sort_array(L, -1);
}

yl_error_type_t yl_lua_sort_array(lua_State *L, int index)
{
    index = lua_absindex(L, index);

    if (!lua_checkstack(L, 10)) {
        lua_pushfstring(L, "could not expand Lua stack");
        return YL_MEMORY_ERROR;
    }

    if (lua_type(L, index) != LUA_TTABLE) {
        lua_pushfstring(L, "Expected an array table, instead got %s", luaL_typename(L, index));
        return YL_TYPE_ERROR;
    }

    lua_pushcfunction(L, yl_lua_error_handler);
    luaL_requiref(L, LUA_TABLIBNAME, luaopen_table, false);
    int type = lua_getfield(L, -1, "sort");
    lua_remove(L, -2); // Delete table module reference.
    // -1: table.sort; -2: error handler; index: array
    if (type != LUA_TFUNCTION) {
        lua_pop(L, 2); // Remove table.sort and error handler.
        lua_pushfstring(L, "Expected table.sort to be a function, instead got %s", lua_typename(L, type));
        return YL_TYPE_ERROR;
    }
    lua_pushvalue(L, index); // Pass array as argument to table.sort.
    lua_pushcfunction(L, yl_lua_compare);
    // -1: compare; -2 array; -3: table.sort; -4: error handler; index: array
    int status = lua_pcall(L, 2, 0, -4);
    if (status != LUA_OK) {
        // -1: error message; -2: error handler; index: array
        lua_remove(L, -2); // Remove error handler.
        switch (status) {
        case LUA_ERRRUN:
            return YL_RUNTIME_ERROR;
        case LUA_ERRMEM:
            return YL_MEMORY_ERROR;
        case LUA_ERRERR:
            return YL_ERROR_HANDLER_ERROR;
        default:
            return YL_RUNTIME_ERROR;
        }
    }
    // -1: error handler; index: array
    lua_pop(L, 1); // Remove the error handler.
    return YL_NO_ERROR;
}

yl_error_type_t yl_lua_get_length(lua_State *L, int index)
{
    index = lua_absindex(L, index);

    if (lua_type(L, index) != LUA_TTABLE) {
        lua_pushfstring(L, "Expected a table, instead got %s", luaL_typename(L, index));
        return YL_TYPE_ERROR;
    }

    if (lua_getmetatable(L, index)) {
        // If there's a metatable, use it as the length operator.
        int type = lua_getfield(L, -1, "__len");
        lua_pop(L, 2); // Remove the function and metatable.
        if (type == LUA_TFUNCTION) {
            lua_len(L, index); // Push the length onto the stack.

            if (!lua_isinteger(L, -1) && !lua_isnil(L, -1)) {
                lua_pushfstring(L, "Expected integer (or nil) length, instead got %s", luaL_typename(L, -1));
                lua_remove(L, -2); // Remove the length, leaving the error message.
                return YL_TYPE_ERROR;
            }

            return YL_NO_ERROR;
        }

        lua_pushnil(L); // No length, this is a mapping.
        return YL_NO_ERROR;
    } else {
        // If there's no metatable, get the "n" field if it exists,
        // or use the raw length operator.
        lua_getfield(L, index, "n");
        if (lua_isinteger(L, -1)) {
            return YL_NO_ERROR;
        } else {
            lua_pop(L, 1); // Remove the nil "n" field value.
            int type = lua_rawgeti(L, index, 1);
            // If there is a 1st element, it's a sequence. Get the length.
            if (type != LUA_TNIL) {
                lua_pop(L, 1); // Remove the value from rawgeti().
                lua_pushinteger(L, lua_rawlen(L, index));
                return YL_NO_ERROR;
            }
            // -1: nil; index: table
            // Assume a fully empty table is a sequence of zero length.
            if (lua_next(L, index) == 0) {
                // lua_next pops the key from the stack
                lua_pushinteger(L, 0);
                return YL_NO_ERROR;
            }
            // lua_next pushed a value and key to the stack
            // -1: value; -2: key; index: table
            lua_pop(L, 2);
        }
    }

    lua_pushnil(L); // No length; this is a mapping.
    return YL_NO_ERROR;
}

int yl_lua_execute_lua(lua_State *L, const char *buf)
{
    int base = lua_gettop(L);
    lua_pushcfunction(L, yl_lua_error_handler);

    const char *retline = lua_pushfstring(L, "return %s;", buf);
    int status = luaL_loadbufferx(L, retline, strlen(retline), buf, "t");
    lua_remove(L, base + 2); // Remove retline.
    if (status == LUA_OK)
        status = lua_pcall(L, 0, 1, base + 1);

    lua_remove(L, base + 1); // Remove the error_handler.
    return status;
}

int yl_lua_execute_lua_function(lua_State *L, const char *fnname, int nargs)
{
    int base = lua_gettop(L) - nargs;
    int status = LUA_OK;

    // First, try to get the value as a global variable.
    int type = lua_getglobal(L, fnname);

    if (type == LUA_TNIL) {
        lua_pop(L, 1);

        status = yl_lua_execute_lua(L, fnname);
        if (status != LUA_OK)
            return status;
    }

    type = lua_type(L, -1);

    if (type != LUA_TFUNCTION) {
        // Clear the stack and return an error message.
        lua_settop(L, base);
        lua_pushfstring(L, "expected `%s` to be a function, but instead got %s", fnname,
                        lua_typename(L, type));
        return LUA_ERRRUN;
    }

    lua_insert(L, base + 1); // Move the function below its argument(s).

    lua_pushcfunction(L, yl_lua_error_handler);
    lua_insert(L, base + 1); // Move the error handler to the bottom.

    status = lua_pcall(L, nargs, 1, base + 1);

    lua_remove(L, base + 1); // Remove the error_handler.
    return status;
}

static int new_tested_table_builder(yl_lua_table_builder_t *table_builder, yaml_event_t *event, yl_error_t *err)
{
    yl_lua_table_builder_t *parent = malloc(sizeof(yl_lua_table_builder_t));
    if (parent == NULL) {
        err->type = YL_MEMORY_ERROR;
        err->line = event->start_mark.line;
        err->column = event->start_mark.column;
        err->context = "While constructing a Lua table, got memory error";
        err->message = "unable to malloc new table builder";
        goto error;
    }
    *parent = *table_builder;
    table_builder->table_index = 0;
    table_builder->sequence_index = 0;
    table_builder->parent = parent;

    return 1;

error:
    return 0;
}

int yl_lua_table_builder(yl_lua_table_builder_t *table_builder, yaml_event_t *event, lua_State *L, yl_error_t *err)
{
    switch (event->type) {
    case YAML_MAPPING_START_EVENT: // Fall through.
    case YAML_SEQUENCE_START_EVENT:
        L = table_builder->L;

        // If we're already in a table, create a nested table builder.
        if (table_builder->table_index != 0)
            if (!new_tested_table_builder(table_builder, event, err))
                goto error;

        lua_newtable(L);
        table_builder->is_mapping = event->type == YAML_MAPPING_START_EVENT;
        table_builder->table_index = lua_gettop(L);
        break;
    case YAML_MAPPING_END_EVENT: // Fall through.
    case YAML_SEQUENCE_END_EVENT:
        // Place the newly constructed table at the top of the stack.
        L = table_builder->L;
        lua_settop(L, table_builder->table_index);

        if (table_builder->parent != NULL) {
            // If this was a nested table, back out to the parent table builder
            // and add this table as a member of the parent by falling through
            // to the "scalar" case.
            yl_lua_table_builder_t *parent = table_builder->parent;
            *table_builder = *parent;
            free(parent);
        } else {
            // If it's not a nested table, we're done.
            break;
        }
        // Fall through.
    case YAML_SCALAR_EVENT:
        ++table_builder->sequence_index;
        if (L == NULL) {
            // If L is NULL, this must be a YAML_SCALAR_EVENT. Convert it to a Lua value.
            L = table_builder->L;
            yl_lua_value_from_scalar(L, event->data.scalar.style, event->data.scalar.length, (char *)event->data.scalar.value);
        }
        if (table_builder->is_mapping) {
            if ((table_builder->sequence_index & 1) == 0) {
                // -1: value; -2: key; table_index: table
                lua_settable(L, table_builder->table_index);
            }
        } else {
            lua_seti(L, table_builder->table_index, table_builder->sequence_index);
        }
        break;
    default:
        err->type = YL_EXECUTION_ERROR;
        err->line = event->start_mark.line;
        err->column = event->start_mark.column;
        err->context = "While constructing a Lua table, got unexpected event";
        err->message = yl_event_name(event->type);
        goto error;
    }
    return 1;

error:
    return 0;
}

void yl_lua_value_from_scalar(lua_State *L, yaml_scalar_style_t style, size_t length, char *value)
{
    if (style == YAML_PLAIN_SCALAR_STYLE) {
        char *end;
        lua_Integer intvalue;
        lua_Number doublevalue;

        if (length == 0 || (length == 1 && value[0] == '~') || (length == 4 && strcmp(value, "null") == 0)) {
            lua_pushnil(L);
        } else if (length == 4 && strcmp(value, "true") == 0) {
            lua_pushboolean(L, true);
        } else if (length == 5 && strcmp(value, "false") == 0) {
            lua_pushboolean(L, false);
        } else if (intvalue = strtoll(value, &end, 0), value + length == end) {
            lua_pushinteger(L, intvalue);
        } else if (doublevalue = strtod(value, &end), value + length == end) {
            lua_pushnumber(L, doublevalue);
        } else {
            lua_pushlstring(L, value, length);
        }
    } else {
        lua_pushlstring(L, value, length);
    }
}
