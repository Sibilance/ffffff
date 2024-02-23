#include "ylt.h"

const char *ylt_yaml_error_names[] = {
    "NO_ERROR",
    "MEMORY_ERROR",
    "READER_ERROR",
    "SCANNER_ERROR",
    "PARSER_ERROR",
    "COMPOSER_ERROR",
    "WRITER_ERROR",
    "EMITTER_ERROR",
};

const char *ylt_yaml_event_names[] = {
    "NO_EVENT",
    "STREAM_START_EVENT",
    "STREAM_END_EVENT",
    "DOCUMENT_START_EVENT",
    "DOCUMENT_END_EVENT",
    "ALIAS_EVENT",
    "SCALAR_EVENT",
    "SEQUENCE_START_EVENT",
    "SEQUENCE_END_EVENT",
    "MAPPING_START_EVENT",
    "MAPPING_END_EVENT",
};


void ylt_delete_context(ylt_context_t *ctx)
{
    yaml_parser_delete(&ctx->parser);
    yaml_emitter_delete(&ctx->emitter);
    yaml_event_delete(&ctx->event);
    for (size_t i = 0; i < ctx->event_buffer.len; ++i) {
        yaml_event_delete(&ctx->event_buffer.events[i]);
    }
    free(ctx->event_buffer.events);
    lua_close(ctx->L);
    *ctx = (ylt_context_t){0};
}


void ylt_parser_error(ylt_context_t *ctx)
{
    luaL_error(ctx->L, "%I:%I: %s: %s: %s",
               (lua_Integer)ctx->parser.problem_mark.line + 1,
               (lua_Integer)ctx->parser.problem_mark.column + 1,
               ylt_yaml_error_names[ctx->parser.error],
               ctx->parser.context,
               ctx->parser.problem);
}


void ylt_emitter_error(ylt_context_t *ctx, const char *msg)
{
    luaL_error(ctx->L, "%I:%I: %s: %s: %s",
               (lua_Integer)ctx->parser.context_mark.line + 1, // Use the parser position; the emitter position is not helpful.
               (lua_Integer)ctx->parser.context_mark.column + 1,
               ylt_yaml_error_names[ctx->emitter.error],
               msg,
               ctx->emitter.problem);
}


void ylt_event_error(ylt_context_t *ctx, const char *msg)
{
    lua_Integer line, column;
    if (ctx->event.type != YAML_NO_EVENT) {
        line = ctx->event.start_mark.line + 1;
        column = ctx->event.start_mark.column + 1;
    } else {
        line = ctx->parser.context_mark.line + 1;
        column = ctx->parser.context_mark.column + 1;
    }
    luaL_error(ctx->L, "%I:%I: EVENT_ERROR: %s: %s", line, column, msg, ylt_yaml_event_names[ctx->event.type]);
}


void ylt_evaluate_stream(ylt_context_t *ctx)
{
    if (ylt_unlikely(ctx->event.type != YAML_NO_EVENT))
        return ylt_event_error(ctx, "Unexpected event already parsed when evaluating stream");

    ylt_parse(ctx);

    if (ylt_unlikely(ctx->event.type != YAML_STREAM_START_EVENT))
        return ylt_event_error(ctx, "Unexpected event at start of stream (expecting STREAM_START_EVENT)");

    ylt_emit(ctx);

    for (;;) {
        ylt_parse(ctx);

        switch (ctx->event.type) {
        case YAML_DOCUMENT_START_EVENT:
            ylt_evaluate_document(ctx);
            break;
        case YAML_STREAM_END_EVENT:
            ylt_emit(ctx);
            return;
        default:
            return ylt_event_error(ctx, "Unexpected event while processing stream");
        }
    }
}


static inline void ylt_evaluate_nested(ylt_context_t *ctx, char *processing_what)
{
    switch (ctx->event.type) {
    case YAML_SEQUENCE_START_EVENT:
        ylt_evaluate_sequence(ctx);
        break;
    case YAML_MAPPING_START_EVENT:
        ylt_evaluate_mapping(ctx);
        break;
    case YAML_SCALAR_EVENT:
        ylt_evaluate_scalar(ctx);
        break;
    case YAML_ALIAS_EVENT:
        ylt_emit(ctx);
        break;
    default:
        ylt_event_error(ctx, lua_pushfstring(ctx->L, "Unexpected event while processing %s", processing_what));
    }
}


void ylt_evaluate_document(ylt_context_t *ctx)
{
    if (ylt_unlikely(ctx->event.type != YAML_DOCUMENT_START_EVENT))
        return ylt_event_error(ctx, "Unexpected event at start of document");

    ylt_output_mode_t initial_output_mode = ctx->output_mode;
    size_t initial_buffer_len = ctx->event_buffer.len;

    // Buffer the document start event in case the document content indicates it should be skipped.
    ylt_buffer_event(ctx);

    ylt_parse(ctx);

    // If the next event has no Lua invocation tag, the content will be passed-through. Flush the buffer (output the
    // DOCUMENT START EVENT) and continue processing.  Otherwise, change to LUA OUTPUT MODE and continue processing.

    bool is_lua_invocation = ylt_is_lua_invocation(ctx);
    if (is_lua_invocation)
        ctx->output_mode = YLT_LUA_OUTPUT_MODE;
    else
        ylt_playback_event_buffer(ctx, initial_buffer_len);

    ylt_evaluate_nested(ctx, "document");

    // If this was a Lua invocation, execute the Lua and check the output. If the output is VOID, discard the
    // entire document. Otherwise, output the entire document starting with the buffered DOCUMENT START EVENT.
    if (is_lua_invocation) {
        ylt_execute_lua(ctx);
        ctx->output_mode = ylt_unlikely(ylt_lua_value_is_void(ctx)) ? YLT_DISCARD_OUTPUT_MODE : initial_output_mode;
        ylt_playback_event_buffer(ctx, initial_buffer_len);
        ylt_render_lua_value(ctx);
    }

    ylt_parse(ctx); // Expect DOCUMENT END EVENT.
    if (ylt_unlikely(ctx->event.type != YAML_DOCUMENT_END_EVENT))
        return ylt_event_error(ctx, "Unexpected event at end of document");
    ylt_emit(ctx); // Output DOCUMENT END EVENT.

    // Finally, restore the original output mode.
    ctx->output_mode = initial_output_mode;
}


void ylt_evaluate_sequence(ylt_context_t *ctx)
{
    if (ylt_unlikely(ctx->event.type != YAML_SEQUENCE_START_EVENT))
        return ylt_event_error(ctx, "Unexpected event at start of sequence");

    ylt_emit(ctx);

    for (ylt_parse(ctx); ctx->event.type != YAML_SEQUENCE_END_EVENT; ylt_parse(ctx)) {
        ylt_output_mode_t initial_output_mode = ctx->output_mode;
        bool is_lua_invocation = ylt_is_lua_invocation(ctx);
        if (is_lua_invocation)
            ctx->output_mode = YLT_LUA_OUTPUT_MODE;

        ylt_evaluate_nested(ctx, "sequence");

        // If this was a Lua invocation, execute the Lua and output the value.
        if (is_lua_invocation) {
            ylt_execute_lua(ctx);
            ctx->output_mode = initial_output_mode;
            ylt_render_lua_value(ctx);
        }
    }

    ylt_parse(ctx); // Expect SEQUENCE END EVENT.
    if (ylt_unlikely(ctx->event.type != YAML_SEQUENCE_END_EVENT))
        return ylt_event_error(ctx, "Unexpected event at end of sequnece");
    ylt_emit(ctx); // Output SEQUENCE END EVENT.
}


void ylt_evaluate_mapping(ylt_context_t *ctx)
{
    if (ylt_unlikely(ctx->event.type != YAML_MAPPING_START_EVENT))
        return ylt_event_error(ctx, "Unexpected event at start of mapping");

    ylt_emit(ctx);

    for (ylt_parse(ctx); ctx->event.type != YAML_MAPPING_END_EVENT; ylt_parse(ctx)) {
        ylt_output_mode_t initial_output_mode = ctx->output_mode;
        size_t initial_buffer_len = ctx->event_buffer.len;
        bool discard_entry = false;

        // Read the key.
        // Always buffer the key in case the value is void.
        bool is_lua_invocation = ylt_is_lua_invocation(ctx);
        ctx->output_mode = is_lua_invocation ? YLT_LUA_OUTPUT_MODE : YLT_BUFFER_OUTPUT_MODE;

        ylt_evaluate_nested(ctx, "mapping key");

        // If this was a Lua invocation, execute the Lua and output the value to the buffer.
        // If the output is VOID, record that we should discard the value.
        if (is_lua_invocation) {
            ylt_execute_lua(ctx);
            discard_entry = ylt_lua_value_is_void(ctx);
            ctx->output_mode = YLT_BUFFER_OUTPUT_MODE;
            ylt_render_lua_value(ctx);
        }

        // Read the value corresponding to the above key.
        ylt_parse(ctx);

        if (ylt_unlikely(discard_entry)) {
            ylt_discard_nested(ctx);
        } else {
            is_lua_invocation = ylt_is_lua_invocation(ctx);
            if (is_lua_invocation) {
                ctx->output_mode = YLT_LUA_OUTPUT_MODE;
            } else {
                // The output can't be VOID here, so flush the buffer and continue.
                ctx->output_mode = initial_output_mode;
                ylt_playback_event_buffer(ctx, initial_buffer_len);
            }

            ylt_evaluate_nested(ctx, "mapping value");

            // If this was a Lua invocation, restore the output mode and execute the Lua. If the
            // output is VOID, discard the buffer. Otherwise, output the buffer and output the value.
            if (is_lua_invocation) {
                ctx->output_mode = initial_output_mode;
                ylt_execute_lua(ctx);
                if (ylt_unlikely(ylt_lua_value_is_void(ctx)))
                    ylt_truncate_event_buffer(ctx, initial_buffer_len);
                else
                    ylt_playback_event_buffer(ctx, initial_buffer_len);
                ylt_render_lua_value(ctx);
            }
        }
    }
}


void ylt_evaluate_scalar(ylt_context_t *ctx)
{
    if (ylt_unlikely(ctx->event.type != YAML_SCALAR_EVENT))
        return ylt_event_error(ctx, "Unexpected event while reading scalar");
}


void ylt_buffer_event(ylt_context_t *ctx)
{
    if (ylt_unlikely(ctx->event_buffer.len == ctx->event_buffer.cap)) {
        ctx->event_buffer.cap = (ctx->event_buffer.cap << 1) || 2;
        yaml_event_t *bigger_events = realloc(ctx->event_buffer.events, ctx->event_buffer.cap * sizeof(yaml_event_t));
        if (ylt_unlikely(bigger_events == NULL)) // Memory allocation failed. Original buffer was not freed.
            ylt_event_error(ctx, "Failed to allocate event buffer");
        ctx->event_buffer.events = bigger_events;
    }

    ctx->event_buffer.events[ctx->event_buffer.len++] = ctx->event;
    ctx->event = (yaml_event_t){0};
}


void ylt_playback_event_buffer(ylt_context_t *ctx, size_t since)
{
    if (ylt_unlikely(ctx->output_mode == YLT_BUFFER_OUTPUT_MODE))
        return; // Nothing to do. This is silly, but allowed for convenience.

    if (ylt_unlikely(ctx->event.type != YAML_NO_EVENT))
        ylt_event_error(ctx, "Unexpected non-empty event when playing back buffer");

    for (size_t i = since; i < ctx->event_buffer.len; ++i) {
        ctx->event = ctx->event_buffer.events[i];
        ctx->event_buffer.events[i] = (yaml_event_t){0}; // Avoid double free if there is an exception in ylt_emit().
        ylt_emit(ctx);
    }

    ctx->event_buffer.len = since;
}


void ylt_truncate_event_buffer(ylt_context_t *ctx, size_t since)
{
    if (ylt_unlikely(since > ctx->event_buffer.len))
        ylt_event_error(ctx, "Truncation length cannot exceed current buffer length");
    for (size_t i = since; i < ctx->event_buffer.len; ++i) {
        yaml_event_delete(&ctx->event_buffer.events[i]);
    }
    ctx->event_buffer.len = since;
}
