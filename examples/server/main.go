package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.ngrok.com/ngrok/v2"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go startHelloServer()

	endpointName := os.Getenv("ENDPOINT_NAME")
	if endpointName == "" {
		endpointName = "hello-server"
	}

	// Internal endpoints use the .internal TLD
	internalURL := fmt.Sprintf("https://%s.internal", endpointName)

	fwd, err := ngrok.Forward(ctx,
		ngrok.WithUpstream("http://localhost:8080"),
		ngrok.WithURL(internalURL),
		ngrok.WithDescription("ngrokd-go example server"),
	)
	if err != nil {
		return fmt.Errorf("failed to create ngrok endpoint: %w", err)
	}

	log.Println("Internal endpoint online:", fwd.URL())
	log.Println("Run client: NGROK_API_KEY=xxx go run examples/client/main.go")

	<-fwd.Done()
	return nil
}

func startHelloServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		fmt.Fprintln(w, "Hello from ngrokd-go!")
	})
	log.Println("Hello server listening on :8080")
	http.ListenAndServe(":8080", nil)
}
