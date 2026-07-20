package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func newFeatureGateTestRouter(t *testing.T, accountFeatures map[string]bool, defaultEnabled bool) *gin.Engine {
	t.Helper()

	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Simulates evo_auth.go having already populated account_features on
	// the request context before RequireFeature runs.
	router.Use(func(c *gin.Context) {
		if accountFeatures != nil {
			ctx := context.WithValue(c.Request.Context(), "account_features", accountFeatures) //nolint:staticcheck
			c.Request = c.Request.WithContext(ctx)
		}
		c.Next()
	})

	router.GET("/agents", RequireFeature("ai_agents", defaultEnabled), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	return router
}

func TestRequireFeature_Blocks(t *testing.T) {
	t.Parallel()

	router := newFeatureGateTestRouter(t, map[string]bool{"ai_agents": false}, true)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/agents", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d. body: %s", w.Code, http.StatusForbidden, w.Body.String())
	}
}

func TestRequireFeature_AllowsWhenEnabled(t *testing.T) {
	t.Parallel()

	router := newFeatureGateTestRouter(t, map[string]bool{"ai_agents": true}, false)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/agents", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d. body: %s", w.Code, http.StatusOK, w.Body.String())
	}
}

func TestRequireFeature_FallsBackToDefaultWhenNoOverride(t *testing.T) {
	t.Parallel()

	router := newFeatureGateTestRouter(t, map[string]bool{}, true)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/agents", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d (default-enabled). body: %s", w.Code, http.StatusOK, w.Body.String())
	}
}

func TestRequireFeature_BlocksWhenNoAccountContextAtAll(t *testing.T) {
	t.Parallel()

	router := newFeatureGateTestRouter(t, nil, false)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/agents", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d. body: %s", w.Code, http.StatusForbidden, w.Body.String())
	}
}
