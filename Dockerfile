FROM golang:latest AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY ./cmd/api/main.go ./cmd/api/main.go
COPY ./cmd/migration/main.go ./cmd/migration/main.go
COPY ./internal/ ./internal/
COPY ./config/ ./config/
COPY ./scripts/ ./scripts/
COPY .env .env

RUN CGO_ENABLED=0 go build -o ./api ./cmd/api/main.go

# Stage 2
FROM alpine:latest

RUN apk --no-cache add ca-certificates curl vim nano

WORKDIR /root/

COPY --from=build /app/api ./
COPY .env .env

EXPOSE 8080

CMD ["./api"]