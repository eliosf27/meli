# build stage
FROM golang:1.14 as builder

ENV GO111MODULE=on
WORKDIR /app
COPY . .
RUN go build

FROM debian:stable-slim

ENV ENVIRONMENT=prod
COPY --from=builder /app /app
WORKDIR /app

ENTRYPOINT ["/app/meli"]
