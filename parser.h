#include <stdbool.h>
#include <stdio.h>

#include "yaml.h"

typedef yaml_parser_t parser_t;
typedef yaml_read_handler_t read_handler_t;
typedef yaml_error_type_t error_type_t;
typedef yaml_event_type_t event_type_t;

typedef struct _event_s {
    event_type_t type;

    size_t line, column;

    error_type_t error;
    const char *error_context;
    const char *error_message;

    const char *tag;
    const char *value;
    bool quoted;

    yaml_event_t _event;
} event_t;

/* Return 1 on success, 0 on error. */
int init_parser_from_string(parser_t *parser, const unsigned char *input, size_t size);
int init_parser_from_file(parser_t *parser, FILE *file);
int init_parser_from_reader(parser_t *parser, read_handler_t *reader, void *data);

int parser_parse(parser_t *parser, event_t *event);

void parser_delete(parser_t *parser);
void event_delete(event_t *event);

const char *error_name(error_type_t error_type);
const char *event_name(event_type_t event_type);
