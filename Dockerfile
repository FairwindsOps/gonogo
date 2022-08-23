FROM alpine:3.16.2

USER nobody
COPY  gonogo /

ENTRYPOINT ["/gonogo"]