# Stage 1: Build the Go binary using Alpine
FROM golang:1.24.1-alpine AS builder

# Install build tools
RUN apk add --no-cache git

# Set up working directory
WORKDIR /app

# Copy and build
COPY . .
RUN go build -o server .

# Stage 2: Minimal runtime image
FROM alpine:3.21

# Create non-root user (optional but good practice)
RUN adduser -D -g '' appuser

WORKDIR /app
COPY --from=builder /app/server .

# Run as non-root user
USER appuser

# Expose port
EXPOSE 8880

# Run the app
ENTRYPOINT ["./server"]

