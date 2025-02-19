# Stage 1: Build the Go application
FROM golang:1.23 AS builder

WORKDIR /app

# Copy Go modules and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# ✅ Build for Linux amd64 with static linking
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o webscraper server.go webscraper.go

# Stage 2: Create a lightweight runtime container
FROM --platform=linux/amd64 ubuntu:22.04

WORKDIR /app

# ✅ Copy the built Go binary
COPY --from=builder /app/webscraper /app/webscraper
RUN chmod +x /app/webscraper

# Expose gRPC port (ChromeDriver runs externally)
EXPOSE 50051

# ✅ Run the web scraper
CMD ["/app/webscraper"]