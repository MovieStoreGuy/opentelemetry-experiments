package cicdreceiver

import (
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/receiver/scraperhelper"

	"github.com/MovieStoreGuy/collector/cicdreceiver/internal/metadata"
)

type Config struct {
	scraperhelper.ScraperControllerSettings `mapstructure:",squash"`
	confighttp.HTTPClientSettings           `mapstructure:",squash"`
	metadata.MetricsBuilderConfig           `mapstructure:",squash"`
}

func newDefaultConfig() component.Config {
	return &Config{
		ScraperControllerSettings: scraperhelper.NewDefaultScraperControllerSettings(metadata.Type),
		HTTPClientSettings:        confighttp.NewDefaultHTTPClientSettings(),
		MetricsBuilderConfig:      metadata.DefaultMetricsBuilderConfig(),
	}
}
