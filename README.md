# LinkedIn MCP Server

Remote-ready MCP server that exposes LinkedIn Advertising tools, resources, and prompts via the streamable HTTP transport. For background on remote connectors, see [Claude’s custom connector docs](https://support.claude.com/en/articles/11503834-building-custom-connectors-via-remote-mcp-servers).

## Configure & Run

1. Export your LinkedIn API token and optionally adjust the server bind address/path:

   ```bash
   export LINKEDIN_ACCESS_TOKEN="<token>"
   ```

2. Build or run directly with Go:

   ```bash
   go run ./...
   # or
   go build -o linkedin-mcp
   ./linkedin-mcp
   ```

The server listens on `http://127.0.0.1:8080` and serves the MCP streamable endpoint at `/mcp` with JSON responses. Clients such as Claude Desktop can connect using the "Remote MCP" connector flow.

## Remote Connector Registration

When adding the server in Claude’s connector UI, specify the base URL `http://127.0.0.1:8080/mcp`. OAuth is not required; the server currently expects the LinkedIn token via environment variables.

## Troubleshooting

- Ensure `LINKEDIN_ACCESS_TOKEN` is set and valid before starting the server.
- Verify the MCP endpoint (`http://127.0.0.1:8080/mcp`) is reachable from clients.
