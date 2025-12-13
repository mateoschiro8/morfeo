FROM golang:1.24.9-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o morfeo .

FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache \
    ca-certificates \
    gcc \
    musl-dev

COPY --from=build /usr/local/go /usr/local/go

ENV PATH="/usr/local/go/bin:${PATH}"

COPY --from=build /app/morfeo .

EXPOSE 8000

ENV GIN_MODE=release
ENV CGO_ENABLED=0
ENV GO111MODULE=off

RUN mkdir -p /app/output /app/input /tmp
RUN chmod -R 777 /app/output /app/input /tmp

