# Use the official Go image as the base image
FROM golang:1.22.5 AS builder

# Create a directory inside the container to store all our application and then make it the working directory
WORKDIR /app

# Copy the go.mod and go.sum files to the working directory
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code from the current directory to the working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o main ./cmd/site/

# Start a new stage from scratch
FROM alpine:3.14

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main /app/main

# Command to run the executable
CMD ["/app/main"]
