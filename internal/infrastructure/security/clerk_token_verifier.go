package security

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/MicahParks/keyfunc/v2"
	"github.com/golang-jwt/jwt/v5"
)

type ClerkTokenVerifier struct {
	jwks          *keyfunc.JWKS
	issuer        string
	audience      string
	requiredScope string
}

func NewClerkTokenVerifier(jwksURL, issuer, audience, requiredScope string) (*ClerkTokenVerifier, error) {
	if jwksURL == "" {
		return nil, fmt.Errorf("CLERK_JWKS_URL is required")
	}
	if issuer == "" {
		return nil, fmt.Errorf("CLERK_ISSUER is required")
	}

	jwks, err := keyfunc.Get(jwksURL, keyfunc.Options{
		RefreshInterval:   time.Hour,
		RefreshRateLimit:  5 * time.Minute,
		RefreshTimeout:    10 * time.Second,
		RefreshUnknownKID: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize JWKS: %w", err)
	}

	return &ClerkTokenVerifier{
		jwks:          jwks,
		issuer:        issuer,
		audience:      audience,
		requiredScope: requiredScope,
	}, nil
}

func (v *ClerkTokenVerifier) Verify(ctx context.Context, tokenString string) (string, error) {
	_ = ctx
	parseOptions := []jwt.ParserOption{
		jwt.WithValidMethods([]string{"RS256", "RS384", "RS512"}),
		jwt.WithIssuer(v.issuer),
	}
	if v.audience != "" {
		parseOptions = append(parseOptions, jwt.WithAudience(v.audience))
	}

	token, err := jwt.Parse(tokenString, v.jwks.Keyfunc, parseOptions...)
	if err != nil {
		return "", fmt.Errorf("token validation failed: %w", err)
	}
	if !token.Valid {
		return "", fmt.Errorf("token is invalid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("token claims are not map claims")
	}

	userID, _ := claims["sub"].(string)
	if strings.TrimSpace(userID) == "" {
		return "", fmt.Errorf("token missing subject")
	}

	if v.requiredScope != "" && !claimHasScope(claims, v.requiredScope) {
		return "", fmt.Errorf("token is missing required scope %q", v.requiredScope)
	}

	return userID, nil
}

func claimHasScope(claims jwt.MapClaims, requiredScope string) bool {
	if scopeRaw, ok := claims["scope"]; ok {
		if scopeString, ok := scopeRaw.(string); ok {
			for _, scope := range strings.Fields(scopeString) {
				if scope == requiredScope {
					return true
				}
			}
		}
	}

	if scopeRaw, ok := claims["scp"]; ok {
		if scopeList, ok := scopeRaw.([]interface{}); ok {
			for _, scope := range scopeList {
				if scopeValue, ok := scope.(string); ok && scopeValue == requiredScope {
					return true
				}
			}
		}
	}

	return false
}
