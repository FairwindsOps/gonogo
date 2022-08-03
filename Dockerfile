FROM alpine:3.16.1

USER nobody
COPY  gonogo /

ENTRYPOINT ["/gonogo"]