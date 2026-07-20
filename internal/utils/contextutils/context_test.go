package contextutils

import (
	"context"
	"testing"
)

func TestIsFeatureEnabled(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		features       map[string]bool
		featureName    string
		defaultEnabled bool
		want           bool
	}{
		{
			name:           "no account_features in context falls back to default (true)",
			features:       nil,
			featureName:    "ai_agents",
			defaultEnabled: true,
			want:           true,
		},
		{
			name:           "no account_features in context falls back to default (false)",
			features:       nil,
			featureName:    "ai_agents",
			defaultEnabled: false,
			want:           false,
		},
		{
			name:           "explicit false override wins over a true default",
			features:       map[string]bool{"ai_agents": false},
			featureName:    "ai_agents",
			defaultEnabled: true,
			want:           false,
		},
		{
			name:           "explicit true override wins over a false default",
			features:       map[string]bool{"ai_agents": true},
			featureName:    "ai_agents",
			defaultEnabled: false,
			want:           true,
		},
		{
			name:           "override for an unrelated feature does not affect this one",
			features:       map[string]bool{"pipelines": false},
			featureName:    "ai_agents",
			defaultEnabled: true,
			want:           true,
		},
		{
			name:           "empty (but non-nil) overrides map falls back to default",
			features:       map[string]bool{},
			featureName:    "ai_agents",
			defaultEnabled: true,
			want:           true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			if tt.features != nil {
				ctx = context.WithValue(ctx, "account_features", tt.features) //nolint:staticcheck // matches the string-key convention already used throughout this package
			}

			got := IsFeatureEnabled(ctx, tt.featureName, tt.defaultEnabled)
			if got != tt.want {
				t.Errorf("IsFeatureEnabled(%v, %q, %v) = %v, want %v", tt.features, tt.featureName, tt.defaultEnabled, got, tt.want)
			}
		})
	}

	t.Run("value of the wrong type in context falls back to default instead of panicking", func(t *testing.T) {
		t.Parallel()

		ctx := context.WithValue(context.Background(), "account_features", "not-a-map") //nolint:staticcheck
		got := IsFeatureEnabled(ctx, "ai_agents", true)
		if got != true {
			t.Errorf("IsFeatureEnabled with wrong-typed context value = %v, want true (default)", got)
		}
	})
}
