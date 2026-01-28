package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	ngrokd "github.com/ishanj12/ngrokd-go"
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	// Target URL - the kubernetes-bound endpoint forwarding to hello-server
	targetURL := os.Getenv("TARGET_URL")
	if targetURL == "" {
		targetURL = "https://hello.example"
	}

	// Create ngrokd dialer (uses NGROK_API_KEY env var)
	dialer, err := ngrokd.NewDialer(ctx, ngrokd.Config{
		DefaultDialer:   &net.Dialer{},
		PollingInterval: 10 * time.Second,
	})
	if err != nil {
		return err
	}
	defer dialer.Close()

	log.Println("Operator ID:", dialer.OperatorID())

	// Discover kubernetes-bound endpoints
	endpoints, err := dialer.DiscoverEndpoints(ctx)
	if err != nil {
		return err
	}

	log.Printf("Found %d endpoint(s)", len(endpoints))
	for _, ep := range endpoints {
		log.Printf("  - %s", ep.URL)
	}

	// Create HTTP client with ngrokd transport
	httpClient := &http.Client{
		Transport: &http.Transport{DialContext: dialer.DialContext},
		Timeout:   30 * time.Second,
	}

	// Connect to the target endpoint
	log.Printf("Connecting to %s...", targetURL)

	resp, err := httpClient.Get(targetURL)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Status: %d\nBody: %s\n", resp.StatusCode, string(body))

	return nil
}
