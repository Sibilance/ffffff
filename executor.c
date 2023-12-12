#include "lua.h"

#include "executor.h"

yl_error_t yl_execute_stream(yaml_parser_t *parser, yaml_emitter_t *emitter)
{
    yl_error_t err;
    yl_event_t event;

    bool done = false;
    while (!done) {
        err = yl_parser_parse(parser, &event);
        if (err.type)
            return err;

        switch (event.type) {
        case YAML_STREAM_START_EVENT:
            // Ignore the start of the stream, as there is no reason the caller
            // already would have consumed this before calling yl_execute_stream().
            break;
        case YAML_DOCUMENT_START_EVENT:
            break;
        case YAML_STREAM_END_EVENT:
            done = true;
            break;
        }

        yl_event_delete(&event);
    }

    return YL_SUCCESS;
}
