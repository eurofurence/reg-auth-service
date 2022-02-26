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
	"github.com/eurofurence/reg-auth-service/web/controller"
	"github.com/go-chi/chi/v5"
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
		dropOffErrorHandler(ctx, w, state, http.StatusBadRequest, "state parameter is missing", "invalid parameters")
		return
	}
	authCode := query.Get("code")
	if authCode == "" {
		dropOffErrorHandler(ctx, w, state, http.StatusBadRequest, "authorization_code parameter is missing", "invalid parameters")
		return
	}

	authRequest, err := database.GetRepository().GetAuthRequestByState(ctx, state)
	if err != nil {
		dropOffErrorHandler(ctx, w, state, http.StatusNotFound, "couldn't load auth request: "+err.Error(), "auth request not found or timed out")
		return
	}

	applicationConfig, err := config.GetApplicationConfig(authRequest.Application)
	if err != nil {
		dropOffErrorHandler(ctx, w, state, http.StatusInternalServerError, "couldn't load application config: "+err.Error(), "internal error")
		return
	}

	accessCode, httpstatus, err := fetchToken(ctx, authCode, *authRequest)
	if err != nil {
		dropOffErrorHandler(ctx, w, state, httpstatus, "couldn't fetch access code: "+err.Error(), "failed to fetch token")
		return
	}

	err = setCookieAndRedirectToDropOffUrl(ctx, w, accessCode, *authRequest, applicationConfig)
	if err != nil {
		dropOffErrorHandler(ctx, w, state, http.StatusInternalServerError, err.Error(), "internal error")
		return
	}
	logging.Ctx(ctx).Info(fmt.Sprintf("OK v1/dropoff(%s)-> %d", state, http.StatusFound))
}

func dropOffErrorHandler(ctx context.Context, w http.ResponseWriter, state string, status int, logMsg string, publicMsg string) {
	logging.Ctx(ctx).Warn(fmt.Sprintf("FAIL v1/dropoff(%s) -> %d: %s", state, status, logMsg))
	w.WriteHeader(status)
	_, _ = w.Write(controller.ErrorResponse(ctx, publicMsg))
}

func fetchToken(ctx context.Context, authCode string, ar entity.AuthRequest) (string, int, error) {
	response, httpstatus, err := IDPClient.TokenWithAuthenticationCodeAndPKCE(ctx, ar.Application, authCode, ar.PkceCodeVerifier)
	if err != nil {
		return "", httpstatus, err
	}
	return response.IdToken, httpstatus, nil
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
