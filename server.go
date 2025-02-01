package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	pb "web-scraper-go/scraperpb"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedNewsScraperServer
}

// SendEmail sends the scraped news using AWS SES
func SendEmail(recipient string, articles []*pb.NewsArticle) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-southeast-2"))
	if err != nil {
		return fmt.Errorf("unable to load AWS SDK config, %v", err)
	}

	client := ses.NewFromConfig(cfg)

	subject := "📰 Daily Tech News from The Verge"
	body := "Here are the top 3 tech news articles from The Verge:\n\n"

	for i, article := range articles {
		body += fmt.Sprintf("%d. %s\n   Author: %s\n   Date: %s\n   URL: %s\n\n",
			i+1, article.Title, article.Author, article.Date, article.Url)
	}

	body += "Stay updated with the latest tech trends! 🚀\n\nBest regards,\nWeb Scraper Bot"

	input := &ses.SendEmailInput{
		Destination: &types.Destination{ // ✅ Use types.Destination
			ToAddresses: []string{recipient},
		},
		Message: &types.Message{ // ✅ Use types.Message
			Body: &types.Body{ // ✅ Use types.Body
				Text: &types.Content{ // ✅ Use types.Content
					Data: aws.String(body),
				},
				Html: &types.Content{
					Data: aws.String("<h1>Daily Tech News</h1><p>" + body + "</p>"),
				},
			},
			Subject: &types.Content{
				Data: aws.String(subject),
			},
		},
		Source: aws.String("wh.tenghe@gmail.com"), // ✅ Make sure this email is verified in AWS SES
	}

	_, err = client.SendEmail(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	log.Println("📧 Email sent successfully to", recipient)
	return nil
}

// gRPC method to get top news and send via email
func (s *server) GetTopNews(ctx context.Context, req *pb.EmailRequest) (*pb.NewsResponse, error) {
	log.Println("📥 Received gRPC request to scrape news for:", req.Email)

	articles := ScrapeTopNews()

	if err := SendEmail(req.Email, articles); err != nil {
		log.Printf("❌ Failed to send email to %s: %v", req.Email, err)
		return nil, err
	}

	log.Println("✅ Successfully scraped and emailed news articles")
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