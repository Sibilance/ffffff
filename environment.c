#include <stdbool.h>

#include "environment.h"

void yl_load_safe_libraries(lua_State *L)
{
    // Load only safe libraries.
    luaL_requiref(L, LUA_TABLIBNAME, luaopen_table, true);
    luaL_requiref(L, LUA_STRLIBNAME, luaopen_string, true);
    luaL_requiref(L, LUA_MATHLIBNAME, luaopen_math, true);
    luaL_requiref(L, LUA_UTF8LIBNAME, luaopen_utf8, true);
    lua_settop(L, 0);
    luaL_requiref(L, LUA_GNAME, luaopen_base, true);
    // Remove unsafe functions from base library.
    lua_pushnil(L);
    lua_setfield(L, 1, "dofile");
    lua_pushnil(L);
    lua_setfield(L, 1, "load");
    lua_pushnil(L);
    lua_setfield(L, 1, "loadfile");
    lua_pushnil(L);
    lua_setfield(L, 1, "require");
    lua_settop(L, 0);
}
