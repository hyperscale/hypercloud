# Manager Builder
FROM node:10.7-alpine as manager-builder
WORKDIR /build/
ARG VERSION
ARG VCS_URL
ARG VCS_REF
ARG BUILD_DATE
RUN apk add --update --no-cache python python-dev
COPY ./ .
WORKDIR /build/cmd/hypercloud-manager
RUN mkdir .npm-global
RUN npm config set prefix '/build/cmd/hypercloud-manager/.npm-global'
RUN export PATH=/build/cmd/hypercloud-manager/.npm-global/bin:$PATH
RUN npm install
RUN npm run build

# Application
FROM node:10.7-alpine
ARG VERSION
ARG VCS_URL
ARG VCS_REF
ARG BUILD_DATE
ENV PORT 8080
ENV NODE_ENV production
ENV NPM_CONFIG_LOGLEVEL warn
RUN apk add --update --no-cache ca-certificates curl && rm -rf /var/cache/apk/*
WORKDIR /opt/hypercloud-manager/
COPY --from=manager-builder /build/cmd/hypercloud-manager/bin ./bin
COPY --from=manager-builder /build/cmd/hypercloud-manager/dist ./dist
COPY --from=manager-builder /build/cmd/hypercloud-manager/config ./config
HEALTHCHECK --interval=5s --timeout=2s CMD curl -f http://localhost:${PORT}/health > /dev/null 2>&1 || exit 1
EXPOSE ${PORT}
CMD [ "node", "bin/server.js" ]

# Metadata
LABEL org.label-schema.vendor="Hyperscale" \
      org.label-schema.url="https://github.com/hyperscale" \
      org.label-schema.name="Hypercloud" \
      org.label-schema.description="A modern cloud platform" \
      org.label-schema.version="v${VERSION}" \
      org.label-schema.vcs-url=${VCS_URL} \
      org.label-schema.vcs-ref=${VCS_REF} \
      org.label-schema.build-date=${BUILD_DATE} \
      org.label-schema.docker.schema-version="1.0"
