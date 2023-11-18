FROM golang:1.21 as builder
WORKDIR /build
COPY webhooky.go /build
RUN go build webhooky.go

FROM ubuntu:22.04
WORKDIR /app
COPY --from=builder /build/webhooky .
ENTRYPOINT ["./webhooky"]
