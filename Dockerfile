FROM golang:1.25.10-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download
RUN CGO_ENABLED=0 go build -o /app/cmd/api /app/cmd/api

FROM alpine:latest

WORKDIR /root

ENV DEPLOY=true

COPY --from=builder /app/config/env/.env /.env
COPY --from=builder /app/cmd/api .

CMD ["./api"]