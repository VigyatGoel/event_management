# Use the official Golang image for building
FROM golang:1.24.2 AS build

WORKDIR /src

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the application
RUN go build cmd/main.go

# Use a minimal image for running the application
FROM alpine:latest

WORKDIR /app

# Copy the built application
COPY --from=build /src main

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./main"]
