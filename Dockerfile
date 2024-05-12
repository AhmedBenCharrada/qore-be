# Use the official Golang image as base
FROM golang:latest AS builder

# Set the current working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o app cmd/main.go

# Use a minimal base image to run the application
FROM alpine:latest

# Set the current working directory inside the container
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/app .

# Expose port 8080 to the outside world
EXPOSE 8080

# Set environment variable for DB_URL
ARG DB_URL
ENV DB_URL=$DB_URL

# Command to run the executable
CMD ["./app"]
