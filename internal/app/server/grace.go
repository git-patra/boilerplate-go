package server

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	defaultGraceTimeout = 30 * time.Second
)

var (
	// ErrGraceShutdownTimeout happens when the server graceful shutdown exceed the given grace timeout.
	ErrGraceShutdownTimeout = fmt.Errorf("server shutdown timed out")
)

// HTTPServer represents an HTTP server
type HTTPServer interface {
	Shutdown(ctx context.Context) error
	Serve(l net.Listener) error
}

// ServeHTTP start the http server on the given address and add graceful shutdown handler
// graceTimeout specify how long we want to wait for the Shutdown to run.
// if graceTimeout = 0, we use default value: 30 second
func ServeHTTP(
	server HTTPServer,
	address string,
	graceTimeout time.Duration,
) error {

	// Create a network listener.
	listener, err := net.Listen("tcp4", address)
	if err != nil {
		return err
	}

	// This channel will be used to tell the goroutine to stop itself without waiting for the signal.
	errCh := make(chan error, 1)
	stoppedCh := WaitTerminateSignal(
		func(ctx context.Context) error {
			if graceTimeout == 0 {
				graceTimeout = defaultGraceTimeout
			}

			stopped := make(chan bool, 1)
			ctx, cancel := context.WithTimeout(ctx, graceTimeout)
			defer cancel()

			go func() {
				if err := server.Shutdown(ctx); err != nil {
					logrus.Error("HTTP server shutdown failed", nil)
				}
				stopped <- true
			}()

			select {
			case <-ctx.Done():
				return ErrGraceShutdownTimeout
			case <-stopped:
				return nil
			}
		},
		errCh,
	)

	logrus.Infof("starting HTTP server on %s", address)

	// Start the HTTP server.
	if err = server.Serve(listener); err != http.ErrServerClosed {
		errCh <- err
		// Don't return here, otherwise the stoppedCh won't be closed.
	}

	<-stoppedCh
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	logrus.Info("HTTP server stopped")
	return nil
}

// WaitTerminateSignal wait for a termination signal, then execute the given handler.
// It returns channel which will be closed after the signal received and the handler executed.
// We can use the signal to wait for the shutdown to be finished.
func WaitTerminateSignal(
	handler func(context.Context) error,
	errCh <-chan error,
) <-chan bool {

	stoppedCh := make(chan bool)

	go func() {
		signalCh := make(chan os.Signal, 1)

		// Wait for one of the signal.
		signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)

		select {
		case err := <-errCh:
			// Got an error from the caller.
			// Exit early.
			logrus.Error(err, "error when waiting for termination signal")
		case <-signalCh:
			// Signal received.
			if err := handler(context.Background()); err != nil {
				logrus.Error(err, "graceful shutdown failed")
			} else {
				logrus.Error("graceful shutdown succeed")
			}
		}

		stoppedCh <- true

	}()

	return stoppedCh
}
