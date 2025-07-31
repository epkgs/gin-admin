package app

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof" //nolint:gosec
	"os"
	"os/signal"
	"syscall"
	"time"

	"gin-admin/internal/configs"
	_ "gin-admin/internal/swagger"
	"gin-admin/pkg/logger"
)

// The Run function initializes and starts a service with configuration and logger, and handles
// cleanup upon exit.
func Run(ctx context.Context, configFile string) error {
	defer func() {
		if err := logger.Sync(ctx); err != nil {
			fmt.Printf("failed to sync logger: %s \n", err.Error())
		}
	}()

	// Load configuration.
	configs.MustLoad(ctx, configFile)

	// Initialize logger.
	cleanLoggerFn, err := logger.InitWithConfig(ctx, &configs.C.Logger)
	if err != nil {
		return err
	}
	ctx = logger.WithTag(ctx, logger.Tag_Main)

	logger.Info(ctx, "starting service ...",
		map[string]any{
			"version": configs.C.Version,
			"pid":     os.Getpid(),
			"config":  configs.C.ConfigFile,
			"env":     configs.C.AppEnv,
			"static":  configs.C.HTTP.StaticDir,
		},
	)

	// Start pprof server.
	if addr := configs.C.Pprof.Addr; addr != "" {
		logger.Info(ctx, "pprof server is listening on "+addr)
		go func() {
			err := http.ListenAndServe(addr, nil)
			if err != nil {
				logger.Error(ctx, "failed to listen pprof server", err)
			}
		}()
	}

	app := New(ctx, configs.C)

	if err := app.Init(ctx); err != nil {
		return err
	}

	return run(ctx, func(ctx context.Context) (func(), error) {
		err := app.InitHttp(ctx)

		cleaner := func() {

			if cleanLoggerFn != nil {
				cleanLoggerFn()
			}

			if err := app.Release(ctx); err != nil {
				logger.Error(ctx, "failed to release app context", err)
			}
		}

		if err != nil {
			return cleaner, err
		}

		return cleaner, nil
	})
}

// The Run function sets up a signal handler and executes a handler function until a termination signal
// is received.
func run(ctx context.Context, handler func(ctx context.Context) (func(), error)) error {
	state := 1
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	cleanFn, err := handler(ctx)
	if err != nil {
		return err
	}

EXIT:
	for {
		sig := <-sc
		logger.Info(ctx, "Received signal", map[string]any{"signal": sig.String()})

		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			state = 0
			break EXIT
		case syscall.SIGHUP:
		default:
			break EXIT
		}
	}

	cleanFn()
	logger.Info(ctx, "Server exit, bye...")
	time.Sleep(time.Millisecond * 100)
	os.Exit(state)
	return nil
}
