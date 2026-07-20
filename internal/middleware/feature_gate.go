package middleware

import (
	"net/http"

	"evo-ai-core-service/internal/httpclient/response"
	"evo-ai-core-service/internal/utils/contextutils"

	"github.com/gin-gonic/gin"
)

// RequireFeature returns a middleware that rejects the request with 403
// unless the requesting Account has featureName enabled (see
// contextutils.IsFeatureEnabled and specs/account-feature-toggles). Mirrors
// evo-ai-crm-community's FeatureGateConcern#require_feature - this is the
// Go-side equivalent, checking the same per-Account overrides independently
// (see specs/account-feature-toggles/04-architecture.md, revised Decision
// C1), since most AI Agents traffic reaches this service directly rather
// than through the CRM's Rails proxy.
func RequireFeature(featureName string, defaultEnabled bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !contextutils.IsFeatureEnabled(c.Request.Context(), featureName, defaultEnabled) {
			response.ErrorResponse(c, "ERR_FEATURE_NOT_AVAILABLE",
				"The "+featureName+" feature is not enabled for this account", nil, http.StatusForbidden)
			c.Abort()
			return
		}

		c.Next()
	}
}
