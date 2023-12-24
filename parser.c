#include "parser.h"

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

int yl_parser_parse(yaml_parser_t *parser, yaml_event_t *event, yl_error_t *err)
{
    *event = (yaml_event_t){0};

    if (!yaml_parser_parse(parser, event)) {
        err->type = (yl_error_type_t)parser->error;
        err->line = parser->problem_mark.line;
        err->column = parser->problem_mark.column;
        err->context = parser->context;
        err->message = parser->problem;

        return 0;
    }

    return 1;
}

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
