package server

import (
	"net/http"
	"time"
)

type Option = func(s *Server) error

// WithServerOverride specifies a custom http.Server to use. In most
// production environments, you will want to use this option.
// This will override the other options such as
//   - WithHost
//   - WithPort
//   - WithReadTimeout
//   - WithWriteTimeout
//   - WithIdleTimeout
func WithServerOverride(svr *http.Server) Option {
	return func(s *Server) error {
		s.activeServer = svr
		return nil
	}
}

func WithHost(host string) Option {
	return func(s *Server) error {
		s.Host = host
		return nil
	}
}

func WithPrintln(fn func(...any)) Option {
	return func(s *Server) error {
		s.println = fn
		return nil
	}
}

func WithPort(port string) Option {
	return func(s *Server) error {
		s.Port = port
		return nil
	}
}

func WithReadTimeout(seconds int) Option {
	return func(s *Server) error {
		s.readTimeout = time.Duration(seconds) * time.Second
		return nil
	}
}

func WithWriteTimeout(seconds int) Option {
	return func(s *Server) error {
		s.writeTimeout = time.Duration(seconds) * time.Second
		return nil
	}
}

func WithIdleTimeout(seconds int) Option {
	return func(s *Server) error {
		s.idleTimeout = time.Duration(seconds) * time.Second
		return nil
	}
}
