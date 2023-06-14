FROM alpine:3.18

USER nobody
COPY  gonogo /

ENTRYPOINT ["/gonogo"]