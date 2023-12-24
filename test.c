#include <stdbool.h>

#include "yaml.h"

#include "test.h"

typedef struct _event_record_s {
    size_t capacity;
    size_t length;
    size_t index;
    yaml_event_t *events;
} event_record_t;

static int record_event(event_record_t *event_record, yaml_event_t *event, yl_error_t *err)
{
    if (event_record->length == event_record->capacity) {
        size_t new_capacity = event_record->capacity << 1;
        if (new_capacity == 0)
            new_capacity = 1;

        yaml_event_t *resized_events = realloc(event_record->events, sizeof(yaml_event_t) * new_capacity);
        if (resized_events == NULL) {
            err->type = YL_MEMORY_ERROR;
            err->line = event->start_mark.line;
            err->column = event->start_mark.column;
            err->context = "While recording events from a stream, got memory error";
            err->message = "unable to realloc event list";
            goto error;
        }

        event_record->events = resized_events;
        event_record->capacity = new_capacity;
    }

    event_record->events[event_record->length++] = *event;
    *event = (yaml_event_t){0}; // Mark the event as consumed to prevent its contents from being freed.

    return 1;

error:
    return 0;
}

static void event_record_delete(event_record_t *event_record)
{
    for (size_t i = 0; i < event_record->length; ++i) {
        yaml_event_delete(&event_record->events[i]);
    }

    if (event_record->events != NULL)
        free(event_record->events);

    *event_record = (event_record_t){0};
}

// static int replay_event(event_record_t *event_record, yaml_event_t *event, yl_error_t *err)
// {
//     (void)err; // Unused.

//     if (event_record->index == event_record->length)
//         goto error;

//     yl_copy_event(&event_record->events[event_record->index++], event);

//     return 1;

// error:
//     return 0;
// }

int yl_test_stream(yl_execution_context_t *ctx)
{
    yaml_event_t next_event = {0};

    yl_execution_context_t saved_ctx = {0};
    bool recording_actual = true;
    event_record_t actual_events = {0};
    event_record_t expected_events = {0};

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

            ctx->consumer = (yl_event_producer_t *)record_event;
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

                event_record_delete(&actual_events);
                event_record_delete(&expected_events);
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

    event_record_delete(&actual_events);
    event_record_delete(&expected_events);

    return 1;

error:
    event_record_delete(&actual_events);
    event_record_delete(&expected_events);
    yaml_event_delete(&next_event);
    return 0;
}
