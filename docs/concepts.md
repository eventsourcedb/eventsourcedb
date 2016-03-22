# EventSourceDB documentation

## Table of Contents:
- [Event](#event)
- [Stream](#stream)
- [Topic Masks](#topic-masks)
- [Publisher](#publisher)
- [Batch Publisher](#batch-publisher])
- [Consumer](#consumer)
- [Realtime](#realtime)
- [Consumers grouping](#consumers-grouping)
- [Correlation](#correlation)
- [Examples](#examples)

## Event
Event is immutable, atomic unit of [EventSourcing][1] pattern. It's a
accomplished fact of state change in the system. Events always signify something
that happened in past-tense. In JSON form an event will look like this:

```json
{
    "id": "someID",
    "stream": "some.stream.name",
    "type": "SomethingHappened",
    "body": {
        "arbitrary": "fields"
    }
}
```

## Stream
Stream is a collection of events grouped together by meaning. An aggregate
root in Domain Driven Design (DDD) world.

Stream names could contain `.` in their names, allowing consumers to address
multiple streams using topic masks.

Example stream names:
- `stream1` just simple name. Used by consumers and publishers.
- `some.topic.stream` name, formed as topic. Used by consumers and publishers.
- `some.topic,some.other.stream` multiple topics. Used by consumers.
- `some.*.stream` topic mask, will match anything but dot. Used by consumers.
- `some.#.stream` topic mask, will match anything including dot (wild-card
    match). Used by consumers.

## Topic Masks
Topic masks is a wild-card on topic names. `*` symbol match anything but dot.
`#` match anything.

Example topic masks in topic names:
`some.*.name` will match `some.stream.name` but not `some.stream.other.name`.
To match latest, `some.#.name` will be needed.

## Publisher
Publishers publish their events into streams using HTTP POST method.

Example request:
```
POST /pub/accounting.book.10

{
    "type": "AccountCredited",
    "body": {
        "user": "alice",
        "ammount": "10",
        "currency": "USD"
    }
}
```

Server will automatically assign an ID.

## Batch Publisher
Multiple JSON documents separated by `\n` are allowed.

## Consumer
Consumers subscribe to streams (single or multiple) to receive published events
in realtime. Also consumers can receive events from past. Consumers could be
grouped (multiple instances of the same consumer).

By default, unless grouped, consumers will receive events in fanout mode.

Consumers within single group will receive events in round-robin mode,
independently from other consumers.

## Realtime
By default consumer will receive events from server indefinitely. Server
will send them as long as consumer is connected. This can be disabled by
supplying HTTP header `X-RTDisabled: 1`. After reaching the head of stream
client will be disconnected.

## Consumers grouping
Consumer grouping used to subscribe multiple instances of the same consumer, to
speed-up (parallelize) processing. Grouping could be enabled by sending HTTP
header `X-GroupID: SomeGroupConsumer`.

## Correlation
Sometimes it's extremely important to process events sequentially grouped by
some field value (for example "user" field). This can be achieved by supplying
HTTP header `X-CorrelationFields: user1,user2`, where `user1` and `user2`
fields are payload inside JSON `body`. This will only work in conjunction with
grouping.

The way it works is: a hash is calculated of those fields, and event is
dispatched to consistent-hashing ring of consumers.

## HTTP Headers
- `User-Agent: MyCoolWorker v0.1` (will be used in status)
- `X-GroupID: SomeGroup`
- `X-CorrelationFields: field1,field2`
- `X-RTDisabled: 1`

## Examples
TBD

[1]: http://martinfowler.com/eaaDev/EventSourcing.html
