package accountid

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	accountIDRequestText = `Before using LinkedIn tools that require an account ID (e.g., campaign search or analytics), collect the user's LinkedIn Ad Account ID.

Ask the user concisely: "What is your LinkedIn Ad Account ID? (numeric value, e.g., 512345678)"

Once provided, use it as the 'accountID' parameter for those tools. You can still use the ad account discovery tool without one.`
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
			promptText = accountIDRequestText
		}
	} else {
		promptText = accountIDRequestText
	}

	return &mcp.GetPromptResult{
		Messages: []*mcp.PromptMessage{
			{
				Role: "system",
				Content: &mcp.TextContent{
					Text: promptText,
				},
			},
		},
	}, nil
}
