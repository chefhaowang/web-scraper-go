package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "web-scraper-go/scraperpb"

	"google.golang.org/grpc"
)

func main() {
	// Connect to gRPC server
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewNewsScraperClient(conn)

	// Set email manually
	email := "1242107568@qq.com" // Replace with the target email

	// Set gRPC context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*180)
	defer cancel()

	// Call gRPC GetTopNews
	resp, err := client.GetTopNews(ctx, &pb.EmailRequest{Email: email})
	if err != nil {
		log.Fatalf("Error calling GetTopNews: %v", err)
	}

	// Print response
	fmt.Printf("✅ News sent successfully to: %s\n", email)
	for _, article := range resp.Articles {
		fmt.Printf("- %s (%s)\n  %s\n\n", article.Title, article.Author, article.Url)
	}
}