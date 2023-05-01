package pyroscopeextension

import (
	"github.com/pyroscope-io/client/pyroscope"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/extension"
)

const (
	typeStr   = "pyroscope"
	stability = component.StabilityLevelBeta
)

func NewFactory() extension.Factory {
	return extension.NewFactory(
		typeStr,
		newDefaultConfig,
		newPyroscopeProfiler,
		stability,
	)
}

func newDefaultConfig() component.Config {
	return &Config{
		RuntimeMutexProfileFraction: 5,
		RuntimeBlockProfileFraction: 5,
		ApplicationName:             "open-telemetry-collector",
		Profiles: append(
			[]pyroscope.ProfileType{},
			pyroscope.DefaultProfileTypes...,
		),
	}
}
