package app

import (
	"context"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func Start() {
	configs := readConfigs()

	server := initServer(configs)

	err := server.Run(context.Background(), &mcp.StdioTransport{})
	if err != nil {
		log.Fatal(err)
	}
}

/*func Start() {
	configs := readConfigs()

	handler := initServer(configs)

	handlerWithLogging := middleware.LoggingHandler(handler)

	log.Printf("MCP server listening on PORT: %s", port)

	// Start the HTTP server with logging handler.
	err := http.ListenAndServe(port, handlerWithLogging)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}*/
