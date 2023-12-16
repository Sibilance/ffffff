#pragma once

#include <stdbool.h>
#include <stddef.h>
#include <stdio.h>

#include "yaml.h"

#include "error.h"

int yl_parser_parse(yaml_parser_t *parser, yaml_event_t *event, yl_error_t *err);

const char *yl_event_name(yaml_event_type_t event_type);
