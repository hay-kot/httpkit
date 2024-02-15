// Package graceful provides a graceful start and shutdown insturmentation for
// go programs based on a simple Plugin interface.
//
// It supports
//   - graceful shutdown of plugins
//   - timeout for plugins to shutdown
//   - os.Signals to listen for (defaults to os.Interrupt and syscall.SIGTERM)
package graceful
