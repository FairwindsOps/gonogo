FROM alpine:3.18.0

USER nobody
COPY  gonogo /

ENTRYPOINT ["/gonogo"]