FROM golang:1.24.9-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o app .

FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates
COPY --from=build /app/app .

EXPOSE 8000

ENV GIN_MODE=release

ENTRYPOINT ["./app", "server"]
