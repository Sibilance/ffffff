#pragma once

#include "executor.h"

int yl_render_scalar(lua_State *L, yaml_event_t *event, yl_error_t *err);

int yl_render_sequence(lua_State *L, yaml_event_t *event, yl_event_record_t *event_record, yl_error_t *err);
