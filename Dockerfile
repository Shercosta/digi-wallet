# Use a Go base image with version 1.23 or higher
FROM golang:1.23-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of app files
COPY . .

# Build the Go app
# Explicitly set GOARCH for cross-compilation if needed.
# For most Docker setups, linux/amd64 is the common target for alpine:latest.
# If your final image (alpine:latest) is ARM64 (e.g., you're on M1/M2 and docker pulls arm64 alpine), then use GOARCH=arm64.
# Let's assume common case is building for amd64 for now, as that's what many public images target.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o digi-wallet .

# Final stage: Smaller runtime image
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/digi-wallet .

# Add this temporary debug line if the issue persists after the above steps
# RUN ls -lh /app/digi-wallet

# Expose the port your Go application listens on
EXPOSE 8080

# Run the Go app
CMD ["./digi-wallet"]