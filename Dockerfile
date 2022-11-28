FROM alpine:3.17.0

USER nobody
COPY  gonogo /

ENTRYPOINT ["/gonogo"]