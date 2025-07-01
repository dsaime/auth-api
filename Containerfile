# Build stage
FROM golang:1.24 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o auth-api cmd/auth-api/main.go

# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/auth-api ./
CMD ["./auth-api"]
ENTRYPOINT ["./auth-api"]