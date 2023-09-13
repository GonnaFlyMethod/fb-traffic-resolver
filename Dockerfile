FROM golang:1.18-alpine AS build

ADD . /build
WORKDIR /build

RUN go mod download && go build -o start_resolver main.go

FROM alpine:3.16 as base

COPY --from=build /build/start_resolver .
CMD ["./start_resolver"]
