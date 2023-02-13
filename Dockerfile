FROM alpine:3.17.2

USER nobody
COPY  gonogo /

ENTRYPOINT ["/gonogo"]