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
	grpcListener, err := a.listenGRPC()
	if err != nil {
		return err
	}

	serverError := make(chan error, 2)
	go func() {
		a.logger.Info("HTTP server started", "address", a.httpServer.Addr)
		err := a.httpServer.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			err = nil
		}
		if err != nil {
			err = fmt.Errorf("serve HTTP: %w", err)
		}
		serverError <- err
	}()
	a.serveGRPC(grpcListener, serverError)

	var runError error
	select {
	case <-ctx.Done():
	case runError = <-serverError:
	}

	return errors.Join(runError, a.shutdown())
}

func (a *App) shutdown() error {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), a.shutdownTimeout)
	defer cancel()

	httpError := a.httpServer.Shutdown(shutdownCtx)
	if httpError != nil {
		httpError = fmt.Errorf("shutdown HTTP server: %w", httpError)
	}
	return errors.Join(httpError, a.shutdownGRPC())
}
