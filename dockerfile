# Use the official Go image as the base image
FROM golang:1.22.5 AS builder

# Create a directory inside the container to store all our application and then make it the working directory
WORKDIR /app

# Copy the go.mod and go.sum files to the working directory
ADD go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code from the current directory to the working Directory inside the container
ADD . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/site/

# Start a new stage from scratch
FROM alpine:3.14

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main /app/main

# Expose port 8080 to the outside world
EXPOSE 8080

RUN apk add --no-cache bash curl && curl -1sLf \
'https://dl.cloudsmith.io/public/infisical/infisical-cli/setup.alpine.sh' | bash \
&& apk add infisical

# Command to run the executable
CMD ["infisical", "run", "--projectId", "5b6ac07a-a07e-48a6-aa70-60bd992cd693", "--", "/app/main"]
