# Use the official Golang image as the base image for building
FROM --platform=linux/amd64 golang:1.22 as builder

WORKDIR /build

# Install required packages including Kafka dependencies
RUN apt-get update && apt-get install -y \
    build-essential \
    librdkafka-dev \
    && rm -rf /var/lib/apt/lists/*

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code into the builder container
COPY . .

# Print contents of /build directory for debugging
RUN ls -l /build

# Build the Go app for Linux (amd64) with CGO enabled
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o ecommerce-service .

# Create a new stage for the final application image (based on Alpine Linux)
FROM alpine:3.18

# Install PostgreSQL client tools
RUN apk add --no-cache postgresql-client

# Set the working directory inside the container
WORKDIR /app

# Copy the built executable from the builder stage
COPY --from=builder /build/ecommerce-service ./ecommerce-service

# Copy configuration files
COPY --from=builder /build/app.env ./app.env

# Copy migration files
COPY --from=builder /build/db/migrations ./db/migrations

# Ensure the ecommerce-service binary is executable
RUN chmod +x ecommerce-service

# Print contents of /app directory for debugging
RUN ls -l /app

# Expose port 8080 if your application needs it
EXPOSE 8080

# Command to run your application
CMD ["./wait-for-db.sh", "db", "./ecommerce-service"]
