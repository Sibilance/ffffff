#include <stdbool.h>
#include <stdio.h>

#include "yaml.h"

typedef yaml_parser_t parser_t;
typedef yaml_read_handler_t read_handler_t;
typedef yaml_error_type_t error_type_t;
typedef yaml_event_type_t event_type_t;

typedef struct _yl_event_s {
    event_type_t type;

    size_t line, column;

    error_type_t error;
    const char *error_context;
    const char *error_message;

    const char *tag;
    const char *value;
    bool quoted;

    yaml_event_t _event;
} yl_event_t;

/* Return 1 on success, 0 on error. */
int yl_init_parser_from_string(parser_t *parser, const unsigned char *input, size_t size);
int yl_init_parser_from_file(parser_t *parser, FILE *file);
int yl_init_parser_from_reader(parser_t *parser, read_handler_t *reader, void *data);

int yl_parser_parse(parser_t *parser, yl_event_t *event);

void yl_parser_delete(parser_t *parser);
void yl_event_delete(yl_event_t *event);

const char *yl_error_name(error_type_t error_type);
const char *yl_event_name(event_type_t event_type);
