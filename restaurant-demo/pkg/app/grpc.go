package app

import (
	"context"
	"fmt"
	"net"
	"time"

	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func (a *App) Run(ctx context.Context) error {
	listener, err := net.Listen("tcp", a.grpcAddress)
	if err != nil {
		return fmt.Errorf("listen gRPC: %w", err)
	}

	serverError := make(chan error, 1)
	go func() {
		a.logger.Info("restaurant gRPC server started", "address", a.grpcAddress)
		serverError <- a.grpcServer.Serve(listener)
	}()

	select {
	case <-ctx.Done():
		return a.shutdown()
	case err := <-serverError:
		if err != nil {
			return fmt.Errorf("serve gRPC: %w", err)
		}
		return nil
	}
}

func (a *App) shutdown() error {
	a.grpcHealth.SetServingStatus("", healthpb.HealthCheckResponse_NOT_SERVING)

	stopped := make(chan struct{})
	go func() {
		a.grpcServer.GracefulStop()
		close(stopped)
	}()

	timer := time.NewTimer(a.shutdownTimeout)
	defer timer.Stop()
	select {
	case <-stopped:
		return nil
	case <-timer.C:
		a.grpcServer.Stop()
		return fmt.Errorf("gRPC graceful shutdown timeout")
	}
}
