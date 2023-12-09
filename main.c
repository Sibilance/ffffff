#include <argp.h>

#include "lua.h"
#include "yaml.h"

const char *yaml_error_names[] = {
    "NO_ERROR",
    "MEMORY_ERROR",
    "READER_ERROR",
    "SCANNER_ERROR",
    "PARSER_ERROR",
    "COMPOSER_ERROR",
    "WRITER_ERROR",
    "EMITTER_ERROR",
};

const char *argp_program_version = "yl 0.0.0";
const char *argp_program_bug_address = "<taliastocks@gmail.com>";
static char doc[] = "Render a YL template.";
static char args_doc[] = "[FILENAME]...";
static struct argp_option options[] = {
    {"file", 'f', "FILE", 0, "Input file to read."},
    {0}};

struct arguments
{
    char *filename;
};

static error_t parse_opt(int key, char *arg, struct argp_state *state)
{
    struct arguments *arguments = state->input;
    switch (key)
    {
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

    if (argp_parse(&argp, argc, argv, 0, 0, &arguments))
    {
        return 1;
    }

    FILE *input;
    if (strcmp(arguments.filename, "-") == 0)
    {
        input = stdin;
    }
    else
    {
        input = fopen(arguments.filename, "rb");
    }

    yaml_parser_t parser;
    yaml_parser_initialize(&parser);
    yaml_parser_set_input_file(&parser, input);

    int done = 0;
    while (!done)
    {
        yaml_event_t event;
        if (!yaml_parser_parse(&parser, &event))
        {
            printf("%zu:%zu: %s: %s",
                   parser.problem_mark.line,
                   parser.problem_mark.column,
                   yaml_error_names[parser.error],
                   parser.problem);
            break;
        }

        /*
          ...
          Process the event.
          ...
        */

        done = (event.type == YAML_STREAM_END_EVENT);

        yaml_event_delete(&event);
    }

    yaml_parser_delete(&parser);

    return 0;
}
