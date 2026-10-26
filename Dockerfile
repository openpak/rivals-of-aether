# syntax=docker/dockerfile:1
# Rivals of Aether's first-party web side: match telemetry (stats.rivalsofaether.com)
# and the in-game news CDN (www.rivalsofaether.com). Plain HTTP on 23570 behind the
# TLS terminator; health on 23571. Netplay (NEX 4.6.4) is served by switch-nex.
FROM --platform=$BUILDPLATFORM golang:1.26-alpine AS build
ARG TARGETOS TARGETARCH
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH GOWORK=off go build -trimpath -buildvcs=false -ldflags="-s -w" -o /out/rivals-of-aetherd ./cmd/rivals-of-aetherd

FROM gcr.io/distroless/static
COPY --from=build /out/rivals-of-aetherd /rivals-of-aetherd
USER nonroot:nonroot
EXPOSE 23570 23571
ENTRYPOINT ["/rivals-of-aetherd"]
