#pragma once

#include <stdbool.h>

#include "lua.h"
#include "yaml.h"

#include "error.h"

const char *yl_event_name(yaml_event_type_t event_type);

typedef struct _yl_event_record_s {
    size_t capacity;
    size_t length;
    size_t index;
    yaml_event_t *events;
} yl_event_record_t;

int yl_copy_event(yaml_event_t *original, yaml_event_t *copy);

const char *yl_event_name(yaml_event_type_t event_type);

int yl_record_event(yl_event_record_t *event_record, yaml_event_t *event, yl_error_t *err);

void yl_event_record_delete(yl_event_record_t *event_record);

int yl_replay_event(yl_event_record_t *event_record, yaml_event_t *event, yl_error_t *err);

char *yl_render_event_record(yl_event_record_t *event_record, yl_error_t *err);

char *yl_copy_anchor(yaml_event_t *event);
