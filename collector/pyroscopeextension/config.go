package pyroscopeextension

import (
	"errors"
	"fmt"

	"github.com/pyroscope-io/client/pyroscope"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configopaque"
	"go.opentelemetry.io/collector/confmap"
	"go.uber.org/multierr"
)

var (
	ErrRequiresPositiveValue = errors.New("required positive value")
	ErrRequiresEndpoint      = errors.New("required endpoint")
	ErrRequiresValidProfile  = errors.New("required valid profile")
)

type Config struct {
	RuntimeMutexProfileFraction int `mapstructure:"runtime_mutex_fraction"`
	RuntimeBlockProfileFaction  int `mapstructure:"runtime_block_fraction"`
	// Endpoint is server address to forward data to from the collector
	Endpoint string `mapstructure:"endpoint"`
	// AuthToken is used to provide authentication for the pyroscope server
	// Required if you're using pyroscope cloud.
	AuthToken configopaque.String `mapstructure:"auth_token"`
	// ApplicationName is used to set what application name should be assocated
	// with profiles being sent to pyroscope.
	ApplicationName string `mapstructure:"application_name"`
	// Tags are a set of attirbutes that help given context on how the
	// collector is configured and its environment.
	Tags map[string]string `mapstructure:"tags"`
	// Profiles are used to control what profiles are of interest
	// when sending data to pyroscope.
	//
	// Default is:
	// - cpu
	// - alloc_objects
	// - alloc_space
	// - inuse_objects
	// - inuse_space
	//
	// Additional profile types are:
	// - goroutines
	// - mutex_count
	// - mutex_duration
	// - block_count
	// - block_duration
	Profiles []pyroscope.ProfileType `mapstructure:"profiles"`
}

var (
	_ component.Config          = (*Config)(nil)
	_ component.ConfigValidator = (*Config)(nil)
	_ confmap.Unmarshaler       = (*Config)(nil)
)

var (
	validProfiles = map[pyroscope.ProfileType]struct{}{
		pyroscope.ProfileCPU:           {},
		pyroscope.ProfileInuseObjects:  {},
		pyroscope.ProfileAllocObjects:  {},
		pyroscope.ProfileInuseSpace:    {},
		pyroscope.ProfileAllocSpace:    {},
		pyroscope.ProfileGoroutines:    {},
		pyroscope.ProfileMutexCount:    {},
		pyroscope.ProfileMutexDuration: {},
		pyroscope.ProfileBlockCount:    {},
		pyroscope.ProfileBlockDuration: {},
	}
)

func (c *Config) Validate() (errs error) {
	if c.RuntimeBlockProfileFaction < 1 {
		errs = multierr.Append(
			errs,
			fmt.Errorf("runtime block profile fraction: %w", ErrRequiresPositiveValue),
		)
	}
	if c.RuntimeMutexProfileFraction < 1 {
		errs = multierr.Append(
			errs,
			fmt.Errorf("runtime mutex profile fraction: %w", ErrRequiresPositiveValue),
		)
	}
	if c.Endpoint == "" {
		errs = multierr.Append(errs, ErrRequiresEndpoint)
	}
	for _, p := range c.Profiles {
		if _, ok := validProfiles[p]; !ok {
			errs = multierr.Append(errs, fmt.Errorf("unknown profile %s: %w", p, ErrRequiresValidProfile))
		}
	}
	return errs
}

func (c *Config) Unmarshal(parser *confmap.Conf) error {
	if parser == nil {
		return nil
	}
	// If the user has set their own list of profiles
	// then the existing list should be cleared to avoid confusion.
	if parser.IsSet("profiles") {
		c.Profiles = c.Profiles[:0]
	}
	return parser.Unmarshal(c, confmap.WithErrorUnused())
}
