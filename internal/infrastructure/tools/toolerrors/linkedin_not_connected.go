package toolerrors

import (
	"fmt"

	"linkedin-mcp/internal/infrastructure/api/gateway"
)

func WrapToolExecutionError(operation string, err error, connectURL string) error {
	if gateway.IsLinkedInNotConnected(err) {
		return fmt.Errorf(
			"cannot %s because no LinkedIn account is connected for this user. Ask the user to open %s, complete LinkedIn connection, and then retry this tool call",
			operation,
			connectURL,
		)
	}

	return fmt.Errorf("failed to %s: %w", operation, err)
}
