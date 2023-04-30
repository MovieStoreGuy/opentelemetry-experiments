package pyroscopeextension

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/extension/extensiontest"
	"go.uber.org/multierr"
)

func TestExtensionLifecycle(t *testing.T) {
	t.Parallel()

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := io.Copy(io.Discard, r.Body)
		assert.NoError(t, multierr.Combine(err, r.Body.Close()), "Must not error when reading content")

		// Blinding do nothing
	}))
	t.Cleanup(s.Close)

	py, err := newPyroscopeProfiler(context.Background(), extensiontest.NewNopCreateSettings(), &Config{
		ApplicationName: "test",
		Endpoint:        s.URL,
	})
	require.NoError(t, err, "Must not error when creating pyroscope extension")

	assert.NoError(t, py.Start(context.Background(), componenttest.NewNopHost()), "Must not error when starting extension")
	assert.ErrorIs(t, py.Start(context.Background(), componenttest.NewNopHost()), ErrAlreadyStarted, "Must error with already started")
	assert.NoError(t, py.Shutdown(context.Background()), "Must not error when trying to shudown")
	assert.ErrorIs(t, py.Shutdown(context.Background()), ErrNeverStarted, "Must report an error if trying to shutdown twice")
}
