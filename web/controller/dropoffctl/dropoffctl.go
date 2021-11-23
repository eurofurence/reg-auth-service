package dropoffctl

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/eurofurence/reg-auth-service/internal/entity"
	"github.com/eurofurence/reg-auth-service/internal/repository/config"
	"github.com/eurofurence/reg-auth-service/internal/repository/database"
	"github.com/eurofurence/reg-auth-service/internal/repository/idp"
	"github.com/eurofurence/reg-auth-service/internal/repository/idp/idpclient"
	"github.com/eurofurence/reg-auth-service/internal/repository/logging"
	"github.com/go-chi/chi"
)

var IDPClient idp.IdentityProviderClient

func Create(server chi.Router) {
	if IDPClient == nil {
		IDPClient = idpclient.New()
	}
	server.Get("/v1/dropoff", dropOffHandler)
}

/* Handle /dropoff requests.
 *
 * Required parameters are:
 *  * state - random-string identifier of this flow
 *  * code  - temporary credential to obtain the access token from the OIDC
 */
func dropOffHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	query := r.URL.Query()
	state := query.Get("state")
	if state == "" {
		dropOffErrorHandler(ctx, w, state, http.StatusBadRequest, "state parameter is missing")
		return
	}
	authCode := query.Get("code")
	if authCode == "" {
		dropOffErrorHandler(ctx, w, state, http.StatusBadRequest, "authorization_code parameter is missing")
		return
	}

	authRequest, err := database.GetRepository().GetAuthRequestByState(ctx, state)
	if err != nil {
		dropOffErrorHandler(ctx, w, state, http.StatusNotFound, "couldn't load auth request: " + err.Error())
		return
	}

	applicationConfig, err := config.GetApplicationConfig(authRequest.Application)
	if err != nil {
		dropOffErrorHandler(ctx, w, state, http.StatusInternalServerError, "couldn't load application config: " + err.Error())
		return
	}

	accessCode, err := fetchAccessCode(ctx, authCode, *authRequest, applicationConfig)
	if err != nil {
		dropOffErrorHandler(ctx, w, state, http.StatusInternalServerError, "couldn't fetch access code: " + err.Error())
		return
	}

	err = setCookieAndRedirectToDropOffUrl(ctx, w, accessCode, *authRequest, applicationConfig)
	if err != nil {
		dropOffErrorHandler(ctx, w, state, http.StatusInternalServerError, err.Error())
		return
	}
	logging.Ctx(ctx).Info(fmt.Sprintf("OK v1/dropoff(%s)-> %d", state, http.StatusFound))
}

func dropOffErrorHandler(ctx context.Context, w http.ResponseWriter, state string, status int, msg string) {
	logging.Ctx(ctx).Warn(fmt.Sprintf("FAIL v1/dropoff(%s) -> %d: %s", state, status, msg))
	// TODO: here we should display some information to the user
	w.WriteHeader(status)
}

func fetchAccessCode(ctx context.Context, authCode string, ar entity.AuthRequest, ac config.ApplicationConfig) (string, error) {
	response, err:= IDPClient.TokenWithAuthenticationCodeAndPKCE(ctx, ar.Application, authCode, ar.PkceCodeVerifier)
	if err != nil  {
		return "", err
	}
	return response.AccessToken, nil
}

func setCookieAndRedirectToDropOffUrl(ctx context.Context, w http.ResponseWriter, accessCode string, authRequest entity.AuthRequest, applicationConfig config.ApplicationConfig) error {
	cookie := &http.Cookie{
		Name:     applicationConfig.CookieName,
		Value:    accessCode,
		Domain:   applicationConfig.CookieDomain,
		Expires:  time.Now().Add(applicationConfig.CookieExpiry),
		Path:     applicationConfig.CookiePath,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)
	w.Header().Set("Location", authRequest.DropOffUrl)
	w.WriteHeader(http.StatusFound)
	return nil
}
