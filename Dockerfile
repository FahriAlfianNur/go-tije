FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY . .
RUN go mod tidy

RUN go build -o /bin/api cmd/api/main.go
RUN go build -o /bin/subscriber cmd/subscriber/main.go
RUN go build -o /bin/publisher cmd/publisher/main.go
RUN go build -o /bin/worker cmd/worker/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /bin/api /app/api
COPY --from=builder /bin/subscriber /app/subscriber
COPY --from=builder /bin/publisher /app/publisher
COPY --from=builder /bin/worker /app/worker
COPY .env.example /app/.env

EXPOSE 8080

CMD ["/app/api"]