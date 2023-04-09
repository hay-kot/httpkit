# HttpKit

A tiny HTTP server toolkit for Go. This module is a collection of packages that provide an opinionated approach to building HTTP servers in Go while still allowing you to use your favorite router mux, and middleware.

## Packages

### server

The server package provides an encapsulated HTTP server with support for Graceful Shutdown, and Background Tasks. This is a bring your own mux approach, so you can use any router mux you want.

It also includes some useful helpers when working with HTTP requests and responses.

- Middleware
  - StripTrailingSlash
- JSON Response helper
- Decode JSON (Strict and Non-Strict)
- Signal Shutdown error

### errchain

The errchain package provides a simple way to handle errors in a chain of handler functions using a custom http.Handler interface that returns an error.

Errchain implements a `ToHandler` method that transforms the custom handler into a standard http.Handler allowing you to mix and match custom handlers with standard http.Handlers and ensures that the custom handler is always compatible with the standard http.Handler interface.

### errtrace

The errtrace packages is intended to work in conjunction with the errchain package. One of the problems with the errchain package is that it can be difficult to trace the error back to the original handler function or core service level function. The errtrace package provides a way to define Traceable errors that provide contextual information like:

- File path (relative or custom)
- Caller function name
- Caller function line number
- Additional message context
- Full wrapped error chain
- Human readable error trace (think stacktrace)

#### Error Trace Examples

`error:` defines a generic error type and not a traceable error.

`trace error:` defines a traceable error type which provides additional context.

_Note: real terminal output is colorized for readability_

```
error:
  failed to do something: error creating user in database
trace error:
  error creating user in database
  cmd/cli/main.go:22  main.ServiceNewUser() -> user repo: error writing to database
trace error:
  error writing to database
  cmd/cli/main.go:12  main.CreateUser() -> wrap: cmd/cli/main.go:10 main.CreateUser: user with id 1 already exists
error:
  cmd/cli/main.go:10 main.CreateUser: user with id 1 already exists
```