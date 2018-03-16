# Dockerfile
FROM node:6.7

WORKDIR /code/webroot

ADD ./webroot .

RUN npm --quiet install
RUN npm --quiet install babelify
RUN PATH=$PATH:node_modules/.bin npm --quiet run build

FROM nginx:1.11-alpine

WORKDIR /usr/share/nginx/html

COPY --from=0 /code/webroot /usr/share/nginx/html/
