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

    yl_event_consumer_t saved_consumer = {0};
    bool recording_actual = true;
    yl_event_record_t actual_events = {0};
    yl_event_record_t expected_events = {0};
    char *expected_rendering = NULL, *actual_rendering = NULL;

    // TODO: Currently just copied from yl_execute_stream().
    // The challenge here is how to implement the !testcases use-case.
    // There is no functionality to "rewind" the event stream, so we
    // need to create one in order to allow the same document to be
    // executed multiple times with different global variables.
    // The output also has to be captured and compared, which is already
    // possible using the callback handler; however, we'll also want
    // to wrap the existing handler so the user can see the test output.
    bool done = false;
    while (!done) {
        if (!ctx->producer(ctx->producer_data, &next_event, &ctx->err))
            goto error;

        switch (next_event.type) {
        case YAML_STREAM_START_EVENT:
            if (!ctx->consumer.callback(ctx->consumer.data, &next_event, NULL, &ctx->err))
                goto error;
            break;
        case YAML_DOCUMENT_START_EVENT: {
            saved_consumer = ctx->consumer;

            ctx->consumer.callback = (yl_event_consumer_callback_t *)yl_record_event;
            ctx->consumer.data = recording_actual ? &actual_events : &expected_events;

            size_t line = next_event.start_mark.line;
            size_t column = next_event.start_mark.column;

            if (!yl_execute_document(ctx, &next_event))
                goto error;

            ctx->consumer = saved_consumer;

            if (recording_actual) {
                recording_actual = false;
            } else {
                actual_rendering = yl_event_record_to_string(&actual_events, &ctx->err);
                if (actual_rendering == NULL)
                    goto error;

                expected_rendering = yl_event_record_to_string(&expected_events, &ctx->err);
                if (expected_rendering == NULL)
                    goto error;

                // Consume both the sets of events.
                for (size_t i = 0; i < actual_events.length; ++i) {
                    if (!ctx->consumer.callback(ctx->consumer.data, &actual_events.events[i], NULL, &ctx->err))
                        goto error;
                }
                for (size_t i = 0; i < expected_events.length; ++i) {
                    if (!ctx->consumer.callback(ctx->consumer.data, &expected_events.events[i], NULL, &ctx->err))
                        goto error;
                }

                if (strcmp(actual_rendering, expected_rendering) != 0) {
                    fprintf(stderr, "%s\n%s\n", actual_rendering, expected_rendering);
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
                recording_actual = true;
            }
        } break;
        case YAML_STREAM_END_EVENT:
            if (!ctx->consumer.callback(ctx->consumer.data, &next_event, NULL, &ctx->err))
                goto error;
            done = true;
            break;
        default:
            ctx->err.type = YL_EXECUTION_ERROR;
            ctx->err.line = next_event.start_mark.line;
            ctx->err.column = next_event.start_mark.column;
            ctx->err.context = "While testing a stream, got unexpected event";
            ctx->err.message = yl_event_name(next_event.type);
            goto error;
        }

        yaml_event_delete(&next_event);
    }

    yl_event_record_delete(&actual_events);
    yl_event_record_delete(&expected_events);

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
