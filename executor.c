#include "lua.h"

#include "executor.h"

int yl_execute_stream(yl_execution_context_t *ctx)
{
    yl_event_t event = {0};

    bool done = false;
    while (!done) {
        if (!yl_parser_parse(&ctx->parser, &event, &ctx->err))
            return 0;

        switch (event.type) {
        case YAML_STREAM_START_EVENT:
            // Ignore the start of the stream, as there is no reason the caller
            // already would have consumed this before calling yl_execute_stream().
            break;
        case YAML_DOCUMENT_START_EVENT:
            if (!yl_execute_document(ctx))
                return 0;
            break;
        case YAML_STREAM_END_EVENT:
            done = true;
            break;
        default:
            ctx->err.type = YL_PARSER_ERROR;
            ctx->err.line = event.line;
            ctx->err.column = event.column;
            ctx->err.context = "While executing a stream, got unexpected event";
            ctx->err.message = yl_event_name(event.type);
            goto error;
        }

        if (!handler(ctx->data, &event, &ctx->err))
            goto error;

        yl_event_delete(&event);
    }

    return 1;

error:
    yl_event_delete(&event);
    return 0;
}

int yl_execute_document(yl_execution_context_t *ctx)
{
    yl_event_t event = {0};

    bool done = false;
    while (!done) {
        if (!yl_parser_parse(&ctx->parser, &event, &ctx->err))
            return 0;

        switch (event.type) {
        case YAML_SCALAR_EVENT:
            break;
        case YAML_SEQUENCE_START_EVENT:
            break;
        case YAML_MAPPING_START_EVENT:
            break;
        case YAML_DOCUMENT_END_EVENT:
            done = true;
            break;
        default:
            ctx->err.type = YL_PARSER_ERROR;
            ctx->err.line = event.line;
            ctx->err.column = event.column;
            ctx->err.context = "While executing a document, got unexpected event";
            ctx->err.message = yl_event_name(event.type);
            goto error;
        }

        if (!handler(ctx->data, &event, &ctx->err))
            goto error;

        yl_event_delete(&event);
    }
    return 1;

error:
    yl_event_delete(&event);
    return 0;
}
