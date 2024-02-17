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

void ylt_try_parse(ylt_context_t *ctx)
{
    if (!yaml_parser_parse(&ctx->parser, &ctx->event))
        return ylt_parser_error(ctx);
}

void ylt_try_emit(ylt_context_t *ctx)
{
    yaml_event_type_t event_type = ctx->event.type;
    if (!yaml_emitter_emit(&ctx->emitter, &ctx->event))
        return ylt_emitter_error(ctx, ylt_yaml_event_names[event_type]);
    ctx->event = (yaml_event_t){0}; // Event has been consumed (either deleted, or added to an internal libyaml queue).
}

void ylt_evaluate_stream(ylt_context_t *ctx)
{
    if (ctx->event.type != YAML_NO_EVENT)
        return ylt_event_error(ctx, "Unexpected event already parsed when evaluating stream");

    ylt_try_parse(ctx);

    if (ctx->event.type != YAML_STREAM_START_EVENT)
        return ylt_event_error(ctx, "Unexpected event at start of stream (expecting STREAM_START_EVENT)");

    ylt_try_emit(ctx);

    bool done = false;
    while (!done) {
        ylt_try_parse(ctx);

        switch (ctx->event.type) {
        case YAML_DOCUMENT_START_EVENT:
            ylt_evaluate_document(ctx);
            break;
        case YAML_STREAM_END_EVENT:
            ylt_try_emit(ctx);
            done = true;
            break;
        default:
            return ylt_event_error(ctx, "Unexpected event while processing stream");
        }

        yaml_event_delete(&ctx->event);
    }
}

void ylt_evaluate_document(ylt_context_t *ctx)
{
    if (ctx->event.type != YAML_DOCUMENT_START_EVENT)
        return ylt_event_error(ctx, "Unexpected event (expecting DOCUMENT_START_EVENT)");

    ylt_try_emit(ctx); // TODO: do we need to buffer until we know whether to skip this document?
}
