#include <argp.h>
#include <stdio.h>

#include "lauxlib.h"
#include "lualib.h"
#include "yaml.h"

#include "executor.h"
#include "parser.h"

const char *argp_program_version = "yl 0.0.0";
const char *argp_program_bug_address = "<taliastocks@gmail.com>";
static char doc[] = "Render a YL template.";
static char args_doc[] = "[FILENAME]...";
static struct argp_option options[] = {
    {"in", 'i', "FILE", 0, "Input file to read from."},
    {"out", 'o', "FILE", 0, "Output file to write to."},
    {0}};

struct arguments {
    FILE *input, *output;
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
    default:
        return ARGP_ERR_UNKNOWN;
    }
    return 0;
}

static struct argp argp = {options, parse_opt, args_doc, doc, 0, 0, 0};

int handler(void *data, yl_event_t *event, yl_error_t *err)
{
    yaml_scalar_style_t style;
    fprintf(stderr, "%zu:%zu: %s\n", event->line, event->column, yl_event_name(event->event.type));
    switch (event->event.type) {
    case YAML_SCALAR_EVENT:
        style = event->event.data.scalar.style;
        bool quoted = style == YAML_DOUBLE_QUOTED_SCALAR_STYLE || style == YAML_SINGLE_QUOTED_SCALAR_STYLE;
        fprintf(stderr, "  TAG: %s\n", event->event.data.scalar.tag);
        fprintf(stderr, "  QUOTED: %d\n  VALUE: %s\n",
                quoted,
                event->event.data.scalar.value);
        break;
    case YAML_SEQUENCE_START_EVENT:
        fprintf(stderr, "  TAG: %s\n", event->event.data.sequence_start.tag);
        break;
    case YAML_MAPPING_START_EVENT:
        fprintf(stderr, "  TAG: %s\n", event->event.data.mapping_start.tag);
        break;
    default:
        break;
    }
    return 1;
}

int main(int argc, char *argv[])
{
    struct arguments args = {
        stdin,
        stdout,
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
    yaml_emitter_t emitter = {0};

    if (!yaml_parser_initialize(&ctx.parser)) {
        fprintf(stderr, "Error initializing parser!\n");
        goto error;
    }
    yaml_parser_set_input_file(&ctx.parser, args.input);

    if (!yaml_emitter_initialize(&emitter)) {
        fprintf(stderr, "Error initializing emitter!\n");
        goto error;
    }
    yaml_emitter_set_output_file(&emitter, args.output);

    ctx.lua = luaL_newstate();
    luaopen_base(ctx.lua);
    luaopen_table(ctx.lua);
    luaopen_string(ctx.lua);
    luaopen_utf8(ctx.lua);
    luaopen_math(ctx.lua);

    ctx.handler = handler;
    if (!yl_execute_stream(&ctx)) {
        fprintf(stderr, "Error executing stream!\n");
        fprintf(stderr, "%zu:%zu: %s: %s: %s\n",
                ctx.err.line,
                ctx.err.column,
                yl_error_name(ctx.err.type),
                ctx.err.context,
                ctx.err.message);
        goto error;
    }

    yaml_parser_delete(&ctx.parser);
    yaml_emitter_delete(&emitter);
    lua_close(ctx.lua);

    return 0;

error:
    yaml_parser_delete(&ctx.parser);
    yaml_emitter_delete(&emitter);
    if (ctx.lua)
        lua_close(ctx.lua);

    return 1;
}
