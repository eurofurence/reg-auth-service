package dropoffctl

import (
	"context"
	"net/http"
	"time"

	"github.com/eurofurence/reg-auth-service/internal/entity"
	"github.com/eurofurence/reg-auth-service/internal/repository/config"
	"github.com/eurofurence/reg-auth-service/internal/repository/database"
	"github.com/eurofurence/reg-auth-service/internal/repository/logging"
	"github.com/go-chi/chi"
)

func Create(server chi.Router) {
	server.Get("/v1/dropoff", dropOffHandler)
}

/* Handle /dropoff requests.
 *
 * Required parameters are:
 *  * state              - random-string identifier of this flow
 *  * authorization_code - temporary credential to obtain the access token from the OIDC
 */
func dropOffHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	query := r.URL.Query()
	state := query.Get("state")
	if state == "" {
		dropOffErrorHandler(ctx, w, http.StatusBadRequest, "state parameter is missing")
		return
	}
	authCode := query.Get("authorization_code")
	if authCode == "" {
		dropOffErrorHandler(ctx, w, http.StatusBadRequest, "authorization_code parameter is missing")
		return
	}

	authRequest, err := database.GetRepository().GetAuthRequestByState(ctx, state)
	if err != nil {
		dropOffErrorHandler(ctx, w, http.StatusNotFound, "couldn't load auth request: " + err.Error())
		return
	}

	applicationConfig, err := config.GetApplicationConfig(authRequest.Application)
	if err != nil {
		dropOffErrorHandler(ctx, w, http.StatusInternalServerError, "couldn't load application config: " + err.Error())
		return
	}

	accessCode, err := fetchAccessCode(ctx, *authRequest, applicationConfig)
	if err != nil {
		dropOffErrorHandler(ctx, w, http.StatusInternalServerError, "couldn't fetch access code: " + err.Error())
		return
	}

	err = setCookieAndRedirectToDropOffUrl(ctx, w, accessCode, *authRequest, applicationConfig)
	if err != nil {
		dropOffErrorHandler(ctx, w, http.StatusInternalServerError, err.Error())
	}
}

func dropOffErrorHandler(ctx context.Context, w http.ResponseWriter, status int, msg string) {
	logging.Ctx(ctx).Error(msg)
	// TODO: here we should display some information to the user
	w.WriteHeader(status)
}

func fetchAccessCode(ctx context.Context, ar entity.AuthRequest, ac config.ApplicationConfig) (string, error) {
	return "dummy_mock_value", nil
}

func setCookieAndRedirectToDropOffUrl(ctx context.Context, w http.ResponseWriter, accessCode string, authRequest entity.AuthRequest, applicationConfig config.ApplicationConfig) error {
	cookie := &http.Cookie{
		Name:  "AccessCode",                     // make this configurable?
		Value: accessCode,
		Domain: "example.com",                   // make this configurable?
		Expires: time.Now().Add(6 * time.Hour),  // make this configurable?
		Secure: true,                            // make this configurable?
		SameSite: http.SameSiteStrictMode,       // make this configurable?
	}
	http.SetCookie(w, cookie)
	w.Header().Set("Location", authRequest.DropOffUrl)
	w.WriteHeader(http.StatusFound)
	return nil
}
