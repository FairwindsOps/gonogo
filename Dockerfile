FROM alpine:3.16.3

USER nobody
COPY  gonogo /

ENTRYPOINT ["/gonogo"]