package contextutils

import (
	"context"
	"errors"

	"evo-ai-core-service/internal/types"

	"github.com/google/uuid"
)

var ErrUnauthorized = errors.New("unauthorized")

func GetUserID(ctx context.Context) (uuid.UUID, error) {
	if userID, ok := ctx.Value("user_id").(uuid.UUID); ok {
		return userID, nil
	}

	return uuid.Nil, errors.New("user_id not found in context")
}

func GetUserEmail(ctx context.Context) (string, error) {
	if email, ok := ctx.Value("email").(string); ok {
		return email, nil
	}

	return "", errors.New("email not found in context")
}

func GetUserName(ctx context.Context) (string, error) {
	if name, ok := ctx.Value("name").(string); ok {
		return name, nil
	}

	return "", errors.New("name not found in context")
}

// GetAccountID returns the tenant Account ID resolved from the validated
// auth token (see internal/middleware/evo_auth.go). Returns an error if the
// token had no associated account - callers that need tenant scoping (e.g.
// pkg/agent/repository) must treat that as a hard failure, not fall back to
// an unscoped query.
func GetAccountID(ctx context.Context) (uuid.UUID, error) {
	if accountID, ok := ctx.Value("account_id").(uuid.UUID); ok {
		return accountID, nil
	}

	return uuid.Nil, errors.New("account_id not found in context")
}

// IsFeatureEnabled reports whether featureName is enabled for the
// requesting Account, per the per-Account feature_overrides resolved from
// the auth-service validation response (see internal/middleware/evo_auth.go
// and specs/account-feature-toggles). When the Account has no explicit
// override for featureName - the common case - defaultEnabled is returned
// instead; it must match the corresponding entry's declared default in
// evo-ai-crm-community/config/features.yml, since this service has no local
// copy of that file to consult.
func IsFeatureEnabled(ctx context.Context, featureName string, defaultEnabled bool) bool {
	features, ok := ctx.Value("account_features").(map[string]bool)
	if !ok || features == nil {
		return defaultEnabled
	}

	if enabled, present := features[featureName]; present {
		return enabled
	}

	return defaultEnabled
}

func GetToken(ctx context.Context) (string, error) {
	if token, ok := ctx.Value("token").(string); ok {
		return token, nil
	}

	return "", errors.New("token not found in context")
}

func GetApiAccessToken(ctx context.Context) (string, error) {
	if token, ok := ctx.Value("api_access_token").(string); ok {
		return token, nil
	}

	return "", errors.New("token not found in context")
}

func GetTokenType(ctx context.Context) (string, error) {
	if tokenType, ok := ctx.Value("token_type").(string); ok {
		return tokenType, nil
	}

	return "", errors.New("token type not found in context")
}

func GetAuthHeaders(ctx context.Context) (interface{}, error) {
	if headers := ctx.Value("auth_headers"); headers != nil {
		return headers, nil
	}

	return nil, errors.New("auth_headers not found in context")
}

// GetUser returns the complete user information from context
func GetUser(ctx context.Context) (types.EvoAuthUser, error) {
	if user, ok := ctx.Value("user").(types.EvoAuthUser); ok {
		return user, nil
	}

	return types.EvoAuthUser{}, errors.New("user not found in context")
}
