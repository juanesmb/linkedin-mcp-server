package app

import (
	"log"
	"net/http"

	"linkedin-mcp/internal/infrastructure/middleware"
)

const port = ":8080"

func Start() {
	configs := readConfigs()

	handler := initServer(configs)

	handlerWithLogging := middleware.LoggingHandler(handler)

	log.Printf("MCP server listening on PORT: %s", port)

	// Start the HTTP server with logging handler.
	err := http.ListenAndServe(port, handlerWithLogging)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

/* Config to start an MCP server using STDIO communication protocol
	func Start() {
	// Create a server
	server := mcp.NewServer(&mcp.Implementation{Name: "LinkedIn", Version: "v1.0.0"}, nil)

	// Add tools
	// mcp.AddTool(server, &mcp.Tool{Name: "greet", Description: "say hi"}, greeter.SayHi)
	mcp.AddTool(server, &mcp.Tool{Name: "search_campaigns", Description: "Search for LinkedIn ad campaigns"}, searchcampaigns.SearchCampaigns)

	// Run the server over stdin/stdout, until the client disconnects.
	err := server.Run(context.Background(), &mcp.StdioTransport{})
	if err != nil {
		log.Fatal(err)
	}
}
*/
