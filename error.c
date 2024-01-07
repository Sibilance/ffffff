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

int yl_lua_error_handler(lua_State *L)
{
    const char *msg = lua_tostring(L, 1);
    luaL_traceback(L, L, msg, 1);
    return 1; // Just return the traceback, discard the message.
}
