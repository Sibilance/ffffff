#include "lua.h"

#include "executor.h"

int yl_execute_stream(yaml_parser_t *parser, yl_event_handler_t *handler, void *data, yl_error_t *err)
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
            if (!yl_execute_document(parser, handler, data, err))
                return 0;
            break;
        case YAML_STREAM_END_EVENT:
            done = true;
            break;
        default:
            break;
        }

        yl_event_delete(&event);
    }

    return 1;
}

int yl_execute_document(yaml_parser_t *parser, yl_event_handler_t *handler, void *data, yl_error_t *err)
{
    yl_event_t event = {0};

    bool done = false;
    while (!done) {
        if (!yl_parser_parse(parser, &event, err))
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
            break;
        }
    }
    return 1;
}
