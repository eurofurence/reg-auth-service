package web

import (
	"github.com/go-chi/chi/middleware"
	"net/http"
	"time"

	"github.com/eurofurence/reg-auth-service/internal/repository/config"
	"github.com/eurofurence/reg-auth-service/internal/repository/logging"
	"github.com/eurofurence/reg-auth-service/web/controller/authctl"
	"github.com/eurofurence/reg-auth-service/web/controller/dropoffctl"
	"github.com/eurofurence/reg-auth-service/web/controller/healthctl"
	"github.com/eurofurence/reg-auth-service/web/filter/corsfilter"
	"github.com/eurofurence/reg-auth-service/web/filter/logreqid"
	"github.com/eurofurence/reg-auth-service/web/filter/reqid"
	"github.com/go-chi/chi"
)

func Create() chi.Router {
	logging.NoCtx().Info("Building routers...")
	server := chi.NewRouter()

	server.Use(middleware.Timeout(60 * time.Second))
	server.Use(reqid.RequestIdMiddleware())
	server.Use(logreqid.LogRequestIdMiddleware())
	server.Use(corsfilter.CorsHeadersMiddleware())

	healthctl.Create(server)
	// add your controllers here
	authctl.Create(server)
	dropoffctl.Create(server)
	return server
}

func Serve(server chi.Router) {
	address := config.ServerAddr()
	logging.NoCtx().Info("Listening on " + address)
	err := http.ListenAndServe(address, server)
	if err != nil {
		logging.NoCtx().Error(err)
	}
}
