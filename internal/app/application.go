package app

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"linkedin-mcp/internal/infrastructure/middleware"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func Start() {
	configs := readConfigs()

	server := initServer(configs)

	handler := mcp.NewStreamableHTTPHandler(func(r *http.Request) *mcp.Server {
		return server
	}, &mcp.StreamableHTTPOptions{JSONResponse: true})

	mux := http.NewServeMux()
	mux.Handle(configs.ServerConfig.Path, handler)

	wrappedHandler := middleware.LoggingHandler(mux)

	httpServer := &http.Server{
		Addr:    configs.ServerConfig.BindAddress,
		Handler: wrappedHandler,
	}

	log.Printf("LinkedIn MCP server (streamable HTTP) listening on %s", configs.ServerConfig.PublicURL)

	shutdownCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	serverErrCh := make(chan error, 1)
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			serverErrCh <- err
		}
	}()

	select {
	case <-shutdownCtx.Done():
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(ctx); err != nil {
			log.Printf("graceful shutdown failed: %v", err)
		}
	case err := <-serverErrCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}
}
