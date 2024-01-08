#include "lauxlib.h"

#include "error.h"

const char *yaml_error_names[] = {
    "NO_ERROR",
    "MEMORY_ERROR",
    "READER_ERROR",
    "SCANNER_ERROR",
    "PARSER_ERROR",
    "COMPOSER_ERROR",
    "WRITER_ERROR",
    "EMITTER_ERROR",
    "EXECUTION_ERROR",
    "SYNTAX_ERROR",
    "RUNTIME_ERROR",
    "ERROR_HANDLER_ERROR",
    "TYPE_ERROR",
    "RENDER_ERROR",
    "ASSERTION_ERROR",
};

const char *yl_error_name(yl_error_type_t error_type)
{
    return yaml_error_names[error_type];
}

yl_error_type_t yl_error_from_lua_error(int status)
{
    switch (status) {
    case LUA_OK:
        return YL_NO_ERROR;
    case LUA_ERRRUN:
        return YL_RUNTIME_ERROR;
    case LUA_ERRSYNTAX:
        return YL_SYNTAX_ERROR;
    case LUA_ERRMEM:
        return YL_MEMORY_ERROR;
    case LUA_ERRERR:
        return YL_ERROR_HANDLER_ERROR;
    default:
        return YL_EXECUTION_ERROR;
    }
}

int yl_lua_error_handler(lua_State *L)
{
    const char *msg = lua_tostring(L, 1);
    luaL_traceback(L, L, msg, 1);
    return 1; // Just return the traceback, discard the message.
}
