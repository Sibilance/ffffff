#include <string.h>

#include "event.h"
#include "render.h"

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

int yl_record_event(yl_event_record_t *event_record, yaml_event_t *event, lua_State *L, yl_error_t *err)
{
    (void)L;

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
        goto error;
    }

    return 1;

error:
    return 0;
}

typedef struct _stringbuilder_s {
    size_t capacity;
    size_t length;
    char *contents;
} stringbuilder_t;

static int stringbuilder_add(stringbuilder_t *builder, char *buffer, size_t size)
{
    if (builder->length + size >= builder->capacity) {
        size_t new_capacity = builder->capacity << 1;
        if (new_capacity == 0)
            new_capacity = 1;
        while (new_capacity < builder->length + size)
            new_capacity <<= 1;

        char *resized_contents = realloc(builder->contents, new_capacity);
        if (resized_contents == NULL)
            return 0;

        builder->contents = resized_contents;
        builder->capacity = new_capacity;
    }

    memcpy(&builder->contents[builder->length], buffer, size);
    builder->length += size;

    return 1;
}

static void stringbuilder_delete(stringbuilder_t *builder)
{
    if (builder->contents)
        free(builder->contents);

    *builder = (stringbuilder_t){0};
}

char *yl_event_record_to_string(yl_event_record_t *event_record, yl_error_t *err)
{
    yaml_emitter_t emitter = {0};
    yaml_event_t event = {0};
    stringbuilder_t stringbuilder = {0};
    char *output = NULL;
    size_t line = 0, column = 0;
    if (event_record->length > event_record->index) {
        line = event_record->events[event_record->index].start_mark.line;
        column = event_record->events[event_record->index].start_mark.column;
    }

    if (!yaml_emitter_initialize(&emitter)) {
        err->type = YL_EMITTER_ERROR;
        err->line = line;
        err->column = column;
        err->context = "Error rendering event record";
        err->message = "failed to initialize emitter";
        goto error;
    }

    yaml_emitter_set_unicode(&emitter, true);
    yaml_emitter_set_encoding(&emitter, YAML_UTF8_ENCODING);
    yaml_emitter_set_output(&emitter, (yaml_write_handler_t *)stringbuilder_add, &stringbuilder);

    yaml_stream_start_event_initialize(&event, YAML_UTF8_ENCODING);
    if (!yaml_emitter_emit(&emitter, &event))
        goto error;

    while (event_record->index < event_record->length) {
        if (!yl_replay_event(event_record, &event, err))
            goto error;

        if (!yaml_emitter_emit(&emitter, &event)) {
            err->type = YL_EMITTER_ERROR;
            err->line = event.start_mark.line;
            err->column = event.start_mark.column;
            err->context = "While rendering event record, failed to emit event";
            err->message = yl_event_name(event.type);
            goto error;
        }
    }

    yaml_stream_end_event_initialize(&event);
    if (!yaml_emitter_emit(&emitter, &event))
        goto error;

    output = strndup(stringbuilder.contents, stringbuilder.length);

    if (output == NULL) {
        err->type = YL_MEMORY_ERROR;
        err->line = line;
        err->column = column;
        err->context = "Error rendering event record";
        err->message = "failed to allocate output string";
        goto error;
    }

    stringbuilder_delete(&stringbuilder);
    yaml_emitter_delete(&emitter);
    return output;

error:
    if (output != NULL)
        free(output);
    stringbuilder_delete(&stringbuilder);
    yaml_emitter_delete(&emitter);
    return NULL;
}

char *yl_copy_anchor(yaml_event_t *event)
{
    char *anchor = NULL;

    switch (event->type) {
    case YAML_SCALAR_EVENT:
        if (event->data.scalar.anchor != NULL)
            anchor = strdup((char *)event->data.scalar.anchor);
        break;
    case YAML_SEQUENCE_START_EVENT:
        if (event->data.sequence_start.anchor)
            anchor = strdup((char *)event->data.sequence_start.anchor);
        break;
    case YAML_MAPPING_START_EVENT:
        if (event->data.mapping_start.anchor)
            anchor = strdup((char *)event->data.mapping_start.anchor);
        break;
    default:
        break;
    }

    return anchor;
}
