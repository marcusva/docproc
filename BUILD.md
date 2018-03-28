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

# Build for Docker

1. Create the base image via docker
2. Build the docproc images via docker-compose
3. Run everything via docker-compose

## Base Image

The docproc base image contains all docproc applications as well as a
[nsqd](http://nsq.io) binary to get docproc up and running for testing with a
message queue system.

Create the base image with the following instructions:

    # docker build -t docproc/base .

The base image is now registered as `docproc/base`.

## Build docproc Images

Create all docproc images with the following instruction:

    # docker-compose build

This creates the following set of docproc images:

* `docproc/fileinput`
* `docproc/preproc`
* `docproc/renderer`
* `docproc/postproc`
* `docproc/output`

Each image can be run individually. Each image runs a local nsqd server to be
used by the individual docproc executable.

## Run everything

All services, including an nsqd, nsqlookupd and nsqadmin instance can be run via

    # docker-compose up

TODO: document ports and directories properly.

### Testing Service concurrency

Scale individual applications via the `--scale <service>=<num>` flag:

    # docker-compose up --scale docproc.preproc=3

will spawn three instances of the `docproc.preproc` service configured in the
docker-compose configuration.