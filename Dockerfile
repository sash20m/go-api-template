FROM golang:1.21 AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY ./cmd/server/main.go ./cmd/server/main.go
COPY ./cmd/server/docs/ ./cmd/server/docs/
COPY ./pkg/ ./pkg/
COPY ./internal/ ./internal/
COPY ./config/ ./config/

RUN CGO_ENABLED=0 go build -o ./server ./cmd/server/main.go

# Stage 2
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=0 /app/server ./
COPY ./migrations/ ./migrations/
COPY .env.prod .env.prod
COPY .env .env
COPY ./logs/ ./logs/

EXPOSE 8080

CMD ["./server"]