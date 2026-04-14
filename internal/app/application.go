package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"linkedin-mcp/internal/infrastructure/middleware"
	"linkedin-mcp/internal/infrastructure/security"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func Start() {
	configs := readConfigs()
	components := initCommonComponents(configs)

	server := initServer(configs, components)
	verifier, err := security.NewClerkTokenVerifier(
		configs.AuthConfig.ClerkJWKSURL,
		configs.AuthConfig.ClerkIssuer,
		configs.AuthConfig.ClerkAudience,
		configs.AuthConfig.RequiredScope,
	)
	if err != nil {
		log.Fatal(err)
	}

	handler := mcp.NewStreamableHTTPHandler(func(r *http.Request) *mcp.Server {
		return server
	}, &mcp.StreamableHTTPOptions{JSONResponse: true})

	mux := http.NewServeMux()
	metadataPath := "/.well-known/oauth-protected-resource"
	resourceMetadataURL := appendURLPath(configs.ServerConfig.PublicURL, metadataPath)
	if resourceMetadataURL == "" {
		resourceMetadataURL = metadataPath
	}

	protectedMCPHandler := middleware.RequireBearerAuth(
		verifier,
		resourceMetadataURL,
		configs.AuthConfig.RequiredScope,
		handler,
	)
	mux.Handle(configs.ServerConfig.Path, protectedMCPHandler)
	mux.HandleFunc(metadataPath, oauthProtectedResourceHandler(configs))
	mux.HandleFunc(metadataPath+configs.ServerConfig.Path, oauthProtectedResourceHandler(configs))

	wrappedHandler := middleware.LoggingHandler(mux)

	httpServer := &http.Server{
		Addr:    configs.ServerConfig.BindAddress,
		Handler: wrappedHandler,
	}

	log.Printf("LinkedIn MCP server (streamable HTTP) listening on path %s (bind %s)", configs.ServerConfig.Path, configs.ServerConfig.BindAddress)

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

func oauthProtectedResourceHandler(configs Configs) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		publicURL := strings.TrimRight(configs.ServerConfig.PublicURL, "/")
		if publicURL == "" {
			scheme := "https"
			if proto := r.Header.Get("X-Forwarded-Proto"); proto != "" {
				scheme = proto
			} else if r.TLS == nil {
				scheme = "http"
			}
			publicURL = fmt.Sprintf("%s://%s", scheme, r.Host)
		}

		resource := appendURLPath(publicURL, configs.ServerConfig.Path)
		response := map[string]any{
			"resource":              resource,
			"authorization_servers": []string{configs.AuthConfig.AuthorizationURL},
			"bearer_methods_supported": []string{
				"header",
			},
		}
		if configs.AuthConfig.RequiredScope != "" {
			response["scopes_supported"] = []string{configs.AuthConfig.RequiredScope}
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "failed to write metadata response", http.StatusInternalServerError)
		}
	}
}

func appendURLPath(baseURL, path string) string {
	base := strings.TrimRight(baseURL, "/")
	normalizedPath := "/" + strings.TrimLeft(path, "/")
	if base == "" {
		return ""
	}
	return base + normalizedPath
}
