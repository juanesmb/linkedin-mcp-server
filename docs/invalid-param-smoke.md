# Invalid Parameter Smoke Scenario

This manual smoke scenario validates that a LinkedIn parameter validation failure keeps context end-to-end and is visible to MCP clients.

## Preconditions
- Jumon web app is running with internal LinkedIn proxy route enabled.
- MCP server is running and configured with valid auth/gateway env vars.
- Test user is authenticated and has a LinkedIn connection.

## Steps
1. Call `search_ad_accounts` with an intentionally malformed search payload that causes LinkedIn to return `PARAM_INVALID`.
2. Observe Jumon proxy response from `/api/internal/linkedin/proxy/[...path]`:
   - HTTP `400`
   - JSON with:
     - `code: "LINKEDIN_PARAM_INVALID"`
     - `message`
     - `providerStatus`
     - `inputErrors[]` including `fieldPath`.
3. Observe MCP server tool response:
   - MCP tool error (not protocol crash)
   - message includes:
     - LinkedIn rejected parameters
     - problematic field path(s)
     - guidance to adjust params and retry.

## Expected Result
The MCP client can retry with corrected parameters using the field-level context, instead of receiving a generic `status 500` failure.
