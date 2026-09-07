package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const (
	readHeaderTimeout = 5 * time.Second
	readTimeout       = 10 * time.Second
	writeTimeout      = 10 * time.Second
	idleTimeout       = 60 * time.Second
)

func newHTTPServer(address string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:              address,
		Handler:           handler,
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}
}

func (a *App) Run(ctx context.Context) error {
	serverError := make(chan error, 1)
	go func() {
		a.logger.Info("HTTP server started", "address", a.httpServer.Addr)
		serverError <- a.httpServer.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		return a.shutdown()
	case err := <-serverError:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return fmt.Errorf("serve HTTP: %w", err)
	}
}

func (a *App) shutdown() error {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), a.shutdownTimeout)
	defer cancel()
	if err := a.httpServer.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown HTTP server: %w", err)
	}
	return nil
}
