FROM alpine:3.17.1

USER nobody
COPY  gonogo /

ENTRYPOINT ["/gonogo"]