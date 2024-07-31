FROM --platform=linux/amd64 golang:1.22 as builder

WORKDIR /build

RUN apt-get update && apt-get install -y \
    build-essential \
    librdkafka-dev \
    && rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o ecommerce-service .

FROM --platform=linux/amd64 golang:1.22 as runner

# Install PostgreSQL client and other dependencies
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    postgresql-client \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=builder /build/ecommerce-service ./ecommerce-service
COPY --from=builder /build/app.env ./app.env
COPY --from=builder /build/db/migrations ./db/migrations
COPY --from=builder /build/wait-for-db.sh ./wait-for-db.sh

RUN chmod +x ecommerce-service wait-for-db.sh
EXPOSE 8080
CMD ["./wait-for-db.sh", "db", "./ecommerce-service"]
