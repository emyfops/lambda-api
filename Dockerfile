FROM golang:1.23.1-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.* ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code
COPY . ./

# Build the application
RUN go build -o main .

# Use the official Alpine image for a lean production container
# https://hub.docker.com/_/alpine
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM alpine:3.20.3
RUN set +x \
    apt-get update \
    apt-get install -y ca-certificates \
    rm -rf /var/lib/apt/lists/*

# Copy the pre-built binary file from the previous stage
COPY --from=builder /app/main /app/

# Expose port 8080 to the outside world
EXPOSE 8080

# Run the server
CMD ["/app/main"]