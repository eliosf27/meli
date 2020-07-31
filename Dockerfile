# build stage
FROM golang:latest as builder

WORKDIR /go/src/meli
COPY . .

RUN go install -mod=vendor -v ./...

FROM debian:stable-slim

COPY --from=builder /go/bin/ /app/
COPY config/* /app/config/

ENV ENVIRONMENT=production

WORKDIR /app/

USER 1000:1000
ENTRYPOINT ["/app/httpserver" ]
