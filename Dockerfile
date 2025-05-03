FROM golang:1.24.1-alpine AS builder

WORKDIR /app

COPY go.* ./

RUN go mod download

COPY . ./
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

WORKDIR /app

EXPOSE 8080
EXPOSE 9100

CMD ["main"]