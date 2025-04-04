# Build Stage
FROM golang:1.21-alpine AS builder
RUN apk add --no-cache build-base gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Final Stage (no need to copy binary)
FROM golang:1.21-alpine

WORKDIR /app

# Copy application code from the builder stage
COPY --from=builder /app .

# Start the application using go run
CMD ["go", "run", "./api"]
