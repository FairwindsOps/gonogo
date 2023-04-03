FROM alpine:3.17.3

USER nobody
COPY  gonogo /

ENTRYPOINT ["/gonogo"]