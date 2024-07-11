# syntax=docker/dockerfile:1
FROM --platform=$BUILDPLATFORM golang:alpine AS build
WORKDIR /src
COPY . .
ARG TARGETOS TARGETARCH
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH CGO_ENABLED=0 go build -a -o GeoIPGS  ./main.go

FROM --platform=${TARGETPLATFORM} alpine:latest
WORKDIR /app
COPY --from=build /src/GeoIPGS /app/GeoIPGS
RUN set -ex \
    && chmod 544 /app/GeoIPGS



ENV TZ=Asia/Shanghai
CMD [ "/app/GeoIPGS" ]
