# Builder
FROM golang:1.8-alpine as builder
WORKDIR /go/src/github.com/hyperscale/hyperpaas/
RUN apk add --update --no-cache ca-certificates curl git make && rm -rf /var/cache/apk/*
RUN curl https://glide.sh/get | sh
COPY ./ .
RUN make deps
RUN make build
RUN chmod +x hyperpaas

# Application
FROM alpine:latest
LABEL maintainer "Axel Etcheverry <axel@etcheverry.biz>"
ENV PORT 8080
RUN apk add --update --no-cache ca-certificates curl && rm -rf /var/cache/apk/*
WORKDIR /root/
COPY --from=builder /go/src/github.com/hyperscale/hyperpaas/hyperpaas .
HEALTHCHECK --interval=5s --timeout=2s CMD curl -f http://localhost:${PORT}/health > /dev/null 2>&1 || exit 1
EXPOSE ${PORT}
VOLUME /var/lib/hyperpaas
CMD [ "./hyperpaas" ]
