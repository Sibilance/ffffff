#include "parser.h"

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

int init_parser_from_string(parser_t *parser, const unsigned char *input, size_t size)
{
    if (!yaml_parser_initialize(parser)) {
        return 0;
    }
    yaml_parser_set_input_string(parser, input, size);
    return 1;
}

int init_parser_from_file(parser_t *parser, FILE *file)
{
    if (!yaml_parser_initialize(parser)) {
        return 0;
    }
    yaml_parser_set_input_file(parser, file);
    return 1;
}

int init_parser_from_reader(parser_t *parser, read_handler_t *reader, void *data)
{
    if (!yaml_parser_initialize(parser)) {
        return 0;
    }
    yaml_parser_set_input(parser, reader, data);
    return 1;
}

int parser_parse(parser_t *parser, event_t *event)
{
    event->type = YAML_NO_EVENT;
    event->error = YAML_NO_ERROR;
    event->error_message = 0;
    event->error_context = 0;
    event->tag = 0;
    event->value = 0;
    event->quoted = false;

    yaml_event_t *_event = &event->_event;

    if (!yaml_parser_parse(parser, _event)) {
        event->line = parser->problem_mark.line + 1;
        event->column = parser->problem_mark.column + 1;
        event->error = parser->error;
        event->error_message = parser->problem;
        event->error_context = parser->context;
        event->tag = 0;
        return 0;
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

    return 1;
}

void parser_delete(parser_t *parser)
{
    yaml_parser_delete(parser);
}

void event_delete(event_t *event)
{
    yaml_event_delete(&event->_event);
}

const char *error_name(error_type_t error_type)
{
    return yaml_error_names[error_type];
}

const char *event_name(event_type_t event_type)
{
    return yaml_event_names[event_type];
}
