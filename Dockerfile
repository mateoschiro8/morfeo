FROM golang:1.24.9-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o morfeo .

FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates
COPY --from=build /app/morfeo .

EXPOSE 8000

ENV GIN_MODE=release

ENTRYPOINT ["./morfeo"]

CMD ["server"]