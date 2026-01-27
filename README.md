# ngrokd-go

A Go SDK for connecting to ngrok-bound Kubernetes endpoints from anywhere. Instead of running a daemon, embed this library directly in your application.

## What It Does

This SDK lets your Go application connect to services exposed through ngrok's Kubernetes bindings. It works by:

1. **Discovering endpoints** - Polls the ngrok API to learn which hostnames are ngrok-bound
2. **Establishing mTLS connections** - Connects to ngrok's cloud service with a client certificate via mTLS
3. **Routing intelligently** - ngrok traffic uses the ngrokd dialer, everything else uses your fallback dialer

The SDK provisions its own mTLS certificate by generating a private key locally and having ngrok sign it. The private key never leaves your machine.

## Installation

```bash
go get github.com/ishanj12/ngrokd-go
```

## Usage

```go
import (
    "context"
    "net"
    "net/http"

    ngrokd "github.com/ishanj12/ngrokd-go"
)

ctx := context.Background()

dialer, _ := ngrokd.NewDialer(ctx, ngrokd.Config{
    APIKey:         "your-api-key",
    FallbackDialer: &net.Dialer{},  // for non-ngrok endpoints
})
defer dialer.Close()

// Populate endpoint cache
dialer.DiscoverEndpoints(ctx)

// Create HTTP client with ngrok-aware transport
client := &http.Client{
    Transport: &http.Transport{DialContext: dialer.DialContext},
}

// Use normally - SDK routes ngrok endpoints through the tunnel
resp, _ := client.Get("https://my-service.example.com/api")
resp.Body.Close()
```

## Configuration

```go
ngrokd.Config{
    APIKey:          "your-api-key",
    FallbackDialer:  &net.Dialer{},       // handles non-ngrok traffic
    RefreshInterval: 5 * time.Minute,     // background endpoint refresh
    RefreshOnMiss:   true,                // re-discover on unknown hostname
    RetryConfig: ngrokd.RetryConfig{
        MaxRetries:     3,
        InitialBackoff: 100 * time.Millisecond,
    },
}
```

## Certificate Storage

Certificates are cached locally to avoid re-provisioning on restart:

- `FileStore` (default) - Saves to `~/.ngrokd-go/certs`
- `MemoryStore` - For ephemeral environments like Fargate or Lambda

## License

MIT
