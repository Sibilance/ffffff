#include "parser.h"

const char *yl_event_names[] = {
    "NO_EVENT",
    "STREAM_START_EVENT",
    "STREAM_END_EVENT",
    "DOCUMENT_START_EVENT",
    "DOCUMENT_END_EVENT",
    "ALIAS_EVENT",
    "SCALAR_EVENT",
    "SEQUENCE_START_EVENT",
    "SEQUENCE_END_EVENT",
    "MAPPING_START_EVENT",
    "MAPPING_END_EVENT",
};

int yl_parser_parse(yaml_parser_t *parser, yaml_event_t *event, yl_error_t *err)
{
    *event = (yaml_event_t){0};

    if (!yaml_parser_parse(parser, event)) {
        err->type = (yl_error_type_t)parser->error;
        err->line = parser->problem_mark.line;
        err->column = parser->problem_mark.column;
        err->context = parser->context;
        err->message = parser->problem;

        return 0;
    }

    return 1;
}

const char *yl_event_name(yaml_event_type_t event_type)
{
    return yl_event_names[event_type];
}
