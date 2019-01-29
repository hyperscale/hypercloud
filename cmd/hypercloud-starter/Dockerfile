# Server Builder
FROM golang:1.11-alpine as starter-builder
WORKDIR /go/src/github.com/hyperscale/hypercloud/
ARG VERSION
ARG VCS_URL
ARG VCS_REF
ARG BUILD_DATE
ENV GO111MODULE on
ENV CGO_ENABLED=0
RUN echo http://mirror.yandex.ru/mirrors/alpine/v3.8/main > /etc/apk/repositories; \
    echo http://mirror.yandex.ru/mirrors/alpine/v3.8/community >> /etc/apk/repositories
RUN echo "ipv6" >> /etc/modules
RUN apk upgrade --update-cache --available && \
    apk add --update --no-cache ca-certificates curl git make && \
    rm -rf /var/cache/apk/*
RUN GO111MODULE=off go get -u github.com/go-bindata/go-bindata/...
COPY ./ .
RUN make build/hypercloud-starter

# Application
FROM alpine:3.8
ARG VERSION
ARG VCS_URL
ARG VCS_REF
ARG BUILD_DATE
ENV PORT 8080
ENV DEBUG true
RUN apk add --update --no-cache ca-certificates curl && rm -rf /var/cache/apk/*
WORKDIR /opt/hypercloud/
COPY --from=starter-builder /go/src/github.com/hyperscale/hypercloud/build/hypercloud-starter .
HEALTHCHECK --interval=5s --timeout=2s CMD curl -f http://localhost:${PORT}/health > /dev/null 2>&1 || exit 1
EXPOSE ${PORT}
VOLUME /var/lib/hypercloud/
ENTRYPOINT ["/opt/hypercloud/hypercloud-starter"]

# Metadata
LABEL org.label-schema.vendor="Hyperscale" \
      org.label-schema.url="https://github.com/hyperscale" \
      org.label-schema.name="Hypercloud Starter" \
      org.label-schema.description="A modern cloud platform" \
      org.label-schema.version="v${VERSION}" \
      org.label-schema.vcs-url=${VCS_URL} \
      org.label-schema.vcs-ref=${VCS_REF} \
      org.label-schema.build-date=${BUILD_DATE} \
      org.label-schema.docker.schema-version="1.0"
