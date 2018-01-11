# Server Builder
FROM golang:1.9-alpine as installer-builder
WORKDIR /go/src/github.com/hyperscale/hyperpaas/
ARG VERSION
ARG VCS_URL
ARG VCS_REF
ARG BUILD_DATE
RUN echo http://mirror.yandex.ru/mirrors/alpine/v3.6/main > /etc/apk/repositories; \
    echo http://mirror.yandex.ru/mirrors/alpine/v3.6/community >> /etc/apk/repositories
RUN echo "ipv6" >> /etc/modules
RUN apk upgrade --update-cache --available && \
    apk add --update --no-cache ca-certificates curl git make && \
    rm -rf /var/cache/apk/*
RUN go get github.com/golang/dep/cmd/dep
COPY ./ .
RUN dep ensure -vendor-only && \
    CGO_ENABLED=0 go build -ldflags "-X github.com/hyperscale/hyperpaas/version.Tag=${VERSION} -X github.com/hyperscale/hyperpaas/version.Revision=${VCS_REF} -X github.com/hyperscale/hyperpaas/version.BuildAt=${BUILD_DATE}" ./cmd/hyperpaas-installer/

# Application
FROM docker:17.12.0-ce-dind
ARG VERSION
ARG VCS_URL
ARG VCS_REF
ARG BUILD_DATE
ENV PORT 8080
ENV DEBUG true
WORKDIR /opt/hyperpaas/
COPY --from=installer-builder /go/src/github.com/hyperscale/hyperpaas/hyperpass-installer .
EXPOSE ${PORT}
VOLUME /var/lib/hyperpaas
ENTRYPOINT ["/usr/local/bin/dockerd-entrypoint.sh", "/opt/hyperpaas/hyperpaas-installer"]

# Metadata
LABEL org.label-schema.vendor="Hyperscale" \
      org.label-schema.url="https://github.com/hyperscale" \
      org.label-schema.name="HyperPaaS" \
      org.label-schema.description="A modern PaaS" \
      org.label-schema.version="v${VERSION}" \
      org.label-schema.vcs-url=${VCS_URL} \
      org.label-schema.vcs-ref=${VCS_REF} \
      org.label-schema.build-date=${BUILD_DATE} \
      org.label-schema.docker.schema-version="1.0"
