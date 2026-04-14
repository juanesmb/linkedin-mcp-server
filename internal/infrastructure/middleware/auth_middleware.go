package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type authContextKey struct{}

type UnauthorizedResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type TokenVerifier interface {
	Verify(ctx context.Context, token string) (string, error)
}

func RequireBearerAuth(verifier TokenVerifier, resourceMetadataURL, requiredScope string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearerToken := extractBearerToken(r.Header.Get("Authorization"))
		if bearerToken == "" {
			writeUnauthorized(w, resourceMetadataURL, requiredScope, "missing bearer token")
			return
		}

		userID, err := verifier.Verify(r.Context(), bearerToken)
		if err != nil {
			writeUnauthorized(w, resourceMetadataURL, requiredScope, err.Error())
			return
		}

		ctx := context.WithValue(r.Context(), authContextKey{}, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(authContextKey{}).(string)
	return userID, ok && userID != ""
}

func writeUnauthorized(w http.ResponseWriter, resourceMetadataURL, requiredScope, detail string) {
	challenge := fmt.Sprintf(`Bearer resource_metadata="%s"`, resourceMetadataURL)
	if requiredScope != "" {
		challenge += fmt.Sprintf(`, scope="%s"`, requiredScope)
	}

	w.Header().Set("WWW-Authenticate", challenge)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)

	_ = json.NewEncoder(w).Encode(UnauthorizedResponse{
		Error:   "unauthorized",
		Message: detail,
	})
}

func extractBearerToken(header string) string {
	parts := strings.Fields(strings.TrimSpace(header))
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}
	return strings.TrimSpace(parts[1])
}
