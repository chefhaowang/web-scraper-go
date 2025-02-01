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
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewNewsScraperClient(conn)

	// Email of the subscriber
	email := "1242107568@qq.com" 

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*180)
	defer cancel()

	resp, err := client.GetTopNews(ctx, &pb.EmailRequest{Email: email})
	if err != nil {
		log.Fatalf("Error calling GetTopNews: %v", err)
	}

	fmt.Println("✅ News sent successfully to", email)

	for i, article := range resp.Articles {
		fmt.Printf("\n%d. Title: %s\n   Author: %s\n   Date: %s\n   URL: %s\n",
			i+1, article.Title, article.Author, article.Date, article.Url)
	}
}