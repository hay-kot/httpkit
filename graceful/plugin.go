package graceful

import "context"

// Plugin defines the interface for a plugin that can be started and stopped via the
// Runner. The plugin should be blocking until the context is cancelled. Note that startup
// order for plugins is non-deterministic and should not be relied upon. If you need to
// will need to manage coordination between plugins.
type Plugin interface {
	// Name returns the name of the plugin. This can be any identifier you like and
	// will only be used for logging and debugging purposes.
	Name() string

	// Start starts the plugin with the context provided for cancellation.
	// if the context is cancelled, the plugin should stop, cleanup and return.
	// The start method _should_ be blocking, until the context is cancelled.
	// If the plugin terminates early with an error, it will cause a shutdown
	// in the system.
	//
	// Errors during shutdown should be logged within your plugin. Errors
	// that occur during the shutdown process will be ignored by the server.
	//
	// Example:
	//   func (p *MyPlugin) Start(ctx context.Context) error {
	//     // do something
	//     <-ctx.Done()
	//     // cleanup / logging
	//     return nil
	//   }
	Start(ctx context.Context) error
}

type pluginFunc struct {
	name  string
	start func(ctx context.Context) error
}

func (p *pluginFunc) Name() string {
	return p.name
}

func (p *pluginFunc) Start(ctx context.Context) error {
	return p.start(ctx)
}

func PluginFunc(name string, start func(ctx context.Context) error) Plugin {
	return &pluginFunc{name: name, start: start}
}
