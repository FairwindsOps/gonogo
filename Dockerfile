FROM alpine:3.19

LABEL org.opencontainers.image.authors="FairwindsOps, Inc." \
    org.opencontainers.image.vendor="FairwindsOps, Inc." \
    org.opencontainers.image.title="gonogo" \
    org.opencontainers.image.description="GoNoGo is a  utility to help users determine upgrade confidence around Kubernetes cluster addons." \
    org.opencontainers.image.documentation="https://gonogo.docs.fairwinds.com/" \
    org.opencontainers.image.source="https://github.com/FairwindsOps/gonogo" \
    org.opencontainers.image.url="https://github.com/FairwindsOps/gonogo" \
    org.opencontainers.image.licenses="Apache License 2.0"

USER nobody
COPY gonogo /

ENTRYPOINT ["/gonogo"]