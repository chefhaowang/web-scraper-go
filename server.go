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

// SendEmail sends the scraped news using AWS SES with a structured format
func SendEmail(recipient string, articles []*pb.NewsArticle) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-southeast-2"))
	if err != nil {
		return fmt.Errorf("unable to load AWS SDK config, %v", err)
	}

	client := ses.NewFromConfig(cfg)

	subject := "📰 Daily Tech News from The Verge"

	plainBody := "Here are the top 3 tech news articles from The Verge:\n\n"
	for i, article := range articles {
		plainBody += fmt.Sprintf("%d. %s\n   Author: %s\n   Date: %s\n   URL: %s\n\n",
			i+1, article.Title, article.Author, article.Date, article.Url)
	}
	plainBody += "Stay updated with the latest tech trends! 🚀\n\nBest regards,\nWeb Scraper Bot"

	htmlBody := `<html>
	<head>
		<style>
			body { font-family: Arial, sans-serif; background-color: #f4f4f4; padding: 20px; }
			.container { background-color: #fff; padding: 20px; border-radius: 8px; box-shadow: 0 2px 5px rgba(0,0,0,0.1); }
			h1 { color: #333; }
			.article { margin-bottom: 20px; padding: 15px; border-bottom: 1px solid #ddd; }
			.article:last-child { border-bottom: none; }
			.title { font-size: 18px; font-weight: bold; color: #0073e6; text-decoration: none; }
			.meta { font-size: 14px; color: #555; margin: 5px 0; }
			.footer { font-size: 12px; color: #888; margin-top: 20px; text-align: center; }
		</style>
	</head>
	<body>
		<div class="container">
			<h1>📰 Top 3 Tech News from The Verge</h1>`

	// Add each article in a structured format
	for i, article := range articles {
		htmlBody += fmt.Sprintf(`
			<div class="article">
				<a href="%s" class="title">%d. %s</a>
				<div class="meta">Author: %s | Date: %s</div>
			</div>`, article.Url, i+1, article.Title, article.Author, article.Date)
	}

	// Footer
	htmlBody += `
			<div class="footer">
				Stay updated with the latest tech trends! 🚀<br>
				<em>Sent by Web Scraper Bot</em>
			</div>
		</div>
	</body>
</html>`

	// SES Email Input (ses sdk calling)
	input := &ses.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: []string{recipient},
		},
		Message: &types.Message{
			Body: &types.Body{
				Text: &types.Content{
					Data: aws.String(plainBody),
				},
				Html: &types.Content{
					Data: aws.String(htmlBody),
				},
			},
			Subject: &types.Content{
				Data: aws.String(subject),
			},
		},
		Source: aws.String("wh.tenghe@gmail.com"), // this email is verified in AWS SES
	}

	// Send the email
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