package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
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

	// Create ngrokd dialer (uses NGROK_API_KEY env var)
	dialer, err := ngrokd.NewDialer(ctx, ngrokd.Config{
		DefaultDialer:   &net.Dialer{},
		PollingInterval: 10 * time.Second,
	})
	if err != nil {
		return err
	}
	defer dialer.Close()

	// Discover kubernetes-bound endpoints
	endpoints, err := dialer.DiscoverEndpoints(ctx)
	if err != nil {
		return err
	}

	if len(endpoints) == 0 {
		log.Println("No endpoints found. Run server first:")
		log.Println("  NGROK_AUTHTOKEN=xxx go run examples/server/main.go")
		return nil
	}

	log.Printf("Found %d endpoint(s)", len(endpoints))

	// Create HTTP client with ngrokd transport
	httpClient := &http.Client{
		Transport: &http.Transport{DialContext: dialer.DialContext},
		Timeout:   30 * time.Second,
	}

	for _, ep := range endpoints {
		log.Printf("Connecting to %s...", ep.URL)

		resp, err := httpClient.Get(ep.URL)
		if err != nil {
			log.Printf("  Error: %v", err)
			continue
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		fmt.Printf("  Status: %d\n  Body: %s\n", resp.StatusCode, string(body))
	}

	return nil
}
