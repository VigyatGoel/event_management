# backend/Dockerfile

# Stage 1: Build the Go binary
FROM golang:1.24.2 AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/main.go

# Stage 2: Final image
FROM alpine:latest

# Create working directory
WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /src/main .

# Make sure it is executable
RUN chmod +x ./main

# Set environment variable for allowed frontend URL
ENV FRONTEND_URL=http://event-frontend:3000

# Expose the port your Go app runs on
EXPOSE 8080

# Start the app
CMD ["./main"]
