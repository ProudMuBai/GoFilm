FROM golang:tip-alpine3.23 AS builder

RUN apk add --no-cache ca-certificates git build-base



WORKDIR /build

COPY ./server/go.mod ./server/go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download


COPY ./server .


RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /build/gofilm ./main.go


FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/gofilm /gofilm

ENTRYPOINT ["/gofilm"]
