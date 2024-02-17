#pragma once

#include <stdbool.h>

#include "lauxlib.h"
#include "lua.h"
#include "lualib.h"
#include "yaml.h"

typedef struct _ylt_context_s {
    yaml_parser_t parser;
    yaml_emitter_t emitter;
    yaml_event_t event;
    lua_State *L;
} ylt_context_t;

void ylt_parser_error(ylt_context_t *ctx);
void ylt_emitter_error(ylt_context_t *ctx, const char *msg);
void ylt_event_error(ylt_context_t *ctx, const char *msg);

void ylt_try_parse(ylt_context_t *ctx);
void ylt_try_emit(ylt_context_t *ctx);

void ylt_evaluate_stream(ylt_context_t *ctx);
void ylt_evaluate_document(ylt_context_t *ctx);
void ylt_evaluate_sequence(ylt_context_t *ctx);
void ylt_evaluate_mapping(ylt_context_t *ctx);
void ylt_evaluate_scalar(ylt_context_t *ctx);
void ylt_evaluate_sequence_to_lua(ylt_context_t *ctx);
void ylt_evaluate_mapping_to_lua(ylt_context_t *ctx);
void ylt_evaluate_scalar_to_lua(ylt_context_t *ctx);
