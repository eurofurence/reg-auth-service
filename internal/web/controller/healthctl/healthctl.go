package healthctl

import (
	"context"
	"encoding/json"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	"github.com/eurofurence/reg-auth-service/internal/web/util/media"
	"net/http"

	"github.com/eurofurence/reg-auth-service/internal/api/v1/health"
	"github.com/go-chi/chi/v5"
	"github.com/go-http-utils/headers"
)

func Create(server chi.Router) {
	server.Get("/", healthGet)
}

func healthGet(w http.ResponseWriter, r *http.Request) {
	dto := health.HealthResultDto{Status: "up"}

	w.Header().Add(headers.ContentType, media.ContentTypeApplicationJson)
	w.WriteHeader(http.StatusOK)
	writeJson(r.Context(), w, dto)
}

func writeJson(ctx context.Context, w http.ResponseWriter, v interface{}) {
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(v)
	if err != nil {
		aulogging.Logger.Ctx(ctx).Warn().WithErr(err).Printf("error while encoding json response: %s", err.Error())
	}
}
