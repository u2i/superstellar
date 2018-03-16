FROM golang:1.6
MAINTAINER Michał Knapik <michal.knapik@u2i.com>

WORKDIR $GOPATH/src/superstellar

ADD . /go/src/superstellar

RUN go get superstellar
RUN go build superstellar
RUN go install superstellar

FROM debian:jessie

EXPOSE 8080

COPY --from=0 /go/bin/superstellar /

CMD /superstellar
