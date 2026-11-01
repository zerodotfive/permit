FROM --platform=$BUILDPLATFORM node:24-trixie AS static

ADD ui /build
WORKDIR /build
RUN npm install; npm run build

FROM --platform=$BUILDPLATFORM golang:1.27-trixie AS build

ADD . /build
COPY --from=static /build/dist /build/ui/dist
WORKDIR /build

ARG TARGETOS
ARG TARGETARCH
ARG GITHUB_REF_NAME
RUN go mod download;  \
    CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -ldflags="-s -w -X github.com/zerodotfive/permit/internal/info.ApplicationVersion=${GITHUB_REF_NAME}" \
       -o bin/permit-server ./cmd/server

FROM scratch AS binary
COPY --from=build /build/bin/permit-server /permit-server

FROM debian:trixie
COPY --from=build /build/bin/permit-server /usr/local/bin/permit-server
RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates curl \
    && rm -rf /var/lib/apt/lists/* \
    && rm -rf /var/cache/apt/archives/*

USER www-data

CMD ["/usr/local/bin/permit-server"]
