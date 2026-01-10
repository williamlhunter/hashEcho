FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY *.go ./
RUN GO111MODULE=off go build -o echo-server .

from alpine:latest
WORKDIR /root/
COPY --from=builder /app/echo-server .
Expose 8080
CMD ["./echo-server"]
