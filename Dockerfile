FROM golang:1.21-alpine as build

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOFLAGS=-mod=vendor

WORKDIR /build

COPY . .

RUN go build -o ./mockster app/main.go

## App
FROM alpine:latest as app

WORKDIR /app

COPY --from=build /build/mockster .

ENV STATIC_PATH=/app/static
COPY app/static ./static

EXPOSE 8080

CMD ["./mockster"]