#include "parser.h"

/**
 * The prototype of an event handler.
 *
 * The event handler is called when the executor has finished processing or producing
 * an event. The handler can write the event to an output file, or accumulate it into
 * some sort of object for further processing.
 *
 * @param[in,out]   data        A pointer to an application data passed to yl_execute_*().
 * @param[in]       event       The event emitted.
 * @param[out]      err         Error details.
 *
 * @returns On success, the handler should return @c 1. If the handler failed,
 * the returned value should be @c 0.
 */
typedef int yl_event_handler_t(void *data, yl_event_t *event, yl_error_t *err);

int yl_execute_stream(yaml_parser_t *parser, yl_event_handler_t *handler, void *data, yl_error_t *err);
int yl_execute_document(yaml_parser_t *parser, yl_event_handler_t *handler, void *data, yl_error_t *err);
