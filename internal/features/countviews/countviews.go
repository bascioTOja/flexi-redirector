package countviews

import (
	"errors"
	"time"

	"flexi-redirector/internal/env"
)

type Feature struct {
	enabled bool

	countAsync   bool          // Perform IncrementViews asynchronously.
	asyncTimeout time.Duration // Defines how long can max wait for IncrementViews.

	// batchUpdates bool              // TODO: implement batch updates
	// flushInterval: 5 * time.Second // TODO: implement batch updates
}

func New() *Feature {
	return &Feature{
		enabled:      true,
		countAsync:   true,
		asyncTimeout: 2 * time.Second,
	}
}

func (f *Feature) Name() string { return "countviews" }

func (f *Feature) Enabled() bool { return f.enabled }

func (f *Feature) Load() error {
	f.enabled = env.Bool("FEATURE_COUNT_VIEWS_ENABLED", true)
	f.countAsync = env.Bool("FEATURE_COUNT_VIEWS_COUNT_ASYNC", true)
	f.asyncTimeout = env.Duration("FEATURE_COUNT_VIEWS_ASYNC_TIMEOUT", 2*time.Second)
	return nil
}

func (f *Feature) Validate() error {
	if !f.enabled {
		return nil
	}
	if f.asyncTimeout <= 0 {
		return errors.New("FEATURE_COUNT_VIEWS_ASYNC_TIMEOUT must be > 0")
	}
	return nil
}

func (f *Feature) CountAsync() bool            { return f.countAsync }
func (f *Feature) AsyncTimeout() time.Duration { return f.asyncTimeout }
