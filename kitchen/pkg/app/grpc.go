package app

import (
	"errors"
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func (a *App) listenGRPC() (net.Listener, error) {
	listener, err := net.Listen("tcp", a.grpcAddress)
	if err != nil {
		return nil, fmt.Errorf("listen gRPC: %w", err)
	}
	return listener, nil
}

func (a *App) serveGRPC(listener net.Listener, serverError chan<- error) {
	go func() {
		a.logger.Info("gRPC server started", "address", a.grpcAddress)
		err := a.grpcServer.Serve(listener)
		if errors.Is(err, grpc.ErrServerStopped) {
			err = nil
		}
		if err != nil {
			err = fmt.Errorf("serve gRPC: %w", err)
		}
		serverError <- err
	}()
}

func (a *App) shutdownGRPC() error {
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
