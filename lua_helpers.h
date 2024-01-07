#pragma once

#include "lua.h"

#include "error.h"

/**
 * Given a table at a given index, returns a sorted list of keys.
 *
 * @param[in,out]   L           A pointer to the Lua state.
 * @param[in]       index       The stack index of the table.
 *
 * @returns On success, returns YL_NO_ERROR and leaves the sorted array of keys
 * on the top of the stack. On error, returns the appropriate error and leaves
 * an error message on the stack.
 */
yl_error_type_t yl_lua_sort_keys(lua_State *L, int index);

/**
 * Given an array at a given index, sort it in-place.
 *
 * @param[in,out]   L           A pointer to the Lua state.
 * @param[in]       index       The stack index of the array.
 *
 * @returns On success, returns YL_NO_ERROR and leaves the stack unchanged. On
 * error, returns the appropriate error and leaves an error message no the stack.
 */
yl_error_type_t yl_lua_sort_array(lua_State *L, int index);

/**
 * Get the length of the table at the top of the stack without removing it.
 *
 * @param[in,out]   L           A pointer to the Lua state.
 * @param[in]       index       The stack index of the table.
 *
 * @returns On success, returns YL_NO_ERROR and leaves the length (or nil)
 * on the top of the stack. On error, returns the appropriate error and leaves
 * an error message on the stack.
 */
yl_error_type_t yl_lua_get_length(lua_State *L, int index);

/**
 * Execute a buffer in the Lua interpreter.
 *
 * @param[in,out]   L           A pointer to the Lua state.
 * @param[in]       buf         Lua code to execute.
 *
 * @returns One of LUA_OK, LUA_ERRSYNTAX, LUA_ERRRUN, LUA_ERRMEM, or LUA_ERRERR.
 * On return, leaves one return value or error message on the Lua stack.
 */
int yl_lua_execute_lua(lua_State *L, const char *buf);

/**
 * Execute a function (by name) in the Lua interpreter. If the function name appears
 * in the global table, it is pulled from there; otherwise, the "name" is executed
 * as Lua code and is expected to return a function.
 *
 * @param[in,out]   L           A pointer to the Lua state.
 * @param[in]       fnname      The "name" of the function to be executed.
 * @param[in]       nargs       The number of arguments on the stack to pass.
 *
 * @returns One of LUA_OK, LUA_ERRSYNTAX, LUA_ERRRUN, LUA_ERRMEM, or LUA_ERRERR.
 * On return, leaves one return value or error message on the Lua stack.
 */
int yl_lua_execute_lua_function(lua_State *L, const char *fnname, int nargs);
