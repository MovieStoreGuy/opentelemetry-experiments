package pyroscopeextension

import (
	"context"
	"errors"
	"runtime"

	"github.com/pyroscope-io/client/pyroscope"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/extension"
	"go.uber.org/zap"
)

var (
	ErrAlreadyStarted = errors.New("already started")
	ErrNeverStarted   = errors.New("never started")
)

type pyrofiler struct {
	conf *Config
	sem  chan struct{}
	log  *zap.Logger

	p *pyroscope.Profiler
}

var (
	_ extension.Extension = (*pyrofiler)(nil)
)

func newPyroscopeProfiler(ctx context.Context, set extension.CreateSettings, cfg component.Config) (extension.Extension, error) {
	return &pyrofiler{
		conf: cfg.(*Config),
		sem:  make(chan struct{}, 1),
		log:  set.Logger,
	}, nil
}

func (py *pyrofiler) Start(ctx context.Context, host component.Host) error {
	select {
	case py.sem <- struct{}{}:
		// Aquired Sem
	default:
		return ErrAlreadyStarted
	}

	runtime.SetBlockProfileRate(py.conf.RuntimeBlockProfileFaction)
	runtime.SetMutexProfileFraction(py.conf.RuntimeMutexProfileFraction)

	p, err := pyroscope.Start(pyroscope.Config{
		ApplicationName: py.conf.ApplicationName,
		Tags:            py.conf.Tags,
		ServerAddress:   py.conf.Endpoint,
		AuthToken:       string(py.conf.AuthToken),
		ProfileTypes:    py.conf.Profiles,
		Logger:          py.log.Sugar(),
	})

	py.p = p

	return err
}

func (py *pyrofiler) Shutdown(_ context.Context) error {
	select {
	case <-py.sem:
		// Released Sem
	default:
		return ErrNeverStarted
	}
	return py.p.Stop()
}
