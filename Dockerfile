FROM alpine:3.16.0

USER nobody
COPY  gonogo /

ENTRYPOINT ["/gonogo"]