# TODOs

* add test setups for nats and beanstalk

* ensure that queue messages stay in the queue or get republished as long as
  no consumer could successfully process them (queue bindings)
  - standard behaviour for NSQ

* necessary for the queue consumptions:
    * thread-safety on Consumer implementations (queue.Writer)
    * thread-safety on Processor implementations (queue.Writer.ProcChain)

* document rules engine properly (inline and usage)

* add docproc.webinput for receiving input data ia HTTP
  * single messages
  * files to be routed into docproc.fileinput

* implement concurrency on message consumption for the queue implementations -
  consume multiple messages in parallel via goroutines per application
  (configurable)

* Support filters/identifiers for different renderers in a single renderer
  instance:
    * if message meets condition X, use Renderer A only
    * if message meets condition Y, use Renderer B only

* add XSL-FO renderer using a saxon/xsltproc invocation for PS, PDF, SVG, ...
* message filters for multiple output queues
* mail handler for output application
* http handler for SOAP and REST receivers

* add Kafka as queue implementation
* add Rabbitmq as queue implementation
