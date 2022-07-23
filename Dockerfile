FROM golang:1.19rc2-alpine AS build

ADD . /build
WORKDIR /build

RUN go mod download && go build -o start_resolver main.go

FROM alpine:3.5 as base

COPY --from=build /build/start_resolver .
CMD ["./start_resolver"]
