# TODOs

* failover for:
  - out-queue|error-queue being unreachable
    - do not wait for Writer.Consume() { ... Publish() } to fail
  - in-queue being unreachable
  - changed directory permissions in docproc.fileinput

* document rules engine properly (inline and usage)

* mail processor

* configurable status codes to accept as default on HTTPSender?
* configurable http methods on HTTPSender?

* add Kafka as queue implementation
* add Rabbitmq as queue implementation
