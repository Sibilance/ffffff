#include "error.h"

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

const char *yl_error_name(yl_error_type_t error_type)
{
    return yaml_error_names[error_type];
}
