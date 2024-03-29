package app

import (
	"context"
	"errors"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	"github.com/StephanHCB/go-autumn-logging-zerolog/loggermiddleware"
	"github.com/eurofurence/reg-auth-service/internal/repository/config"
	"github.com/eurofurence/reg-auth-service/internal/repository/idp"
	"github.com/eurofurence/reg-auth-service/internal/web/controller/authctl"
	"github.com/eurofurence/reg-auth-service/internal/web/controller/dropoffctl"
	"github.com/eurofurence/reg-auth-service/internal/web/controller/healthctl"
	"github.com/eurofurence/reg-auth-service/internal/web/controller/logoutctl"
	"github.com/eurofurence/reg-auth-service/internal/web/controller/userinfoctl"
	"github.com/eurofurence/reg-auth-service/internal/web/middleware"
	"github.com/go-chi/chi/v5"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func CreateRouter(ctx context.Context) chi.Router {
	aulogging.Logger.NoCtx().Debug().Print("Setting up router")
	server := chi.NewRouter()

	server.Use(middleware.AddRequestIdToContextAndResponse)
	server.Use(loggermiddleware.AddZerologLoggerToContext)
	server.Use(middleware.RequestLogger)
	server.Use(middleware.PanicRecoverer)
	server.Use(middleware.CorsHandling)
	server.Use(middleware.TokenValidator)

	idpClient := idp.New()

	healthctl.Create(server)
	// add your controllers here
	authctl.Create(server)
	dropoffctl.Create(server, idpClient)
	userinfoctl.Create(server, idpClient)
	logoutctl.Create(server)
	return server
}

func newServer(ctx context.Context, router chi.Router) *http.Server {
	aulogging.Logger.NoCtx().Debug().Print("setting up server")
	return &http.Server{
		Addr:         config.ServerAddr(),
		Handler:      router,
		ReadTimeout:  config.ServerReadTimeout(),
		WriteTimeout: config.ServerWriteTimeout(),
		IdleTimeout:  config.ServerIdleTimeout(),
		BaseContext: func(l net.Listener) context.Context {
			return ctx
		},
	}
}

func runServerWithGracefulShutdown() error {
	ctx, cancel := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	handler := CreateRouter(ctx)
	srv := newServer(ctx, handler)

	go func() {
		<-sig
		defer cancel()
		aulogging.Logger.NoCtx().Debug().Print("Stopping services now")

		tCtx, tcancel := context.WithTimeout(ctx, time.Second*5)
		defer tcancel()

		if err := srv.Shutdown(tCtx); err != nil {
			aulogging.Logger.NoCtx().Error().WithErr(err).Printf("Couldn't shutdown server gracefully: %s", err.Error())
			os.Exit(3)
		}
	}()

	aulogging.Logger.NoCtx().Info().Print("Running service on ", config.ServerAddr())
	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		aulogging.Logger.NoCtx().Error().WithErr(err).Printf("Server closed unexpectedly: %s", err.Error())
		return err
	}

	return nil
}
