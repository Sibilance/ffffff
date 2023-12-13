#pragma once

#include <stddef.h>

#include "yaml.h"

#define YL_SUCCESS ((yl_error_t){0})

typedef enum _yl_error_type_e {
    YL_NO_ERROR = YAML_NO_ERROR,
    YL_MEMORY_ERROR = YAML_MEMORY_ERROR,
    YL_READER_ERROR = YAML_READER_ERROR,
    YL_SCANNER_ERROR = YAML_SCANNER_ERROR,
    YL_PARSER_ERROR = YAML_PARSER_ERROR,
    YL_COMPOSER_ERROR = YAML_COMPOSER_ERROR,
    YL_WRITER_ERROR = YAML_WRITER_ERROR,
    YL_EMITTER_ERROR = YAML_EMITTER_ERROR,
} yl_error_type_t;

typedef struct _yl_error_s {
    yl_error_type_t type;

    size_t line, column;

    const char *context;
    const char *message;
} yl_error_t;

const char *yl_error_name(yl_error_type_t error_type);
