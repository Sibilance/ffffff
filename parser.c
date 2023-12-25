#include "parser.h"

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
