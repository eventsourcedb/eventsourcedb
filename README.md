# EventSourceDB

## What is it about?
This is a database to provide easy [EventSourcing][1] adoption at your
organizarion, to safely store the source of truth about state of your
application.

## Status
ALPHA-GRADE SOFTWARE. NOT FOR PRODUCTION USE (YET).

[![Build Status](https://travis-ci.org/eventsourcedb/eventsourcedb.svg?branch=master)](https://travis-ci.org/eventsourcedb/eventsourcedb)

## Features
* [ ] JSON-enabled streaming HTTP2 API.
* [ ] Pub/Sub.
* [ ] Proxy-friendly.
* [ ] Operations-friendly (backup, restore, scale as easy as adding new http proxies).
* [ ] Uses [Bolt][2] as transactional, fast, and reliable storage.
* [ ] Append-only workflow.
* [ ] Yet it's possible to delete retired streams.

## Install & Use
Build requirments:
* make
* [Glide][3]
* [Go][4] 1.6+

Following command will install binaries into your `$GOPATH/bin`.

    $ make deps install

## Why HTTP2?
HTTP2 brings many advantages for such RPC-ish system, multiplexing multiple
requests over a single TCP connection, flow control, header compression,
binary encoding. Yet everyone are familiar with HTTP semantics.

## Why not gRPC then?
To keep things simple. Simple requests, simple JSON, simple proxies. And to
maintain some level of compatibility with HTTP1 clients (although, not
recommended to use).

## Client libraries
At the moment we don't have specific API client libraries but any
[HTTP2 library][5] will work.

It's possible to use HTTP1 but highly not recommended, since it will introduce
huge overhead on connections.

## Roadmap
Our [Roadmap][6]

## License
[Apache 2.0][7]

[1]: http://martinfowler.com/eaaDev/EventSourcing.html
[2]: https://github.com/boltdb/bolt
[3]: https://glide.sh
[4]: https://golang.org/dl/
[5]: https://github.com/http2/http2-spec/wiki/Implementations
[6]: https://github.com/eventsourcedb/eventsourcedb/blob/master/ROADMAP.md
[7]: https://github.com/eventsourcedb/eventsourcedb/blob/master/LICENSE
