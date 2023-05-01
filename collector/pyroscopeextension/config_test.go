package pyroscopeextension

import (
	"path"
	"testing"

	"github.com/pyroscope-io/client/pyroscope"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/confmap/confmaptest"
)

func TestConfigMapstructure(t *testing.T) {
	t.Parallel()

	assert.NoError(t, componenttest.CheckConfigStruct(&Config{}), "Must not error when parsing config")
}

func TestConfigLoaded(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		id     component.ID
		expect *Config
		err    error
	}{
		{
			id: component.NewIDWithName(typeStr, "min-values"),
			expect: &Config{
				RuntimeMutexProfileFraction: 5,
				RuntimeBlockProfileFraction: 5,
				Endpoint:                    "http://pyroscope-server:4040",
				ApplicationName:             "open-telemetry-collector",
				Profiles:                    pyroscope.DefaultProfileTypes,
			},
			err: nil,
		},
		{
			id: component.NewIDWithName(typeStr, "all-values"),
			expect: &Config{
				RuntimeMutexProfileFraction: 5,
				RuntimeBlockProfileFraction: 5,
				Endpoint:                    "http://pyroscope-server:4040",
				AuthToken:                   "mycloudtoken",
				ApplicationName:             "open-telemetry collector",
				Tags: map[string]string{
					"service.name":           "opentelemetry-demo",
					"service.version":        "v1.0.4",
					"deployment.environment": "test",
				},
				Profiles: []pyroscope.ProfileType{
					pyroscope.ProfileCPU,
					pyroscope.ProfileAllocObjects,
					pyroscope.ProfileAllocSpace,
					pyroscope.ProfileInuseObjects,
					pyroscope.ProfileInuseSpace,
					pyroscope.ProfileGoroutines,
					pyroscope.ProfileMutexCount,
					pyroscope.ProfileMutexDuration,
					pyroscope.ProfileBlockCount,
					pyroscope.ProfileBlockDuration,
				},
			},
		},
		{
			id:     component.NewIDWithName(typeStr, "missing-endpoint"),
			expect: newDefaultConfig().(*Config),
			err:    ErrRequiresEndpoint,
		},
		{
			id: component.NewIDWithName(typeStr, "negative-runtime-fractions"),
			expect: &Config{
				RuntimeMutexProfileFraction: -1,
				RuntimeBlockProfileFraction: -1,
				Endpoint:                    "http://pyroscope-server:4040",
				ApplicationName:             "open-telemetry-collector",
				Profiles:                    pyroscope.DefaultProfileTypes,
			},
			err: ErrRequiresPositiveValue,
		},
		{
			id: component.NewIDWithName(typeStr, "invalid-profile"),
			expect: &Config{
				RuntimeMutexProfileFraction: 5,
				RuntimeBlockProfileFraction: 5,
				Endpoint:                    "http://pyroscope-server:4040",
				ApplicationName:             "open-telemetry-collector",
				Profiles:                    []pyroscope.ProfileType{"*"},
			},
			err: ErrRequiresValidProfile,
		},
	} {
		tc := tc
		t.Run(tc.id.String(), func(t *testing.T) {
			t.Parallel()

			cm, err := confmaptest.LoadConf(path.Join("testdata", "config.yml"))
			require.NoError(t, err, "Must not error when loading configuration")

			factory := NewFactory()
			cfg := factory.CreateDefaultConfig()

			sub, err := cm.Sub(tc.id.String())
			require.NoError(t, err, "Must not error trying to fetch details")
			require.NoError(t, component.UnmarshalConfig(sub, cfg))

			assert.ErrorIs(t, component.ValidateConfig(cfg), tc.err, "Must match the expected error")
			assert.EqualValues(t, tc.expect, cfg, "Must match the expected config")
		})
	}
}
