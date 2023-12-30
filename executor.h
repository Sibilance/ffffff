#pragma once

#include "lua.h"

#include "event.h"
#include "parser.h"

/**
 * The prototype of an event producer.
 *
 * The event producer is called when the executor needs a new event to process.
 *
 * @param[in,out]   data        A pointer to an application data.
 * @param[out]      event       The event produced.
 * @param[out]      err         Error details.
 *
 * @returns On success, the producer should return @c 1. If the producer failed,
 * the returned value should be @c 0.
 */
typedef int yl_event_producer_t(void *data, yaml_event_t *event, yl_error_t *err);

/**
 * The prototype of an event consumer.
 *
 * The event consumer is called when the executor has finished processing or producing
 * an event. The consumer can write the event to an output file, or accumulate it into
 * some sort of object for further processing.
 *
 * @param[in,out]   data        A pointer to an application data.
 * @param[in]       event       The event emitted.
 * @param[in]       L           A pointer to the Lua state, if it contains data to
 *                              augment the event.
 * @param[out]      err         Error details.
 *
 * @returns On success, the handler should return @c 1. If the handler failed,
 * the returned value should be @c 0.
 */
typedef int yl_event_consumer_t(void *data, yaml_event_t *event, lua_State *L, yl_error_t *err);

typedef struct _yl_execution_context_s {
    yl_event_producer_t *producer;
    void *producer_data;
    lua_State *lua;
    yl_event_consumer_t *consumer;
    void *consumer_data;
    yl_error_t err;
} yl_execution_context_t;

int yl_execute_stream(yl_execution_context_t *ctx);
int yl_execute_document(yl_execution_context_t *ctx, yaml_event_t *event);
int yl_execute_sequence(yl_execution_context_t *ctx, yaml_event_t *event);
int yl_execute_mapping(yl_execution_context_t *ctx, yaml_event_t *event);
int yl_execute_scalar(yl_execution_context_t *ctx, yaml_event_t *event);
