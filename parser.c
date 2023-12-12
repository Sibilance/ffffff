#include "parser.h"

const char *yaml_event_names[] = {
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

int yl_init_parser_from_string(yl_parser_t *parser, const unsigned char *input, size_t size)
{
    if (!yaml_parser_initialize(parser)) {
        return 0;
    }
    yaml_parser_set_input_string(parser, input, size);
    return 1;
}

int yl_init_parser_from_file(yl_parser_t *parser, FILE *file)
{
    if (!yaml_parser_initialize(parser)) {
        return 0;
    }
    yaml_parser_set_input_file(parser, file);
    return 1;
}

int yl_init_parser_from_reader(yl_parser_t *parser, yl_read_handler_t *reader, void *data)
{
    if (!yaml_parser_initialize(parser)) {
        return 0;
    }
    yaml_parser_set_input(parser, reader, data);
    return 1;
}

yl_error_t yl_parser_parse(yl_parser_t *parser, yl_event_t *event)
{
    *event = (yl_event_t){0};

    yaml_event_t *_event = &event->_event;

    if (!yaml_parser_parse(parser, _event)) {
        return (yl_error_t){
            (yl_error_type_t)parser->error,
            parser->problem_mark.line + 1,
            parser->problem_mark.column + 1,
            parser->context,
            parser->problem,
        };
    }

    event->type = _event->type;
    event->line = _event->start_mark.line + 1;
    event->column = _event->start_mark.column + 1;

    switch (_event->type) {
    case YAML_SCALAR_EVENT:
        event->tag = (const char *)_event->data.scalar.tag;
        event->value = (const char *)_event->data.scalar.value;
        yaml_scalar_style_t style = _event->data.scalar.style;
        event->quoted = style == YAML_DOUBLE_QUOTED_SCALAR_STYLE || style == YAML_SINGLE_QUOTED_SCALAR_STYLE;
        break;
    case YAML_SEQUENCE_START_EVENT:
        event->tag = (const char *)_event->data.sequence_start.tag;
        break;
    case YAML_MAPPING_START_EVENT:
        event->tag = (const char *)_event->data.mapping_start.tag;
        break;
    default:
        break;
    }

    return (yl_error_t){0};
}

void yl_parser_delete(yl_parser_t *parser)
{
    yaml_parser_delete(parser);
}

void yl_event_delete(yl_event_t *event)
{
    if (event->type) {
        yaml_event_delete(&event->_event);
        *event = (yl_event_t){0};
    }
}

const char *yl_event_name(yl_event_type_t event_type)
{
    return yaml_event_names[event_type];
}
