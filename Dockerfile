# Build stage
FROM golang:1.25.5-alpine3.23 AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy dependency files and download
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build as a static binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gauth ./cmd/gauth

# Final stage
FROM alpine:3.23.2

# Basic requirements
RUN apk --no-cache add tzdata ca-certificates

WORKDIR /root/

# Copy binary
COPY --from=builder /app/gauth /usr/local/bin/gauth

# Setup standard config dir
RUN mkdir -p /root/.gauth

# Application setup
ENTRYPOINT ["gauth"]
