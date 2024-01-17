#pragma once

#include "executor.h"
#include "parser.h"

int yl_test_stream(yl_execution_context_t *ctx);

int yl_test_record_document(yl_execution_context_t *ctx, yaml_event_t *next_event, yl_event_record_t *event_record, bool execute);

int yl_test_passthrough_document(yl_execution_context_t *ctx, yaml_event_t *next_event);

int yl_test_testcase(lua_State *L);
