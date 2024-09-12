#include <ctype.h>
#include <stdlib.h>
#include <string.h>

#include "error.h"
#include "executor.h"

int yl_execute_stream(yl_execution_context_t *ctx)
{
    yaml_event_t next_event = {0};

    bool done = false;
    while (!done) {
        if (!ctx->producer.callback(ctx->producer.data, &next_event, &ctx->err))
            goto error;

        switch (next_event.type) {
        case YAML_STREAM_START_EVENT:
            if (!ctx->consumer.callback(ctx->consumer.data, &next_event, &ctx->err))
                goto error;
            break;
        case YAML_DOCUMENT_START_EVENT:
            if (!yl_execute_document(ctx, &next_event))
                goto error;
            break;
        case YAML_STREAM_END_EVENT:
            if (!ctx->consumer.callback(ctx->consumer.data, &next_event, &ctx->err))
                goto error;
            done = true;
            break;
        default:
            ctx->err.type = YL_EXECUTION_ERROR;
            ctx->err.line = next_event.start_mark.line;
            ctx->err.column = next_event.start_mark.column;
            ctx->err.context = "While executing a stream, got unexpected event";
            ctx->err.message = yl_event_name(next_event.type);
            goto error;
        }

        yaml_event_delete(&next_event);
    }

    return 1;

error:
    yaml_event_delete(&next_event);
    return 0;
}

int yl_execute_document(yl_execution_context_t *ctx, yaml_event_t *event)
{
    (void)ctx;
    (void)event;
    return 0;
}

int yl_execute_sequence(yl_execution_context_t *ctx, yaml_event_t *event)
{
    (void)ctx;
    (void)event;
    return 0;
}

int yl_execute_mapping(yl_execution_context_t *ctx, yaml_event_t *event)
{
    (void)ctx;
    (void)event;
    return 0;
}

int yl_execute_scalar(yl_execution_context_t *ctx, yaml_event_t *event)
{
    (void)ctx;
    (void)event;
    return 0;
}
