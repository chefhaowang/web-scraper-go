# Stage 1: Build the Go application
FROM golang:1.23 AS builder

WORKDIR /app

# Copy Go modules and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build for Linux amd64
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o webscraper server.go webscraper.go

# Stage 2: Create a lightweight runtime container
FROM --platform=linux/amd64 ubuntu:22.04

WORKDIR /app

# Install required system dependencies (NO Chrome/ChromeDriver)
RUN apt-get update && apt-get install -y \
    curl \
    libnss3 \
    libx11-xcb1 \
    fonts-liberation \
    xdg-utils \
    && rm -rf /var/lib/apt/lists/*

# Copy the built Go binary
COPY --from=builder /app/webscraper /app/webscraper
RUN chmod +x /app/webscraper

# Expose the gRPC port (ChromeDriver runs externally)
EXPOSE 50051

# Run the web scraper (expects ChromeDriver to be running separately)
CMD ["/app/webscraper"]