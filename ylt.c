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

void ylt_evaluate_document(ylt_context_t *ctx)
{
    if (ctx->event.type != YAML_DOCUMENT_START_EVENT)
        return ylt_event_error(ctx, "Unexpected event (expecting DOCUMENT_START_EVENT)");

    ylt_output_mode_t initial_output_mode = ctx->output_mode;
    size_t initial_buffer_len = ctx->event_buffer.len;

    // Buffer the document start event in case the document content indicates it should be skipped.
    ylt_buffer_event(ctx);

    ylt_parse(ctx);

    // If the next event has no Lua invocation tag, the content will be passed-through. Flush the buffer (output the
    // DOCUMENT START EVENT) and continue processing.
    // Otherwise, change to LUA OUTPUT MODE and continue processing. On return, check the Lua stack for the return value.
    // If the return value is VOID, purge the buffer (discard the DOCUMENT START EVENT), restore the initial output mode,
    // and return. Otherwise, flush the buffer (output the DOCUMENT START EVENT), restore the initial output mode, render
    // the Lua value, and continue processing.

    switch (ctx->event.type) {
    case YAML_SEQUENCE_START_EVENT:
        if (ylt_is_lua_invocation(ctx->event.data.sequence_start.tag))
            ctx->output_mode = YLT_LUA_OUTPUT_MODE;
        else
            ylt_playback_event_buffer(ctx);
        ylt_evaluate_sequence(ctx);
        break;
    case YAML_MAPPING_START_EVENT:
        if (ylt_is_lua_invocation(ctx->event.data.mapping_start.tag))
            ctx->output_mode = YLT_LUA_OUTPUT_MODE;
        else
            ylt_playback_event_buffer(ctx);
        ylt_evaluate_mapping(ctx);
        break;
    case YAML_SCALAR_EVENT:
        if (ylt_is_lua_invocation(ctx->event.data.scalar.tag))
            ctx->output_mode = YLT_LUA_OUTPUT_MODE;
        else
            ylt_playback_event_buffer(ctx);
        ylt_evaluate_scalar(ctx);
        break;
    default:
        return ylt_event_error(ctx, "Unexpected event while processing document");
    }

    // If we were in LUA OUTPUT MODE, this was a Lua invocation. Restore the output mode to EMITTER OUTPUT MODE.
    // Evaluate the Lua invocation and check the stack for the return value. If the return value is VOID, purge the
    // buffer (discard the DOCUMENT START EVENT) and return (skipping the document). Otherwise, flush the buffer
    // (output the DOCUMENT START EVENT) and then render the Lua value.
    bool discard_document_end = false;
    if (ctx->output_mode == YLT_LUA_OUTPUT_MODE) {
        ctx->output_mode = initial_output_mode;
        ylt_evaluate_lua(ctx);
        if (ylt_unlikely(ylt_lua_value_is_void(ctx))) {
            ylt_truncate_event_buffer(ctx, initial_buffer_len);
            discard_document_end = true;
        } else {
            ylt_playback_event_buffer(ctx);
            ylt_render_lua_value(ctx);
        }
    }

    ylt_parse(ctx); // Expect DOCUMENT END EVENT.
    if (ylt_unlikely(ctx->event.type != YAML_DOCUMENT_END_EVENT))
        return ylt_event_error(ctx, "Unexpected event (expecting DOCUMENT_END_EVENT)");
    if (ylt_unlikely(discard_document_end))
        yaml_event_delete(&ctx->event);
    else
        ylt_emit(ctx); // Output DOCUMENT END EVENT.
}

void ylt_buffer_event(ylt_context_t *ctx)
{
    if (ylt_unlikely(ctx->event_buffer.len == ctx->event_buffer.cap)) {
        ctx->event_buffer.cap = (ctx->event_buffer.cap << 1) || 2;
        yaml_event_t *bigger_events = realloc(ctx->event_buffer.events, ctx->event_buffer.cap * sizeof(yaml_event_t));
        if (bigger_events == NULL) // Memory allocation failed. Original buffer was not freed.
            ylt_event_error(ctx, "Failed to allocate event buffer");
        ctx->event_buffer.events = bigger_events;
    }

    ctx->event_buffer.events[ctx->event_buffer.len++] = ctx->event;
    ctx->event = (yaml_event_t){0};
}

void ylt_playback_event_buffer(ylt_context_t *ctx)
{
    if (ylt_unlikely(ctx->output_mode == YLT_BUFFER_OUTPUT_MODE))
        return; // Nothing to do. This is silly, but allowed for convenience.

    if (ylt_unlikely(ctx->event.type != YAML_NO_EVENT))
        ylt_event_error(ctx, "Unexpected non-empty event when playing back buffer");

    for (size_t i = 0; i < ctx->event_buffer.len; ++i) {
        ctx->event = ctx->event_buffer.events[i];
        ctx->event_buffer.events[i] = (yaml_event_t){0}; // Avoid double free if there is an exception in ylt_emit().
        ylt_emit(ctx);
    }

    ctx->event_buffer.len = 0;
}

void ylt_truncate_event_buffer(ylt_context_t *ctx, size_t len)
{
    if (ylt_unlikely(len > ctx->event_buffer.len))
        ylt_event_error(ctx, "Truncation length cannot exceed current buffer length");
    for (size_t i = len; i < ctx->event_buffer.len; ++i) {
        yaml_event_delete(&ctx->event_buffer.events[i]);
    }
    ctx->event_buffer.len = len;
}
