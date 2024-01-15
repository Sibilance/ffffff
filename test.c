#include <stdbool.h>
#include <string.h>

#include "lua.h"
#include "yaml.h"

#include "event.h"
#include "render.h"
#include "test.h"

int yl_test_stream(yl_execution_context_t *ctx)
{
    yaml_event_t next_event = {0};

    yl_event_record_t actual_events = {0};
    yl_event_record_t expected_events = {0};
    char *expected_rendering = NULL, *actual_rendering = NULL;

    if (!ctx->producer(ctx->producer_data, &next_event, &ctx->err))
        goto error;

    if (next_event.type != YAML_STREAM_START_EVENT) {
        ctx->err.type = YL_PARSER_ERROR;
        ctx->err.line = next_event.start_mark.line;
        ctx->err.column = next_event.start_mark.column;
        ctx->err.context = "While trying to read a stream, got unexpected event";
        ctx->err.message = yl_event_name(next_event.type);
        goto error;
    }

    if (!ctx->consumer.callback(ctx->consumer.data, &next_event, NULL, &ctx->err))
        goto error;

    while (true) {
        size_t line = 0, column = 0;
        if (!ctx->producer(ctx->producer_data, &next_event, &ctx->err))
            goto error;

        switch (next_event.type) {
        case YAML_DOCUMENT_START_EVENT:
            line = next_event.start_mark.line;
            column = next_event.start_mark.column;

            if (!yl_test_record_document(ctx, &next_event, &actual_events))
                goto error;

            break;

        case YAML_STREAM_END_EVENT:
            if (!ctx->consumer.callback(ctx->consumer.data, &next_event, NULL, &ctx->err))
                goto error;
            goto done;

        default:
            ctx->err.type = YL_PARSER_ERROR;
            ctx->err.line = next_event.start_mark.line;
            ctx->err.column = next_event.start_mark.column;
            ctx->err.context = "While trying to read a document, got unexpected event";
            ctx->err.message = yl_event_name(next_event.type);
            goto error;
        }

        actual_rendering = yl_event_record_to_string(&actual_events, &ctx->err);
        if (actual_rendering == NULL)
            goto error;
        for (size_t i = 0; i < actual_events.length; ++i) {
            if (!ctx->consumer.callback(ctx->consumer.data, &actual_events.events[i], NULL, &ctx->err))
                goto error;
        }

        if (!ctx->producer(ctx->producer_data, &next_event, &ctx->err))
            goto error;

        switch (next_event.type) {
        case YAML_DOCUMENT_START_EVENT:
            if (!yl_test_record_document(ctx, &next_event, &expected_events))
                goto error;
            break;

        default:
            ctx->err.type = YL_PARSER_ERROR;
            ctx->err.line = next_event.start_mark.line;
            ctx->err.column = next_event.start_mark.column;
            ctx->err.context = "While trying to read a document, got unexpected event";
            ctx->err.message = yl_event_name(next_event.type);
            goto error;
        }

        expected_rendering = yl_event_record_to_string(&expected_events, &ctx->err);
        if (expected_rendering == NULL)
            goto error;
        for (size_t i = 0; i < expected_events.length; ++i) {
            if (!ctx->consumer.callback(ctx->consumer.data, &expected_events.events[i], NULL, &ctx->err))
                goto error;
        }

        if (strcmp(actual_rendering, expected_rendering) != 0) {
            ctx->err.type = YL_ASSERTION_ERROR;
            ctx->err.line = line;
            ctx->err.column = column;
            ctx->err.context = "While comparing rendered documents";
            ctx->err.message = "actual document differs from expected document";
            goto error;
        }

        free(actual_rendering);
        actual_rendering = NULL;
        free(expected_rendering);
        expected_rendering = NULL;
        yl_event_record_delete(&actual_events);
        yl_event_record_delete(&expected_events);
    }
done:
    return 1;

error:
    if (actual_rendering != NULL)
        free(actual_rendering);
    if (expected_rendering != NULL)
        free(expected_rendering);
    yl_event_record_delete(&actual_events);
    yl_event_record_delete(&expected_events);
    yaml_event_delete(&next_event);
    return 0;
}

int yl_test_record_document(yl_execution_context_t *ctx, yaml_event_t *next_event, yl_event_record_t *event_record)
{
    yl_event_consumer_t saved_consumer = ctx->consumer;

    ctx->consumer.callback = (yl_event_consumer_callback_t *)yl_record_event;
    ctx->consumer.data = event_record;

    if (!yl_execute_document(ctx, next_event))
        goto error;

    ctx->consumer = saved_consumer;
    yaml_event_delete(next_event);

    return 1;

error:
    ctx->consumer = saved_consumer;
    yaml_event_delete(next_event);

    return 0;
}
