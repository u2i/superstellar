FROM golang:1.6
MAINTAINER Michał Knapik <michal.knapik@u2i.com>

COPY docker/utils/entrypoint.sh /
RUN chmod +x /entrypoint.sh

VOLUME "$GOPATH/src/superstellar"

WORKDIR $GOPATH/src/superstellar/superstellar_utils

EXPOSE 8080

ENTRYPOINT /entrypoint.sh
