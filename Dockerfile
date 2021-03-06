FROM golang as base

WORKDIR /app

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o app ./cmd/web

FROM alpine:latest as certs

RUN apk --update add ca-certificates

FROM scratch as app

COPY entries.json /
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=base app /

ENV ENTRIES_PATH /entries.json

ENTRYPOINT ["/app"]
