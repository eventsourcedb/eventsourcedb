# Roadmap

## 0.1.0
* First release
* HTTP2 JSON API
* boltdb as a storage
* pub/sub

## 0.2.0
* scan and subscribe to all events
* consume events by type
* consume events by stream
* consume last N events
* consume events starting from ID
* consume events by time ranges

## 0.3.0
* archive retired streams to AWS S3 / OpenStack Swift
* query archived streams (without restore)
* restore archived streams
* replicate data

## 0.4.0
* build client libraries for API
* parallelize querying using replicas
