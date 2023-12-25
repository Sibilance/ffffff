#include "event.h"

const char *yl_event_names[] = {
    "NO_EVENT",
    "STREAM_START_EVENT",
    "STREAM_END_EVENT",
    "DOCUMENT_START_EVENT",
    "DOCUMENT_END_EVENT",
    "ALIAS_EVENT",
    "SCALAR_EVENT",
    "SEQUENCE_START_EVENT",
    "SEQUENCE_END_EVENT",
    "MAPPING_START_EVENT",
    "MAPPING_END_EVENT",
};

int yl_copy_event(yaml_event_t *original, yaml_event_t *copy)
{
    switch (original->type) {
    case YAML_DOCUMENT_START_EVENT:
        return yaml_document_start_event_initialize(copy,
                                                    original->data.document_start.version_directive,
                                                    original->data.document_start.tag_directives.start,
                                                    original->data.document_start.tag_directives.end,
                                                    original->data.document_start.implicit);
    case YAML_ALIAS_EVENT:
        return yaml_alias_event_initialize(copy, original->data.alias.anchor);
    case YAML_SCALAR_EVENT:
        return yaml_scalar_event_initialize(copy,
                                            original->data.scalar.anchor,
                                            original->data.scalar.tag,
                                            original->data.scalar.value,
                                            original->data.scalar.length,
                                            original->data.scalar.plain_implicit,
                                            original->data.scalar.quoted_implicit,
                                            original->data.scalar.style);
    case YAML_SEQUENCE_START_EVENT:
        return yaml_sequence_start_event_initialize(copy,
                                                    original->data.sequence_start.anchor,
                                                    original->data.sequence_start.tag,
                                                    original->data.sequence_start.implicit,
                                                    original->data.sequence_start.style);
    case YAML_MAPPING_START_EVENT:
        return yaml_mapping_start_event_initialize(copy,
                                                   original->data.mapping_start.anchor,
                                                   original->data.mapping_start.tag,
                                                   original->data.mapping_start.implicit,
                                                   original->data.mapping_start.style);
    default:
        *copy = *original; // No need to copy internal fields.
        return 1;
    }
}

const char *yl_event_name(yaml_event_type_t event_type)
{
    return yl_event_names[event_type];
}

int yl_record_event(yl_event_record_t *event_record, yaml_event_t *event, yl_error_t *err)
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

void yl_event_record_delete(yl_event_record_t *event_record)
{
    for (size_t i = 0; i < event_record->length; ++i) {
        yaml_event_delete(&event_record->events[i]);
    }

    if (event_record->events != NULL)
        free(event_record->events);

    *event_record = (yl_event_record_t){0};
}

int yl_replay_event(yl_event_record_t *event_record, yaml_event_t *event, yl_error_t *err)
{
    if (event_record->index == event_record->length)
        goto error;

    yaml_event_t *original_event = &event_record->events[event_record->index++];
    if (!yl_copy_event(original_event, event)) {
        err->type = YL_MEMORY_ERROR; // Probably.
        err->line = original_event->start_mark.line;
        err->column = original_event->start_mark.column;
        err->context = "While replaying an event, ran into problem";
        err->message = "memory error";
    }

    return 1;

error:
    return 0;
}
