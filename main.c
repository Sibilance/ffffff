#include <argp.h>
#include <stdio.h>

#include "lauxlib.h"
#include "lualib.h"
#include "yaml.h"

#include "executor.h"
#include "parser.h"
#include "test.h"

const char *argp_program_version = "yl 0.0.0";
const char *argp_program_bug_address = "https://github.com/Sibilance/ffffff/issues";
static char doc[] = "Render a YL template.";
static char args_doc[] = "[FILENAME]...";
static struct argp_option options[] = {
    {"in", 'i', "FILE", 0, "Input file to read from."},
    {"out", 'o', "FILE", 0, "Output file to write to."},
    {"debug", 'd', 0, 0, "Instead of generating YAML, output debug information"},
    {"test", 't', 0, 0, "Run the input as a test case. Test cases alternate input and expected output "
                        "documents in the stream. Test cases can be parameterized by preceeding them "
                        "with a document containing a sequence of mappings annotated with !testcases. "
                        "Keys from these mappings are set as global variables in each test case. The "
                        "number of expected output documents must equal the length of the !testcases "
                        "sequence."},
    {0}};

struct arguments {
    FILE *input, *output;
    bool debug;
    bool test;
};

static error_t parse_opt(int key, char *arg, struct argp_state *state)
{
    struct arguments *arguments = state->input;
    switch (key) {
    case 'i':
        if (strcmp(arg, "-") != 0) {
            arguments->input = fopen(arg, "rb");
        }
        break;
    case 'o':
        if (strcmp(arg, "-") != 0) {
            arguments->output = fopen(arg, "wb");
        }
        break;
    case 'd':
        arguments->debug = true;
        break;
    case 't':
        arguments->test = true;
        break;
    default:
        return ARGP_ERR_UNKNOWN;
    }
    return 0;
}

static struct argp argp = {options, parse_opt, args_doc, doc, 0, 0, 0};

int debug_handler(lua_State *L, yaml_event_t *event, yl_error_t *err)
{
    yaml_scalar_style_t style;
    fprintf(stderr, "%zu:%zu: %s\n", event->start_mark.line + 1, event->start_mark.column + 1, yl_event_name(event->type));
    if (lua_gettop(L)) {
        int type = lua_type(L, 1);
        switch (type) {
        case LUA_TNUMBER:
            if (lua_isinteger(L, 1))
                fprintf(stderr, "  LUA INTEGER: %lld\n", lua_tointeger(L, 1));
            else
                fprintf(stderr, "  LUA FLOAT: %#.17g\n", lua_tonumber(L, 1));
            break;
        case LUA_TBOOLEAN:
            fprintf(stderr, "  LUA BOOL: %s\n", lua_toboolean(L, 1) ? "true" : "false");
            break;
        case LUA_TSTRING:
            fprintf(stderr, "  LUA STRING: %s\n", lua_tostring(L, 1));
            break;
        case LUA_TTABLE:
            fprintf(stderr, "  LUA TABLE\n");
            break;
        case LUA_TNIL:
            fprintf(stderr, "  LUA NIL\n");
            break;
        default:
            fprintf(stderr, "  LUA UNEXPECTED TYPE: %s\n", lua_typename(L, type));
        }
        lua_settop(L, 0); // Clear the stack.
    }
    switch (event->type) {
    case YAML_SCALAR_EVENT:
        style = event->data.scalar.style;
        bool quoted = style == YAML_DOUBLE_QUOTED_SCALAR_STYLE || style == YAML_SINGLE_QUOTED_SCALAR_STYLE;
        fprintf(stderr, "  TAG: %s\n", event->data.scalar.tag);
        fprintf(stderr, "  QUOTED: %d\n  VALUE: %s\n",
                quoted,
                event->data.scalar.value);
        break;
    case YAML_SEQUENCE_START_EVENT:
        fprintf(stderr, "  TAG: %s\n", event->data.sequence_start.tag);
        break;
    case YAML_MAPPING_START_EVENT:
        fprintf(stderr, "  TAG: %s\n", event->data.mapping_start.tag);
        break;
    default:
        break;
    }
    return 1;
}

int emitter_handler(yaml_emitter_t *emitter, yaml_event_t *event, yl_error_t *err)
{
    if (!yaml_emitter_emit(emitter, event))
        goto error;

    // Mark the event as consumed to prevent double free.
    *event = (yaml_event_t){0};
    return 1;

error:
    // Mark the event as consumed to prevent double free.
    *event = (yaml_event_t){0};
    err->type = (yl_error_type_t)emitter->error;
    err->line = emitter->line;
    err->column = emitter->column;
    err->context = "While emitting YAML, encountered error";
    err->message = emitter->problem;
    return 0;
}

int main(int argc, char *argv[])
{
    struct arguments args = {
        stdin,
        stdout,
        false,
        false,
    };

    if (argp_parse(&argp, argc, argv, 0, 0, &args)) {
        return 1;
    }

    if (!args.input) {
        fprintf(stderr, "Error opening input file!\n");
        return 1;
    }
    if (!args.output) {
        fprintf(stderr, "Error opening output file!\n");
        return 1;
    }

    yl_execution_context_t ctx = {0};
    yaml_parser_t parser = {0};
    yaml_emitter_t emitter = {0};

    if (!yaml_parser_initialize(&parser)) {
        fprintf(stderr, "Error initializing parser!\n");
        goto error;
    }
    yaml_parser_set_input_file(&parser, args.input);

    ctx.producer = (yl_event_producer_t *)yl_parser_parse;
    ctx.producer_data = &parser;

    if (!yaml_emitter_initialize(&emitter)) {
        fprintf(stderr, "Error initializing emitter!\n");
        goto error;
    }
    yaml_emitter_set_unicode(&emitter, true);
    yaml_emitter_set_encoding(&emitter, YAML_UTF8_ENCODING);
    yaml_emitter_set_output_file(&emitter, args.output);

    ctx.lua = luaL_newstate();
    if (ctx.lua == NULL) {
        fprintf(stderr, "Error initializing lua!\n");
        goto error;
    }

    // Load only safe libraries.
    luaL_requiref(ctx.lua, LUA_GNAME, luaopen_base, true);
    luaL_requiref(ctx.lua, LUA_TABLIBNAME, luaopen_table, true);
    luaL_requiref(ctx.lua, LUA_STRLIBNAME, luaopen_string, true);
    luaL_requiref(ctx.lua, LUA_MATHLIBNAME, luaopen_math, true);
    luaL_requiref(ctx.lua, LUA_UTF8LIBNAME, luaopen_utf8, true);
    lua_settop(ctx.lua, 0);
    lua_pushglobaltable(ctx.lua);
    // Remove unsafe functions from base library.
    lua_pushnil(ctx.lua);
    lua_setfield(ctx.lua, 1, "dofile");
    lua_pushnil(ctx.lua);
    lua_setfield(ctx.lua, 1, "load");
    lua_pushnil(ctx.lua);
    lua_setfield(ctx.lua, 1, "loadfile");
    lua_pushnil(ctx.lua);
    lua_setfield(ctx.lua, 1, "require");
    lua_settop(ctx.lua, 0);

    if (args.debug) {
        ctx.handler = (yl_event_handler_t *)debug_handler;
        ctx.handler_data = ctx.lua;
    } else {
        ctx.handler = (yl_event_handler_t *)emitter_handler;
        ctx.handler_data = &emitter;
    }

    if (args.test) {
        if (!yl_test_stream(&ctx)) {
            fprintf(stderr, "Error testing stream!\n");
            fprintf(stderr, "%zu:%zu: %s: %s: %s\n",
                    ctx.err.line,
                    ctx.err.column,
                    yl_error_name(ctx.err.type),
                    ctx.err.context,
                    ctx.err.message);
            goto error;
        }
    } else if (!yl_execute_stream(&ctx)) {
        fprintf(stderr, "Error executing stream!\n");
        fprintf(stderr, "%zu:%zu: %s: %s: %s\n",
                ctx.err.line,
                ctx.err.column,
                yl_error_name(ctx.err.type),
                ctx.err.context,
                ctx.err.message);
        goto error;
    }

    yaml_parser_delete(&parser);
    yaml_emitter_delete(&emitter);
    lua_close(ctx.lua);

    return 0;

error:
    yaml_parser_delete(&parser);
    yaml_emitter_delete(&emitter);
    if (ctx.lua)
        lua_close(ctx.lua);

    return 1;
}
