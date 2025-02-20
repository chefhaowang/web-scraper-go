package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	pb "web-scraper-go/scraperpb"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

var client pb.NewsScraperClient

func getNewsHandler(w http.ResponseWriter, r *http.Request) {
	// Get email from query parameters
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "Missing 'email' query parameter", http.StatusBadRequest)
		return
	}

	// Set gRPC context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*180)
	defer cancel()

	// Call gRPC GetTopNews
	resp, err := client.GetTopNews(ctx, &pb.EmailRequest{Email: email})
	if err != nil {
		log.Printf("Error calling GetTopNews: %v", err)
		http.Error(w, "Failed to fetch news", http.StatusInternalServerError)
		return
	}

	// Create response structure
	response := struct {
		Message string          `json:"message"`
		Email   string          `json:"email"`
		News    []*pb.NewsArticle `json:"news"`
	}{
		Message: "✅ News sent successfully",
		Email:   email,
		News:    resp.Articles,
	}

	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	// Connect to gRPC server
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()
	client = pb.NewNewsScraperClient(conn)

	// Set up HTTP server
	r := mux.NewRouter()
	r.HandleFunc("/get-news", getNewsHandler).Methods("GET")

	fmt.Println("🚀 REST API running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}