# TODOs

* add test setups for nats and beanstalk

* failover for:
  - out-queue|error-queue being unreachable
    - do not wait for Writer.Consume() { ... Publish() } to fail
  - in-queue being unreachable
  - changed directory permissions in docproc.fileinput

* nats: ensure that queue messages stay in the queue or get republished as
  long as no consumer could successfully process them

* document rules engine properly (inline and usage)

* add docproc.webinput for receiving input data ia HTTP
  * single messages
  * files to be routed into docproc.fileinput

* implement concurrency on message consumption for the queue implementations -
  consume multiple messages in parallel via goroutines per application
  (configurable)

* mail processor
* configurable headers for the HTTPSender

* configurable status codes to accept as default on HTTPSender?
* configurable http methods on HTTPSender?

* add Kafka as queue implementation
* add Rabbitmq as queue implementation
