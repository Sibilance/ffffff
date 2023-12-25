#include <stdbool.h>

#include "yaml.h"

#include "event.h"
#include "test.h"

int yl_test_stream(yl_execution_context_t *ctx)
{
    yaml_event_t next_event = {0};

    yl_execution_context_t saved_ctx = {0};
    bool recording_actual = true;
    yl_event_record_t actual_events = {0};
    yl_event_record_t expected_events = {0};

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
            if (!ctx->consumer(ctx->consumer_data, &next_event, &ctx->err))
                goto error;
            break;
        case YAML_DOCUMENT_START_EVENT: {
            saved_ctx = *ctx;

            ctx->consumer = (yl_event_producer_t *)yl_record_event;
            ctx->consumer_data = recording_actual ? &actual_events : &expected_events;

            if (!yl_execute_document(ctx, &next_event))
                goto error;

            ctx->consumer = saved_ctx.consumer;
            ctx->consumer_data = saved_ctx.consumer_data;

            if (recording_actual) {
                for (size_t i = 0; i < actual_events.length; ++i) {
                    if (!ctx->consumer(ctx->consumer_data, &actual_events.events[i], &ctx->err))
                        goto error;
                }
                recording_actual = false;
            } else {
                for (size_t i = 0; i < expected_events.length; ++i) {
                    if (!ctx->consumer(ctx->consumer_data, &expected_events.events[i], &ctx->err))
                        goto error;
                }
                // TODO: compare outputs

                yl_event_record_delete(&actual_events);
                yl_event_record_delete(&expected_events);
                recording_actual = true;
            }
        } break;
        case YAML_STREAM_END_EVENT:
            if (!ctx->consumer(ctx->consumer_data, &next_event, &ctx->err))
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
    yl_event_record_delete(&actual_events);
    yl_event_record_delete(&expected_events);
    yaml_event_delete(&next_event);
    return 0;
}
