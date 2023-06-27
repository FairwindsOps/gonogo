FROM alpine:3.18

LABEL org.opencontainers.image.authors="FairwindsOps, Inc." \
    org.opencontainers.image.vendor="FairwindsOps, Inc." \
    org.opencontainers.image.title="gonogo" \
    org.opencontainers.image.description="GoNoGo is a  utility to help users determine upgrade confidence around Kubernetes cluster addons." \
    org.opencontainers.image.documentation="https://gonogo.docs.fairwinds.com/" \
    org.opencontainers.image.source="https://github.com/FairwindsOps/gonogo" \
    org.opencontainers.image.url="https://github.com/FairwindsOps/gonogo" \
    org.opencontainers.image.licenses="Apache License 2.0"

WORKDIR /usr/local/bin

RUN apk -U upgrade
RUN apk --no-cache add ca-certificates

RUN addgroup -S gonogo && adduser -u 1200 -S gonogo -G gonogo
USER 1200
COPY gonogo .

WORKDIR /opt/app

ENTRYPOINT ["gonogo"]