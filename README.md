# LinkedIn MCP Server

A Model Context Protocol (MCP) server that exposes LinkedIn Advertising capabilities—searching ad accounts, exploring campaigns, and retrieving analytics insights. Use it to connect MCP-compatible clients (e.g., Claude Desktop) to the LinkedIn Ads API via a single, structured interface.

## Highlights
- Streamable HTTP transport for remote connector support.
- Tools, prompts, and resources tailored to LinkedIn Ads workflows.
- Written in Go using the official [modelcontextprotocol/go-sdk](https://github.com/modelcontextprotocol/go-sdk/tree/main).

## Prerequisites
- Go 1.25+
- LinkedIn Ads REST access token with the necessary scopes
- Optional: Docker for containerized builds

## Quick Start (local)
```bash
git clone https://github.com/your-org/linkedin-mcp.git
cd linkedin-mcp/server

export LINKEDIN_ACCESS_TOKEN="<your_linkedin_token>"

go run ./...
# or
go build -o linkedin-mcp
./linkedin-mcp
```
The server listens on `http://127.0.0.1:8080/mcp`. You can point an MCP client at that URL to begin calling tools.

## Docker
```bash
docker build -t linkedin-mcp:local .
docker run --rm -p 8080:8080 \
  -e LINKEDIN_ACCESS_TOKEN="$LINKEDIN_ACCESS_TOKEN" \
  linkedin-mcp:local
```

## Configuration
Environment variables:
- `LINKEDIN_ACCESS_TOKEN` (required): LinkedIn Ads API bearer token
- `PORT` (optional): port to bind (default `8080`)
- `MCP_SERVER_HOST` (optional): host interface (default `0.0.0.0`)
- `MCP_SERVER_PATH` (optional): MCP endpoint path (default `/mcp`)

## Remote MCP Connector (Claude Desktop example)
1. Run or deploy the server (see Cloud Run guide below).
2. In Claude Desktop → Settings → Connectors → Add Remote MCP.
3. Use the base URL `https://your-domain/mcp` (or `http://127.0.0.1:8080/mcp` for local testing).
4. Call the `system_guidelines` prompt first to understand available tools and required arguments.

## Testing
```bash
go test ./...
```

## Contributing
Issues and pull requests are welcome. Please:
- Include context about LinkedIn API usage or MCP behavior.
- Add or update tests when touching business logic.
- Follow Go formatting (`gofmt`) before submitting.
