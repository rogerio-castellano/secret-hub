# Stage 1: Build the binary
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o secret-hub .

# Stage 2: Minimal runtime image
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/secret-hub .

# Set default entrypoint
ENTRYPOINT ["/app/secret-hub"]
