FROM golang:1.23 AS build

ENV GOTOOLCHAIN=auto

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o app .

FROM alpine:3.19
WORKDIR /root/
COPY --from=build /app/app .
ENTRYPOINT ["./app", "server"]
