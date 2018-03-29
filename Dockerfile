#### Build stage ####
FROM golang AS build

RUN wget -q -O /bin/dep https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 \
    && chmod +x /bin/dep
RUN wget -q -O nsq.tgz https://github.com/nsqio/nsq/releases/download/v1.0.0-compat/nsq-1.0.0-compat.linux-amd64.go1.8.tar.gz \
    && tar -xzf nsq.tgz \
    && cp nsq-1.0.0-compat.linux-amd64.go1.8/bin/nsqd /usr/local/bin

# Build the docproc applications
ARG TAGS="beanstalk nats nsq"
ENV SRC_DIR=/go/src/github.com/marcusva/docproc
ADD . $SRC_DIR

RUN cd $SRC_DIR \
    && /bin/dep ensure \
    && go build -tags "$TAGS" -v -o /app/docproc.fileinput ./docproc.fileinput \
    && go build -tags "$TAGS" -v -o /app/docproc.proc ./docproc.proc

RUN (cd $SRC_DIR && go test ./...) || (echo "tests failed" && false)

#### Image creation ####
FROM golang

COPY --from=build /usr/local/bin/nsqd /usr/local/bin/nsqd
RUN ln -s /usr/local/bin/nsqd /bin/nsqd
COPY --from=build /app /app
