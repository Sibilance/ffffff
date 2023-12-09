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
    if (!yaml_parser_parse(parser, &event->_event)) {
        event->type = YAML_NO_EVENT;
        event->line = parser->problem_mark.line + 1;
        event->column = parser->problem_mark.column + 1;
        event->error = parser->error;
        event->error_message = parser->problem;
        event->error_context = parser->context;
        return 0;
    }

    yaml_event_t *_event = &event->_event;
    event->type = _event->type;
    event->line = _event->start_mark.line + 1;
    event->column = _event->start_mark.column + 1;
    event->error = YAML_NO_ERROR;
    event->error_message = 0;
    event->error_context = 0;
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
