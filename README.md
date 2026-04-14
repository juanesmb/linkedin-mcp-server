# LinkedIn MCP Server

A Model Context Protocol (MCP) server that exposes LinkedIn Advertising capabilities—searching ad accounts, exploring campaigns, and retrieving analytics insights. Use it to connect MCP-compatible clients (e.g., Claude Desktop) to the LinkedIn Ads API via a single, structured interface.

## Highlights
- Streamable HTTP transport for remote connector support.
- Tools and resources tailored to LinkedIn Ads workflows.
- Server-level MCP instructions guide tool usage and sequencing.
- MCP OAuth protected resource metadata and `WWW-Authenticate` bearer challenges.
- Clerk bearer-token validation with user-scoped delegation to the Jumon gateway.
- Written in Go using the official [modelcontextprotocol/go-sdk](https://github.com/modelcontextprotocol/go-sdk/tree/main).

## Prerequisites
- Go 1.25+
- A Clerk app issuing access tokens for your MCP client
- Jumon web app deployment with internal gateway endpoints enabled
- Optional: Docker for containerized builds

## Quick Start (local)
```bash
git clone https://github.com/your-org/linkedin-mcp.git
cd linkedin-mcp/server

cp .env.example .env
# Fill in Clerk + gateway values, then:
set -a; source .env; set +a

go run ./...
# or
go build -o linkedin-mcp
./linkedin-mcp
```
The server listens on `http://127.0.0.1:8080/mcp` and exposes OAuth protected resource metadata at `/.well-known/oauth-protected-resource`.

## Docker
```bash
docker build -t linkedin-mcp:local .
docker run --rm -p 8080:8080 \
  --env-file .env \
  linkedin-mcp:local
```

## Configuration
Environment variables:
- `CLERK_ISSUER` (required): Clerk token issuer URL
- `CLERK_JWKS_URL` (required): Clerk JWKS endpoint used for JWT signature validation
- `CLERK_AUDIENCE` (optional): expected audience in incoming access tokens
- `MCP_REQUIRED_SCOPE` (optional): required token scope (enforced if provided)
- `AUTHORIZATION_SERVER_URL` (optional): authorization server URL advertised in metadata (defaults to `CLERK_ISSUER`)
- `JUMON_GATEWAY_BASE_URL` (required): Jumon web base URL (for `/api/internal/*` calls)
- `JUMON_GATEWAY_INTERNAL_SECRET` (required): internal secret sent as `x-gateway-secret`
- `PORT` (optional): port to bind (default `8080`)
- `MCP_SERVER_HOST` (optional): host interface (default `0.0.0.0`)
- `MCP_SERVER_PATH` (optional): MCP endpoint path (default `/mcp`)
- `PUBLIC_BASE_URL` (optional): absolute public URL used in metadata/challenges (recommended in production)

By default the HTTP server binds to `0.0.0.0:8080` and serves MCP on `/mcp`.
Server instructions are loaded from `internal/app/instructions/server_instructions.md` at startup; if the file is missing or empty, startup fails.

## Runtime auth flow
1. MCP client sends `Authorization: Bearer <access_token>` to `/mcp`.
2. Server validates token signature/claims against Clerk JWKS and extracts the Clerk user ID (`sub`).
3. LinkedIn tool execution is delegated to Jumon internal endpoints:
   - `GET /api/internal/connections/linkedin/current`
   - `POST /api/internal/linkedin/proxy/*`
   - `POST /api/internal/linkedin/refresh`
4. Jumon decrypts user provider tokens and performs LinkedIn API requests.

The MCP server never stores or uses static LinkedIn access tokens.

## Remote MCP Connector (Claude Desktop example)
1. Run or deploy the server (see Cloud Run guide below).
2. In Claude Desktop → Settings → Connectors → Add Remote MCP.
3. Use the base URL `https://your-domain/mcp` (or `http://127.0.0.1:8080/mcp` for local testing).
4. Complete OAuth with Clerk when prompted by the client.
5. Follow the server `instructions` and read analytics resources (`linkedin://analytics/parameters` and `linkedin://analytics/metrics`) before calling `get_analytics`.

## Testing
```bash
go test ./...
```

## Contributing
Issues and pull requests are welcome. Please:
- Include context about LinkedIn API usage or MCP behavior.
- Add or update tests when touching business logic.
- Follow Go formatting (`gofmt`) before submitting.
