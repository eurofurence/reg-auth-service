package logoutctl

import (
	"context"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	"github.com/eurofurence/reg-auth-service/internal/web/controller"
	"net/http"
	"time"

	"github.com/eurofurence/reg-auth-service/internal/repository/config"
	"github.com/go-chi/chi/v5"
)

func Create(server chi.Router) {
	server.Get("/v1/logout", logoutHandler)
}

/* Handle /logout requests.
 *
 * Required parameters are:
 *  * app_name  - the name of the application that the user wants to be authenticated for
 *
 * Redirects to app_name's default dropoff url after cookie deletion.
 */
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	query := r.URL.Query()
	regAppName := query.Get("app_name")
	if regAppName == "" {
		logoutErrorHandler(ctx, w, regAppName, http.StatusBadRequest, "app_name parameter is missing", "invalid parameters")
		return
	}

	applicationConfig, err := config.GetApplicationConfig(regAppName)
	if err != nil {
		logoutErrorHandler(ctx, w, regAppName, http.StatusNotFound, "app_name is unknown", "invalid parameters")
		return
	}

	err = clearCookieAndRedirectToDropOffUrl(ctx, w, applicationConfig)
	if err != nil {
		logoutErrorHandler(ctx, w, regAppName, http.StatusInternalServerError, err.Error(), "internal error")
		return
	}
	aulogging.Logger.Ctx(ctx).Info().Printf("OK v1/logout(%s)-> %d", regAppName, http.StatusFound)
}

func logoutErrorHandler(ctx context.Context, w http.ResponseWriter, appName string, status int, logMsg string, publicMsg string) {
	aulogging.Logger.Ctx(ctx).Warn().Printf("FAIL v1/logout(%s) -> %d: %s", appName, status, logMsg)
	w.WriteHeader(status)
	_, _ = w.Write(controller.ErrorResponse(ctx, publicMsg, ""))
}

func clearCookieAndRedirectToDropOffUrl(ctx context.Context, w http.ResponseWriter, applicationConfig config.ApplicationConfig) error {
	cookie := &http.Cookie{
		Name:     applicationConfig.CookieName,
		Value:    "",
		Domain:   applicationConfig.CookieDomain,
		Expires:  time.Now(),
		MaxAge:   -1,
		Path:     applicationConfig.CookiePath,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)

	if config.OidcAccessTokenCookieName() != "" {
		// additional cookie needed for this service
		accessCookie := &http.Cookie{
			Name:     config.OidcAccessTokenCookieName(),
			Value:    "",
			Domain:   applicationConfig.CookieDomain,
			Expires:  time.Now(),
			MaxAge:   -1,
			Path:     applicationConfig.CookiePath,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
		}
		http.SetCookie(w, accessCookie)
	}

	w.Header().Set("Location", applicationConfig.DefaultDropoffUrl)
	w.WriteHeader(http.StatusFound)
	return nil
}
