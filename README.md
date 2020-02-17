# Flame

Torrent client for DashoTV

[![Build Status](https://travis-ci.org/dashotv/flame.svg?branch=master)](https://travis-ci.org/dashotv/flame)

### make [test]

Default target, runs a simple test. If running locally, will load .env if it exists.

```
# .env
export FLAME_URL="<url to utorrent gui>"
```

### make server

Run an instance of the server:

* Sends nats message with flame.Response every second
* Caches response every second
* Web service responds with cached message

### make receiver

Run a receiver:

* Subscribes to nats message and prints flame.Response
