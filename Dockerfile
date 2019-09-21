FROM golang:1.13 as builder
WORKDIR /build
COPY webhooky.go /build
RUN go build webhooky.go

FROM debian:buster-slim
WORKDIR /app
COPY --from=builder /build/webhooky .
ENTRYPOINT ["./webhooky"]
