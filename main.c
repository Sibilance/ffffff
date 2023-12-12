#include <argp.h>
#include <stdio.h>

#include "parser.h"

const char *argp_program_version = "yl 0.0.0";
const char *argp_program_bug_address = "<taliastocks@gmail.com>";
static char doc[] = "Render a YL template.";
static char args_doc[] = "[FILENAME]...";
static struct argp_option options[] = {
    {"file", 'f', "FILE", 0, "Input file to read."},
    {0}};

struct arguments {
    char *filename;
};

static error_t parse_opt(int key, char *arg, struct argp_state *state)
{
    struct arguments *arguments = state->input;
    switch (key) {
    case 'f':
        arguments->filename = arg;
        break;
    default:
        return ARGP_ERR_UNKNOWN;
    }
    return 0;
}

static struct argp argp = {options, parse_opt, args_doc, doc, 0, 0, 0};

int main(int argc, char *argv[])
{
    struct arguments arguments;

    arguments.filename = "-";

    if (argp_parse(&argp, argc, argv, 0, 0, &arguments)) {
        return 1;
    }

    FILE *input;
    if (strcmp(arguments.filename, "-") == 0) {
        input = stdin;
    } else {
        input = fopen(arguments.filename, "rb");
    }

    if (!input) {
        printf("Error opening file!\n");
        return 1;
    }

    yl_parser_t parser;
    if (!yl_init_parser_from_file(&parser, input)) {
        printf("Error initializing parser!\n");
        return 1;
    }

    yl_event_t event;
    int done = 0;
    while (!done) {
        if (!yl_parser_parse(&parser, &event)) {
            printf("%zu:%zu: %s: %s: %s\n",
                   event.line,
                   event.column,
                   yl_error_name(event.error),
                   event.error_context,
                   event.error_message);
            break;
        }

        printf("%zu:%zu: %s\n", event.line, event.column, yl_event_name(event.type));
        printf("  TAG: %s\n", event.tag);
        switch (event.type) {
        case YAML_SCALAR_EVENT:
            printf("  QUOTED: %d\n  VALUE: %s\n",
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

    return 0;
}
