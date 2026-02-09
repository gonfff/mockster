FROM node:22-alpine AS ui-build

WORKDIR /ui

COPY package.json package-lock.json tsconfig.json ./
COPY app/static/ts ./app/static/ts

RUN npm ci
RUN npm run build:ui


FROM golang:1.25-alpine AS build

ENV CGO_ENABLED=0
ENV GOOS=linux

WORKDIR /build

COPY . .

COPY --from=ui-build /ui/app/static/script.js ./app/static/script.js

RUN go build -o ./mockster app/main.go

## App
FROM alpine:latest AS app

WORKDIR /app

COPY --from=build /build/mockster .

ENV STATIC_PATH=/app/static
COPY app/static ./static

EXPOSE 8080

CMD ["./mockster"]
