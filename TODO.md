# TODOs

* add test setups for nats and beanstalk

* ensure that queue messages stay in the queue or get republished as long as
  no consumer could successfully process them (queue bindings)
  - standard behaviour for NSQ

* document rules engine properly (inline and usage)

* add docproc.webinput for receiving input data ia HTTP
  * single messages
  * files to be routed into docproc.fileinput

* implement concurrency on message consumption for the queue implementations -
  consume multiple messages in parallel via goroutines per application
  (configurable)

* mail processor
* http processor for SOAP and REST receivers

* add Kafka as queue implementation
* add Rabbitmq as queue implementation
