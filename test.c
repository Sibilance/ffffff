#include "test.h"

int yl_test_stream(yl_execution_context_t *ctx)
{
    yaml_event_t next_event = {0};

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
        if (!yl_parser_parse(&ctx->parser, &next_event, &ctx->err))
            goto error;

        switch (next_event.type) {
        case YAML_STREAM_START_EVENT:
            if (!ctx->handler(ctx->data, &next_event, &ctx->err))
                goto error;
            break;
        case YAML_DOCUMENT_START_EVENT:
            if (!yl_execute_document(ctx, &next_event))
                goto error;
            break;
        case YAML_STREAM_END_EVENT:
            if (!ctx->handler(ctx->data, &next_event, &ctx->err))
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
    }

    return 1;

error:
    yaml_event_delete(&next_event);
    return 0;
}
