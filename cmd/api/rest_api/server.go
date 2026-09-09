package restapi

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lbenedar/fresta/cmd/api/core"
)

type AppServer core.Application

func ListenAndServe(app *AppServer) error {
	srv := &http.Server{
		Addr:     app.Cfg.Port,
		Handler:  app.routes(),
		ErrorLog: slog.NewLogLogger(app.Slogger.Handler(), slog.LevelError),
		// TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	shutdownError := make(chan error)
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		app.Slogger.Info("Caught signal", "signal", s.String())

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		err = app.FoundryApp.Shutdown()
		if err != nil {
			shutdownError <- err
		}

		app.FoundryApp.Transport.DB.Close()

		shutdownError <- nil
	}()

	app.Slogger.Info("Starting server", "addr", srv.Addr, "env", app.Cfg.Env)
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	app.Slogger.Info("Stopped server", "addr", srv.Addr)
	// err = <-shutdownError
	if err != nil {
		return err
	}
	return nil
}
