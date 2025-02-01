package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "web-scraper-go/scraperpb"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedNewsScraperServer
}

func (s *server) GetTopNews(ctx context.Context, req *pb.Empty) (*pb.NewsResponse, error) {
    log.Println("📥 Received gRPC request to scrape news")

    articles := ScrapeTopNews()

    log.Println("✅ Successfully scraped news articles")
    return &pb.NewsResponse{Articles: articles}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterNewsScraperServer(grpcServer, &server{})

	fmt.Println("🚀 gRPC server running on port 50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}