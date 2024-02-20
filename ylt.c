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
    if (ctx->event.type != YAML_NO_EVENT)
        return ylt_event_error(ctx, "Unexpected event already parsed when evaluating stream");

    ylt_parse(ctx);

    if (ctx->event.type != YAML_STREAM_START_EVENT)
        return ylt_event_error(ctx, "Unexpected event at start of stream (expecting STREAM_START_EVENT)");

    ylt_emit(ctx);

    while (true) {
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

    // Buffer the document start event in case the document content indicates it should be skipped.
    ylt_buffer_event(ctx);

    ylt_parse(ctx);

    // If the next event has no Lua invocation tag, the content will be passed-through. Flush the buffer (output the
    // DOCUMENT START EVENT) and continue processing.
    // Otherwise, change to LUA OUTPUT MODE and continue processing. On return, check the Lua stack for the return value.
    // If the return value is VOID, purge the buffer (discard the DOCUMENT START EVENT), restore EMIT OUTPUT MODE, and
    // return. Otherwise, flush the buffer (output the DOCUMENT START EVENT), restore EMIT OUTPUT MODE, and then render
    // the Lua value.

    switch (ctx->event.type) {
    case YAML_SEQUENCE_START_EVENT:
        if (ylt_is_lua_invocation(ctx->event.data.sequence_start.tag))
            ctx->output_mode = YLT_LUA_OUTPUT_MODE;
        else
            ylt_flush_event_buffer(ctx);
        ylt_evaluate_sequence(ctx);
        break;
    case YAML_MAPPING_START_EVENT:
        if (ylt_is_lua_invocation(ctx->event.data.mapping_start.tag))
            ctx->output_mode = YLT_LUA_OUTPUT_MODE;
        else
            ylt_flush_event_buffer(ctx);
        ylt_evaluate_mapping(ctx);
        break;
    case YAML_SCALAR_EVENT:
        if (ylt_is_lua_invocation(ctx->event.data.scalar.tag))
            ctx->output_mode = YLT_LUA_OUTPUT_MODE;
        else
            ylt_flush_event_buffer(ctx);
        ylt_evaluate_scalar(ctx);
        break;
    default:
        return ylt_event_error(ctx, "Unexpected event while processing document");
    }

    ylt_parse(ctx);
    if (ctx->event.type != YAML_DOCUMENT_END_EVENT)
        return ylt_event_error(ctx, "Unexpected event (expecting DOCUMENT_END_EVENT)");

    // If we were in LUA OUTPUT MODE, this was a Lua invocation. Restore the output mode to EMITTER OUTPUT MODE.
    // Evaluate the Lua invocation and check the stack for the return value. If the return value is VOID, purge the
    // buffer (discard the DOCUMENT START EVENT) and return (skipping the document). Otherwise, flush the buffer
    // (output the DOCUMENT START EVENT) and then render the Lua value.
    if (ctx->output_mode == YLT_LUA_OUTPUT_MODE) {
        ctx->output_mode = YLT_EMITTER_OUTPUT_MODE;
        ylt_evaluate_lua(ctx);
        if (ylt_lua_value_is_void(ctx)) {
            ylt_truncate_event_buffer(ctx, 0);
            yaml_event_delete(&ctx->event); // Discard DOCUMENT END EVENT as well.
        } else {
            ylt_flush_event_buffer(ctx);
            ylt_render_lua_value(ctx);
            ylt_emit(ctx); // Output DOCUMENT END EVENT.
        }
    }
}
