# Building docproc

docproc requires [Go](http://golang.org) to be built. Aside from that no other
build tools are necessary.

To build the individual applications, you can stick to Go's way of building and
installing applications, e.g.

    # go build -v github.com/marcusva/docproc/docproc.fileinput

or

    # go install -v github.com/marcusva/docproc/docproc.fileinput


## Queue Support

docproc relies on a message queue implementation. It currently supports the
following message queues:

* beanstalk - http://kr.github.io/beanstalkd/
* NSQ - http://nsq.io/
* NATS - https://nats.io/

Use the appropriate build tags to build docproc with support for one or more
of those:

* beanstalk: `beanstalk`
* NSQ: `nsq`
* NATS: `nats`

They can be added at build time via Go's `-tags` parameter:

    # go build -tags nsq,beanstalk -v github.com/marcusva/docproc/docproc.preproc

Other message queues can be easily supported by implementing the `ReadQueue`
and `WriteQueue` interfaces of the `docproc/common/queue` package.

# Docker Builds

You can find example setups for Docker in the test/dockerfiles directory.
