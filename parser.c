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

int yl_parser_parse(yaml_parser_t *parser, yl_event_t *event, yl_error_t *err)
{
    *event = (yl_event_t){0};

    yaml_event_t *_event = &event->event;

    if (!yaml_parser_parse(parser, _event)) {
        err->type = (yl_error_type_t)parser->error;
        err->line = parser->problem_mark.line + 1;
        err->column = parser->problem_mark.column + 1;
        err->context = parser->context;
        err->message = parser->problem;

        return 0;
    }

    event->type = _event->type;
    event->line = _event->start_mark.line + 1;
    event->column = _event->start_mark.column + 1;

    if (_event->type == YAML_SCALAR_EVENT) {
        yaml_scalar_style_t style = _event->data.scalar.style;
        event->quoted = style == YAML_DOUBLE_QUOTED_SCALAR_STYLE || style == YAML_SINGLE_QUOTED_SCALAR_STYLE;
    }

    return 1;
}

void yl_event_delete(yl_event_t *event)
{
    if (event->type) {
        yaml_event_delete(&event->event);
        *event = (yl_event_t){0};
    }
}

const char *yl_event_name(yaml_event_type_t event_type)
{
    return yl_event_names[event_type];
}
