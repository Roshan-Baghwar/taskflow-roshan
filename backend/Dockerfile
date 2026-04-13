FROM golang:1.25-alpine AS builder
WORKDIR /app

# Copy dependency files first
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o taskflow ./cmd/server

# Runtime stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/taskflow .
COPY --from=builder /app/migrations ./migrations
EXPOSE 8080
CMD ["./taskflow"]