package pyroscopeextension

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/extension/extensiontest"
)

func TestFactory(t *testing.T) {
	t.Parallel()

	f := NewFactory()
	assert.Equal(t, component.Type("pyroscope"), f.Type(), "Must match the expected value")
	assert.Equal(t, component.StabilityLevelBeta, f.ExtensionStability(), "Must match the expected stability")
}

func TestFactoryNewExtension(t *testing.T) {
	t.Parallel()

	f := NewFactory()

	ext, err := f.CreateExtension(
		context.Background(),
		extensiontest.NewNopCreateSettings(),
		newDefaultConfig(),
	)
	assert.NoError(t, err, "Must not error when creating extension")
	assert.NotNil(t, ext, "Must have a valid extension")
}
