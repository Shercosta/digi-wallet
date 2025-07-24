# Use a Go base image with version 1.23 or higher
FROM golang:1.23-alpine AS builder

# Set the working direcctory inside the container
WORKDIR /app

# Copy go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of app files
COPY . .

# Build the Go app
# CGO_ENABLED=0 is important for static binaries
# good for smaller images with 'scratch' or 'alpine' base in the final stage
# Force rebuild comment (remove this line after debugging)
RUN CGO_ENABLED=0 GOOS=linux go build -o digi-wallet .

# Final stage: Smaller runtime image
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/digi-wallet .

# --- ADD THIS LINE TEMPORARILY FOR DEBUGGING ---
RUN ls -lh /app/digi-wallet

# Expose the port your Go application listens on
EXPOSE 8080

# Run the Go app
CMD ["./digi-wallet"]