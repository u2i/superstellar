FROM golang:1.6
MAINTAINER Michał Knapik <michal.knapik@u2i.com>

WORKDIR $GOPATH/src/superstellar/superstellar_utils

ADD . /go/src/superstellar

RUN go get
RUN go build
RUN go install

EXPOSE 8080

ENTRYPOINT /go/bin/superstellar_utils backend 100 50ms
