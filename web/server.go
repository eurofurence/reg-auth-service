package web

import (
	"github.com/eurofurence/reg-backend-template-test/internal/repository/config"
	"github.com/eurofurence/reg-backend-template-test/internal/repository/logging"
	"github.com/eurofurence/reg-backend-template-test/web/controller/healthctl"
	"github.com/eurofurence/reg-backend-template-test/web/filter/corsfilter"
	"github.com/eurofurence/reg-backend-template-test/web/filter/logreqid"
	"github.com/eurofurence/reg-backend-template-test/web/filter/reqid"
	"github.com/go-chi/chi"
	"net/http"
)

func Create() chi.Router {
	logging.NoCtx().Info("Building routers...")
	server := chi.NewRouter()

	server.Use(reqid.RequestIdMiddleware())
	server.Use(logreqid.LogRequestIdMiddleware())
	server.Use(corsfilter.CorsHeadersMiddleware())

	healthctl.Create(server)
	// add your controllers here
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
