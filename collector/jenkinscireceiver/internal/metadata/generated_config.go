// Code generated by mdatagen. DO NOT EDIT.

package metadata

import "go.opentelemetry.io/collector/confmap"

// MetricConfig provides common config for a particular metric.
type MetricConfig struct {
	Enabled bool `mapstructure:"enabled"`

	enabledSetByUser bool
}

func (ms *MetricConfig) Unmarshal(parser *confmap.Conf) error {
	if parser == nil {
		return nil
	}
	err := parser.Unmarshal(ms, confmap.WithErrorUnused())
	if err != nil {
		return err
	}
	ms.enabledSetByUser = parser.IsSet("enabled")
	return nil
}

// MetricsConfig provides config for jenkins metrics.
type MetricsConfig struct {
	JenkinsJobCommitDelta MetricConfig `mapstructure:"jenkins.job.commit_delta"`
	JenkinsJobDuration    MetricConfig `mapstructure:"jenkins.job.duration"`
	JenkinsJobsCount      MetricConfig `mapstructure:"jenkins.jobs.count"`
}

func DefaultMetricsConfig() MetricsConfig {
	return MetricsConfig{
		JenkinsJobCommitDelta: MetricConfig{
			Enabled: true,
		},
		JenkinsJobDuration: MetricConfig{
			Enabled: true,
		},
		JenkinsJobsCount: MetricConfig{
			Enabled: true,
		},
	}
}

// MetricsBuilderConfig is a configuration for jenkins metrics builder.
type MetricsBuilderConfig struct {
	Metrics MetricsConfig `mapstructure:"metrics"`
}

func DefaultMetricsBuilderConfig() MetricsBuilderConfig {
	return MetricsBuilderConfig{
		Metrics: DefaultMetricsConfig(),
	}
}
