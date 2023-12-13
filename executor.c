#include "lua.h"

#include "executor.h"

int yl_execute_stream(yaml_parser_t *parser, yaml_emitter_t *emitter, yl_error_t *err)
{
    yl_event_t event;

    bool done = false;
    while (!done) {
        if (!yl_parser_parse(parser, &event, err))
            return 0;

        switch (event.type) {
        case YAML_STREAM_START_EVENT:
            // Ignore the start of the stream, as there is no reason the caller
            // already would have consumed this before calling yl_execute_stream().
            break;
        case YAML_DOCUMENT_START_EVENT:
            if (!yl_execute_document(parser, emitter, err))
                return 0;
            break;
        case YAML_STREAM_END_EVENT:
            done = true;
            break;
        }

        yl_event_delete(&event);
    }

    return 1;
}

int yl_execute_document(yaml_parser_t *parser, yaml_emitter_t *emitter, yl_error_t *err)
{
    return 1;
}
