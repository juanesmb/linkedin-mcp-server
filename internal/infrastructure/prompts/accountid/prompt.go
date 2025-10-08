package accountid

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Prompt struct{}

func NewPrompt() *Prompt {
	return &Prompt{}
}

func (p *Prompt) GetPrompt(ctx context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	var promptText string
	if req.Params.Arguments != nil {
		if accountID, exists := req.Params.Arguments["accountID"]; exists && accountID != "" {
			promptText = fmt.Sprintf("Great! You have provided the LinkedIn Ad Account ID: %s. You can now use the LinkedIn tools with this Account ID as the 'accountID' parameter.", accountID)
		} else {
			promptText = `You are working with LinkedIn advertising tools. Before using any LinkedIn tools, you MUST first ask the user for their LinkedIn Ad Account ID.

WORKFLOW:
1. Ask the user: "What is your LinkedIn Ad Account ID? (This should be a numeric value, e.g., 512345678)"
2. Wait for the user to provide the Account ID
3. Only then proceed to use the LinkedIn tools with the provided Account ID

The Account ID is required for both tools and should be provided as the 'accountID' parameter.`
		}
	} else {
		promptText = `You are working with LinkedIn advertising tools. Before using any LinkedIn tools, you MUST first ask the user for their LinkedIn Ad Account ID.

WORKFLOW:
1. Ask the user: "What is your LinkedIn Ad Account ID? (This should be a numeric value, e.g., 512345678)"
2. Wait for the user to provide the Account ID
3. Only then proceed to use the LinkedIn tools with the provided Account ID

The Account ID is required for both tools and should be provided as the 'accountID' parameter.`
	}

	return &mcp.GetPromptResult{
		Messages: []*mcp.PromptMessage{
			{
				Role: "user",
				Content: &mcp.TextContent{
					Text: promptText,
				},
			},
		},
	}, nil
}
