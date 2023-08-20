# build application
FROM golang:1.21.0-alpine AS builder

ENV GOPRIVATE=github.com/vbyazilim

ARG GITHUB_ACCESS_TOKEN
ARG BUILD_INFORMATION

# hadolint ignore=DL3018
RUN apk add --update --no-cache git \
    && git config --global url.https://${GITHUB_ACCESS_TOKEN}@github.com/.insteadOf https://github.com/

WORKDIR /build
COPY ./go.mod /build/

# COPY ./go.mod ./go.sum /build/
# RUN go mod download

COPY . /build
RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -ldflags="-X 'github.com/vbyazilim/kvstore/src/releaseinfo.BuildInformation=${BUILD_INFORMATION}'" -o app ./cmd/server

# get certificates
FROM alpine:3.18.3 AS certs

# hadolint ignore=DL3018
RUN apk add --update --no-cache ca-certificates

FROM busybox:1.36
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /build/app /kvstoreapp

EXPOSE 8000
CMD ["/kvstoreapp"]
