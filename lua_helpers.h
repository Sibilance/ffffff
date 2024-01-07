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
 * Compare the top two values on the Lua stack, allowing values of different types
 * to be compared in a consistent way.
 *
 * @param[in,out]   L           A pointer to the Lua state.
 *
 * @returns On success, returns the number of results leaves the boolean result on
 * the stack.
 */
int yl_lua_compare(lua_State *L);

/**
 * Get the length of the table at the top of the stack without removing it.
 *
 * @param[in,out]   L       A pointer to the Lua state.
 * @param[in]       index   The stack index of the table.
 *
 * @returns On success, returns YL_NO_ERROR and leaves the length (or nil)
 * on the top of the stack. On error, returns the appropriate error and leaves
 * an error message on the stack.
 */
yl_error_type_t yl_lua_get_length(lua_State *L, int index);
