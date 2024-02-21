#pragma once

#include <stdbool.h>

#include "lauxlib.h"
#include "lua.h"
#include "lualib.h"
#include "yaml.h"

#if defined(__GNUC__)
#define ylt_likely(x) (__builtin_expect(!!(x), 1))
#define ylt_unlikely(x) (__builtin_expect(!!(x), 0))
#else
#define ylt_likely(x) (!!(x))
#define ylt_unlikely(x) (!!(x))
#endif

const char *ylt_yaml_error_names[];
const char *ylt_yaml_event_names[];
const void *ylt_void_sentinel = &ylt_void_sentinel;

typedef enum _ylt_output_mode_e {
    YLT_EMITTER_OUTPUT_MODE,
    YLT_BUFFER_OUTPUT_MODE,
    YLT_LUA_OUTPUT_MODE,
} ylt_output_mode_t;

typedef struct _ylt_event_buffer_s {
    size_t cap;
    size_t len;
    yaml_event_t *events;
} ylt_event_buffer_t;

typedef struct _ylt_context_s {
    yaml_parser_t parser;
    yaml_emitter_t emitter;
    yaml_event_t event;
    ylt_event_buffer_t event_buffer;
    lua_State *L;
    ylt_output_mode_t output_mode;
} ylt_context_t;

void ylt_delete_context(ylt_context_t *ctx);

void ylt_parser_error(ylt_context_t *ctx);
void ylt_emitter_error(ylt_context_t *ctx, const char *msg);
void ylt_event_error(ylt_context_t *ctx, const char *msg);

void ylt_evaluate_stream(ylt_context_t *ctx);
void ylt_evaluate_document(ylt_context_t *ctx);
void ylt_evaluate_sequence(ylt_context_t *ctx);
void ylt_evaluate_mapping(ylt_context_t *ctx);
void ylt_evaluate_scalar(ylt_context_t *ctx);

void ylt_buffer_event(ylt_context_t *ctx);
void ylt_playback_event_buffer(ylt_context_t *ctx);
void ylt_truncate_event_buffer(ylt_context_t *ctx, size_t len);

void ylt_evaluate_lua(ylt_context_t *ctx);
void ylt_render_lua_value(ylt_context_t *ctx);

static inline void ylt_parse(ylt_context_t *ctx)
{
    if (ylt_unlikely(ctx->event.type != YAML_NO_EVENT))
        return ylt_event_error(ctx, "Unexpected non-empty event when parsing");
    if (ylt_unlikely(!yaml_parser_parse(&ctx->parser, &ctx->event)))
        return ylt_parser_error(ctx);
}

static inline void ylt_emit(ylt_context_t *ctx)
{
    // TODO: take into account output mode
    yaml_event_type_t event_type = ctx->event.type;
    if (ylt_unlikely(!yaml_emitter_emit(&ctx->emitter, &ctx->event)))
        return ylt_emitter_error(ctx, ylt_yaml_event_names[event_type]);
    // Event has been consumed (either deleted, or copied (by value) to an internal libyaml queue).
    // Mark it as empty.
    ctx->event = (yaml_event_t){0};
}

static inline bool ylt_is_lua_invocation(yaml_char_t *tag)
{
    return ylt_unlikely(tag) && ylt_likely(tag[0] == '!') && ylt_unlikely(tag[1] != '!');
}

static inline bool ylt_lua_value_is_void(ylt_context_t *ctx)
{
    return ylt_unlikely(lua_touserdata(ctx->L, -1) == ylt_void_sentinel);
}
