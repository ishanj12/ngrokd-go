package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
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

	// Create kubernetes-bound agent endpoint
	ln, err := ngrok.Listen(ctx,
		ngrok.WithURL("https://hello-server.example"),
		ngrok.WithBindings("kubernetes"),
		ngrok.WithDescription("ngrokd-go example server"),
	)
	if err != nil {
		return fmt.Errorf("failed to create ngrok endpoint: %w", err)
	}

	log.Println("Endpoint online:", ln.URL())
	log.Println("Run client: NGROK_API_KEY=xxx go run examples/client/main.go")

	// Serve hello world
	return http.Serve(ln, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		fmt.Fprintln(w, "Hello from ngrokd-go!")
	}))
}
