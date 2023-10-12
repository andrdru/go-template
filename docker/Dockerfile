FROM golang:1.21-alpine

RUN go install github.com/rubenv/sql-migrate/...@v1.5.1 &&\
    go install github.com/a8m/envsubst/cmd/envsubst@latest
