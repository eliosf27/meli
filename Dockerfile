# build stage
FROM golang:latest as builder

WORKDIR /go/src/meli
COPY . .

RUN go install -mod=vendor -v ./...

FROM debian:stable-slim

RUN apt-get update; apt-get install -y wget
RUN mkdir -p /usr/local/share/ca-certificates/cacert.org
RUN wget -P /usr/local/share/ca-certificates/cacert.org http://www.cacert.org/certs/root.crt http://www.cacert.org/certs/class3.crt
RUN update-ca-certificates

COPY --from=builder /go/bin/ /app/
COPY config/* /app/config/

ENV ENVIRONMENT=develop

WORKDIR /app/

USER 1000:1000
ENTRYPOINT ["/app/httpserver" ]
