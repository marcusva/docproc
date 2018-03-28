#### Build stage ####
FROM golang AS build

# Build nsq
RUN go get -v github.com/nsqio/nsq/... \
    && wget -O /bin/dep https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 \
    && chmod +x /bin/dep \
    && cd /go/src/github.com/nsqio/nsq \
    && /bin/dep ensure \
    && ./test.sh \
    && CGO_ENABLED=0 make DESTDIR=/opt PREFIX=/nsq GOFLAGS='-ldflags="-s -w"' install;

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

COPY --from=build /opt/nsq/bin/nsqd /usr/local/bin/nsqd
RUN ln -s /usr/local/bin/nsqd /bin/nsqd
COPY --from=build /app /app
