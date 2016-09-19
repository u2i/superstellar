FROM node:6.7
MAINTAINER Michał Knapik <michal.knapik@u2i.com>

COPY docker/frontend/entrypoint.sh /
RUN chmod +x /entrypoint.sh

VOLUME /code

WORKDIR /code/webroot

EXPOSE 8090

ENTRYPOINT /entrypoint.sh
