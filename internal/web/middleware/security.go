package middleware

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	"github.com/eurofurence/reg-auth-service/internal/api/v1/errorapi"
	"github.com/eurofurence/reg-auth-service/internal/repository/config"
	"github.com/eurofurence/reg-auth-service/internal/web/util/ctxvalues"
	"github.com/eurofurence/reg-auth-service/internal/web/util/media"
	"github.com/go-http-utils/headers"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// --- getting the values from the request ---

func fromCookie(r *http.Request, cookieName string) string {
	if cookieName == "" {
		// ok if not configured, don't accept cookies then
		return ""
	}

	authCookie, _ := r.Cookie(cookieName)
	if authCookie == nil {
		// missing cookie is not considered an error, either
		return ""
	}

	return authCookie.Value
}

func fromAuthHeader(r *http.Request) string {
	headerValue := r.Header.Get(headers.Authorization)

	if !strings.HasPrefix(headerValue, "Bearer ") {
		return ""
	}

	return strings.TrimPrefix(headerValue, "Bearer ")
}

// --- validating the individual pieces ---

// important - if any of these return an error, you must abort processing via "return" and log the error message

func recordAccessTokenInContextUnchecked(ctx context.Context, accessTokenValue string) (success bool) {
	if accessTokenValue != "" {
		ctxvalues.SetAccessToken(ctx, accessTokenValue) // required for userinfo call to IDP
		return true
	}
	return false
}

func keyFuncForKey(rsaPublicKey *rsa.PublicKey) func(token *jwt.Token) (interface{}, error) {
	return func(token *jwt.Token) (interface{}, error) {
		return rsaPublicKey, nil
	}
}

type CustomClaims struct {
	Email         string   `json:"email"`
	EmailVerified bool     `json:"email_verified"`
	Groups        []string `json:"groups,omitempty"`
	Name          string   `json:"name"`
}

type AllClaims struct {
	jwt.RegisteredClaims
	CustomClaims
}

func checkIdToken_MustReturnOnError(ctx context.Context, idTokenValue string) (success bool, err error) {
	if idTokenValue != "" {
		tokenString := strings.TrimSpace(idTokenValue)

		errorMessage := ""
		for _, key := range config.OidcKeySet() {
			claims := AllClaims{}
			token, err := jwt.ParseWithClaims(tokenString, &claims, keyFuncForKey(key), jwt.WithValidMethods([]string{"RS256", "RS512"}))
			if err == nil && token.Valid {
				parsedClaims, ok := token.Claims.(*AllClaims)
				if ok {
					if config.OidcAllowedAudience() != "" {
						if len(parsedClaims.Audience) != 1 || parsedClaims.Audience[0] != config.OidcAllowedAudience() {
							return false, errors.New("token audience does not match")
						}
					}

					if config.OidcAllowedIssuer() != "" {
						if parsedClaims.Issuer != config.OidcAllowedIssuer() {
							return false, errors.New("token issuer does not match")
						}
					}

					ctxvalues.SetIdToken(ctx, idTokenValue)
					ctxvalues.SetEmail(ctx, parsedClaims.Email)
					ctxvalues.SetEmailVerified(ctx, parsedClaims.EmailVerified)
					ctxvalues.SetName(ctx, parsedClaims.Name)
					ctxvalues.SetSubject(ctx, parsedClaims.Subject)
					for _, group := range parsedClaims.Groups {
						ctxvalues.SetAuthorizedAsGroup(ctx, group)
					}

					return true, nil
				}
				errorMessage = "empty claims substructure"
			} else if err != nil {
				errorMessage = err.Error()
			} else {
				errorMessage = "token parsed but invalid"
			}
		}
		return false, errors.New(errorMessage)
	}
	return false, nil
}

func allow(actualMethod string, actualUrlPath string, allowedMethod string, allowedUrlPath string) bool {
	return actualMethod == allowedMethod && actualUrlPath == allowedUrlPath
}

func skipAuthCheckCompletely(method string, urlPath string) bool {
	// positive list for request URLs and Methods where the complete check can be skipped
	return allow(method, urlPath, http.MethodGet, "/v1/auth") || // login step 1
		allow(method, urlPath, http.MethodGet, "/v1/dropoff") || // login step 2
		allow(method, urlPath, http.MethodGet, "/v1/logout") || // logout
		allow(method, urlPath, http.MethodGet, "/") // healthcheck
}

// --- top level ---

func checkAllAuthentication_MustReturnOnError(ctx context.Context, method string, urlPath string, authHeaderValue string, idTokenCookieValue string, accessTokenCookieValue string) error {
	if skipAuthCheckCompletely(method, urlPath) {
		return nil
	}

	// try authorization header (gives only access token, so MUST use userinfo endpoint in controller to return useful info)
	success := recordAccessTokenInContextUnchecked(ctx, authHeaderValue)
	if success {
		return nil
	}

	// now try cookie pair
	success, err := checkIdToken_MustReturnOnError(ctx, idTokenCookieValue)
	if err != nil {
		return fmt.Errorf("invalid id token in cookie: %s", err.Error())
	}
	if success {
		success2 := recordAccessTokenInContextUnchecked(ctx, accessTokenCookieValue)
		if success2 {
			return nil
		}
	}

	// not supplying authorization is not a valid use case, there are endpoints that allow anonymous access
	return fmt.Errorf("failed to provide any authorization either via auth header or via cookies")
}

// --- middleware validating the values and adding to context values ---

func TokenValidator(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		authHeaderValue := fromAuthHeader(r)
		idTokenCookieValue := fromCookie(r, config.OidcIdTokenCookieName())
		accessTokenCookieValue := fromCookie(r, config.OidcAccessTokenCookieName())

		err := checkAllAuthentication_MustReturnOnError(ctx, r.Method, r.URL.Path, authHeaderValue, idTokenCookieValue, accessTokenCookieValue)
		if err != nil {
			UnauthenticatedError(ctx, w, r, "authorization failed to check out during local validation - please see logs for details", err.Error())
			return
		}

		// WARNING - at this point we might still have an unverified access token!

		next.ServeHTTP(w, r)
		return
	}
	return http.HandlerFunc(fn)
}

// --- accessors see ctxvalues ---

func UnauthenticatedError(ctx context.Context, w http.ResponseWriter, r *http.Request, details string, logMessage string) {
	aulogging.Logger.Ctx(ctx).Warn().Print(logMessage)
	ErrorHandler(ctx, w, r, "auth.unauthorized", http.StatusUnauthorized, url.Values{"details": []string{details}})
}

func ErrorHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, msg string, status int, details url.Values) {
	timestamp := time.Now().Format(time.RFC3339)
	response := errorapi.ErrorDto{Message: msg, Timestamp: timestamp, Details: details, RequestId: ctxvalues.RequestId(ctx)}
	w.Header().Set(headers.ContentType, media.ContentTypeApplicationJson)
	w.WriteHeader(status)
	WriteJson(ctx, w, response)
}

func WriteJson(ctx context.Context, w http.ResponseWriter, v interface{}) {
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(v)
	if err != nil {
		aulogging.Logger.Ctx(ctx).Warn().WithErr(err).Printf("error while encoding json response: %s", err.Error())
	}
}
