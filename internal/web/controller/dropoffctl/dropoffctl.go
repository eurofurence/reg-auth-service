package dropoffctl

import (
	"context"
	"fmt"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	"github.com/eurofurence/reg-auth-service/internal/web/controller"
	"net/http"
	"time"

	"github.com/eurofurence/reg-auth-service/internal/entity"
	"github.com/eurofurence/reg-auth-service/internal/repository/config"
	"github.com/eurofurence/reg-auth-service/internal/repository/database"
	"github.com/eurofurence/reg-auth-service/internal/repository/idp"
	"github.com/go-chi/chi/v5"
)

var IDPClient idp.IdentityProviderClient

func Create(server chi.Router, idpClient idp.IdentityProviderClient) {
	if IDPClient == nil {
		IDPClient = idpClient
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
		dropOffErrorHandler(ctx, w, state, http.StatusBadRequest, "state parameter is missing", "invalid parameters", config.ErrorUrl())
		return
	}

	authRequest, err := database.GetRepository().GetAuthRequestByState(ctx, state)
	if err != nil {
		dropOffErrorHandler(ctx, w, state, http.StatusNotFound, "couldn't load auth request: "+err.Error(), "auth request not found or timed out", config.ErrorUrl())
		return
	}

	applicationConfig, err := config.GetApplicationConfig(authRequest.Application)
	if err != nil {
		dropOffErrorHandler(ctx, w, state, http.StatusInternalServerError, "couldn't load application config: "+err.Error(), "internal error", config.ErrorUrl())
		return
	}

	errorCode := query.Get("error")
	errorDescription := query.Get("error_description")
	if errorCode != "" {
		dropOffErrorHandler(ctx, w, state, http.StatusBadRequest, fmt.Sprintf("error parameter set (%s|%s)", errorCode, errorDescription), fmt.Sprintf("%s: %s", errorCode, errorDescription), applicationConfig.DefaultDropoffUrl)
		return
	}

	authCode := query.Get("code")
	if authCode == "" {
		dropOffErrorHandler(ctx, w, state, http.StatusBadRequest, "authorization_code parameter is missing", "invalid parameters", config.ErrorUrl())
		return
	}

	idToken, accessToken, httpstatus, err := fetchToken(ctx, authCode, *authRequest)
	if err != nil {
		dropOffErrorHandler(ctx, w, state, httpstatus, "couldn't fetch access codes: "+err.Error(), "failed to fetch token", config.ErrorUrl())
		return
	}

	err = setCookiesAndRedirectToDropOffUrl(ctx, w, idToken, accessToken, *authRequest, applicationConfig)
	if err != nil {
		dropOffErrorHandler(ctx, w, state, http.StatusInternalServerError, err.Error(), "internal error", "")
		return
	}
	aulogging.Logger.Ctx(ctx).Info().Printf("OK dropoff[%s]", state)
}

func dropOffErrorHandler(ctx context.Context, w http.ResponseWriter, state string, status int, logMsg string, publicMsg string, retryUrl string) {
	aulogging.Logger.Ctx(ctx).Warn().Printf("FAIL dropoff[%s]: %s", state, logMsg)
	w.WriteHeader(status)
	_, _ = w.Write(controller.ErrorResponse(ctx, publicMsg, retryUrl))
}

func fetchToken(ctx context.Context, authCode string, ar entity.AuthRequest) (string, string, int, error) {
	response, httpstatus, err := IDPClient.TokenWithAuthenticationCodeAndPKCE(ctx, ar.Application, authCode, ar.PkceCodeVerifier)
	if err != nil {
		return "", "", httpstatus, err
	}
	return response.IdToken, response.AccessToken, httpstatus, nil
}

func setCookiesAndRedirectToDropOffUrl(ctx context.Context, w http.ResponseWriter, idToken string, accessToken string, authRequest entity.AuthRequest, applicationConfig config.ApplicationConfig) error {
	sameSite := http.SameSiteStrictMode
	httpOnly := true // https://stackoverflow.com/questions/71819265/httponly-cookie-and-fetch
	secure := true
	if config.SendInsecureCookies() {
		secure = false
		aulogging.Logger.Ctx(ctx).Warn().Print("sending insecure cookies. This configuration is not intended for production use, only for local development!")
	}
	if config.IsCorsDisabled() {
		sameSite = http.SameSiteNoneMode
		aulogging.Logger.Ctx(ctx).Warn().Print("sending Same Site Policy None cookies to work with disabled CORS. This configuration is not intended for production use, only for local development!")
	}
	if config.SendNonHttpOnlyCookies() {
		httpOnly = false
		aulogging.Logger.Ctx(ctx).Warn().Print("sending non-http-only cookies. This configuration is not intended for production use, only for local development!")
	}

	// first set the cookie wanted by the application
	applicationCookie := &http.Cookie{
		Name:     applicationConfig.CookieName,
		Value:    idToken,
		Domain:   applicationConfig.CookieDomain,
		Expires:  time.Now().Add(applicationConfig.CookieExpiry),
		Path:     applicationConfig.CookiePath,
		Secure:   secure,
		HttpOnly: httpOnly,
		SameSite: sameSite,
	}
	http.SetCookie(w, applicationCookie)

	if config.OidcAccessTokenCookieName() != "" {
		// additional cookie needed for this service
		accessCookie := &http.Cookie{
			Name:     config.OidcAccessTokenCookieName(),
			Value:    accessToken,
			Domain:   applicationConfig.CookieDomain,
			Expires:  time.Now().Add(applicationConfig.CookieExpiry),
			Path:     applicationConfig.CookiePath,
			Secure:   secure,
			HttpOnly: httpOnly,
			SameSite: sameSite,
		}
		http.SetCookie(w, accessCookie)
	}

	w.Header().Set("Location", authRequest.DropOffUrl)
	w.WriteHeader(http.StatusFound)
	return nil
}
