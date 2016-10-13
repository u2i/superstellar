FROM debian:jessie

EXPOSE 8080

COPY superstellar /

CMD /superstellar
