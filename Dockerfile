FROM golang:latest AS builder

RUN mkdir /build
ADD . /build
WORKDIR /build

ENV TRUST_URL=

RUN go build -o bin/main .

FROM debian:latest
COPY --from=builder /build/bin /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 5000

CMD ["/app/main"]