package features

import "fmt"

// Module is a self-contained feature module.
// Each module loads its own env config, validates it and provides a runtime handle.
//
// - Name is used for errors/logging.
// - Enabled should be a cheap check.
// - Load is called once during startup.
// - Validate must ensure the module is safe to use when enabled.
type Module interface {
	Name() string
	Enabled() bool
	Load() error
	Validate() error
}

// Manager loads and validates all registered feature modules.
type Manager struct {
	modules []Module
}

func NewManager(modules ...Module) *Manager {
	return &Manager{modules: modules}
}

func (manager *Manager) LoadAndValidate() error {
	for _, mod := range manager.modules {
		if err := mod.Load(); err != nil {
			return fmt.Errorf("feature %s: load failed: %w", mod.Name(), err)
		}
		if err := mod.Validate(); err != nil {
			return fmt.Errorf("feature %s: invalid config: %w", mod.Name(), err)
		}
	}
	return nil
}
