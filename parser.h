#include <stdbool.h>
#include <stddef.h>
#include <stdio.h>

#include "yaml.h"

#include "error.h"

typedef yaml_parser_t yl_parser_t;
typedef yaml_read_handler_t yl_read_handler_t;
typedef yaml_event_type_t yl_event_type_t;

typedef struct _yl_event_s {
    yl_event_type_t type;

    size_t line, column;

    const char *tag;
    const char *value;
    bool quoted;

    yaml_event_t _event;
} yl_event_t;

/* Return 1 on success, 0 on error. */
int yl_init_parser_from_string(yl_parser_t *parser, const unsigned char *input, size_t size);
int yl_init_parser_from_file(yl_parser_t *parser, FILE *file);
int yl_init_parser_from_reader(yl_parser_t *parser, yl_read_handler_t *reader, void *data);

yl_error_t yl_parser_parse(yl_parser_t *parser, yl_event_t *event);

void yl_parser_delete(yl_parser_t *parser);
void yl_event_delete(yl_event_t *event);

const char *yl_event_name(yl_event_type_t event_type);
