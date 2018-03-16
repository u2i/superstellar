FROM node:6.7
MAINTAINER Michał Knapik <michal.knapik@u2i.com>

WORKDIR /code/webroot

ADD . /code

EXPOSE 8090

RUN npm install
ENTRYPOINT npm run dev
