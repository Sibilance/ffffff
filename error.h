#include <stddef.h>

#include "yaml.h"

typedef enum _yl_error_type_e {
    YL_NO_ERROR = YAML_NO_ERROR,
    YL_MEMORY_ERROR,
    YL_READER_ERROR,
    YL_SCANNER_ERROR,
    YL_PARSER_ERROR,
    YL_COMPOSER_ERROR,
    YL_WRITER_ERROR,
    YL_EMITTER_ERROR,
} yl_error_type_t;

typedef struct _yl_exec_error_s {
    size_t line, column;

} yl_exec_error;

const char *yl_error_name(yl_error_type_t error_type);
