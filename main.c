#include <argp.h>
#include <stdio.h>

#include "yaml.h"

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

    yaml_parser_t parser;
    if (!yaml_parser_initialize(&parser)) {
        fprintf(stderr, "Error initializing parser!\n");
        return 1;
    }
    yaml_parser_set_input_file(&parser, args.input);

    yaml_emitter_t emitter;
    if (!yaml_emitter_initialize(&emitter)) {
        fprintf(stderr, "Error initializing emitter!\n");
    }
    yaml_emitter_set_output_file(&emitter, args.output);

    yl_event_t event;
    int done = 0;
    while (!done) {
        yl_error_t err = yl_parser_parse(&parser, &event);
        if (err.type) {
            fprintf(stderr, "%zu:%zu: %s: %s: %s\n",
                    err.line,
                    err.column,
                    yl_error_name(err.type),
                    err.error_context,
                    err.error_message);
            break;
        }

        fprintf(stderr, "%zu:%zu: %s\n", event.line, event.column, yl_event_name(event.type));
        fprintf(stderr, "  TAG: %s\n", event.tag);
        switch (event.type) {
        case YAML_SCALAR_EVENT:
            fprintf(stderr, "  QUOTED: %d\n  VALUE: %s\n",
                    event.quoted,
                    event.value);
            break;
        default:
            break;
        }

        done = (event.type == YAML_STREAM_END_EVENT);

        yl_event_delete(&event);
    }

    yaml_parser_delete(&parser);
    yaml_emitter_delete(&emitter);

    return 0;
}
