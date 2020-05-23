#### Build stage ####
FROM golang AS build

RUN wget -q -O nsq.tgz https://github.com/nsqio/nsq/releases/download/v1.0.0-compat/nsq-1.0.0-compat.linux-amd64.go1.8.tar.gz \
    && tar -xzf nsq.tgz \
    && cp nsq-1.0.0-compat.linux-amd64.go1.8/bin/nsqd /usr/local/bin

# Build the docproc applications
ARG TAGS="beanstalk nats nsq"
ENV SRC_DIR=/go/src/github.com/marcusva/docproc
ADD . $SRC_DIR

RUN (cd $SRC_DIR && CGO_ENABLED=0 make all test) || (echo "tests failed" && false)
RUN cd $SRC_DIR && PREFIX=/app make install

#### Image creation ####
FROM alpine

# curl is used for integration testing
RUN apk add --update curl

COPY --from=build /usr/local/bin/nsqd /usr/local/bin/nsqd
RUN ln -s /usr/local/bin/nsqd /bin/nsqd
COPY --from=build /app /app
