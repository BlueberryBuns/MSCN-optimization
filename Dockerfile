FROM golang:1.19-alpine3.15 AS builder

WORKDIR /app

COPY go.mod ./
COPY *.go ./


RUN CGO_ENABLED=0 GOOS=linux go build -a -o main.out .

FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main.out ./
CMD ["./main.out"]
