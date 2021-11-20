package logreqid

import (
	"net/http"

	"github.com/eurofurence/reg-auth-service/internal/repository/logging"
	"github.com/eurofurence/reg-auth-service/web/filter/reqid"
)

func logRequestIdHandler(next http.Handler) func(w http.ResponseWriter, r *http.Request) {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		newCtx := logging.CreateContextWithLoggerForRequestId(ctx, reqid.GetRequestID(ctx))
		r = r.WithContext(newCtx)

		next.ServeHTTP(w, r)
	}
	return handlerFunc
}

// would not need this extra layer in the absence of parameters

func LogRequestIdMiddleware() func(http.Handler) http.Handler {
	middlewareCreator := func(next http.Handler) http.Handler {
		return http.HandlerFunc(logRequestIdHandler(next))
	}
	return middlewareCreator
}
