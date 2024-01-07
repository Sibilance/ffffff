#pragma once

#include "executor.h"

int yl_render_event(yl_event_consumer_t *consumer, yaml_event_t *event, lua_State *L, yl_error_t *err);

int yl_render_scalar(yl_event_consumer_t *consumer, yaml_event_t *event, lua_State *L, yl_error_t *err);

int yl_render_sequence(yl_event_consumer_t *consumer, yaml_event_t *event, lua_State *L, yl_error_t *err);
