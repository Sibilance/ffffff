#pragma once

#include <stdbool.h>
#include <stddef.h>
#include <stdio.h>

#include "yaml.h"

#include "error.h"

typedef yaml_event_type_t yl_event_type_t;
typedef yaml_read_handler_t yl_read_handler_t;

typedef struct _yl_event_s {
    yl_event_type_t type;

    size_t line, column;

    const char *tag;
    const char *value;
    bool quoted;

    yaml_event_t _event;
} yl_event_t;

int yl_parser_parse(yaml_parser_t *parser, yl_event_t *event, yl_error_t *err);

void yl_event_delete(yl_event_t *event);

const char *yl_event_name(yaml_event_type_t event_type);
