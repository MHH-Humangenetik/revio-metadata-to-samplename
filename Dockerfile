# syntax=docker/dockerfile:1
FROM --platform=$BUILDPLATFORM golang:1.25.1 AS build

ARG TARGETOS
ARG TARGETARCH

COPY . /app
WORKDIR /app
RUN mkdir -p build && \
    go mod download && \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -a -ldflags '-extldflags "-static"' -v -o /app/build/revio-metadata-to-samplename

FROM scratch
COPY --from=build /app/build/revio-metadata-to-samplename /revio-metadata-to-samplename
WORKDIR /
ENTRYPOINT ["/revio-metadata-to-samplename"]
