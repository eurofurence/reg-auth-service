package middleware

import (
	"crypto/rsa"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	"github.com/eurofurence/reg-auth-service/internal/repository/config"
	"github.com/eurofurence/reg-auth-service/internal/web/util/ctxvalues"
	"github.com/go-http-utils/headers"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strings"
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

	return "Bearer " + authCookie.Value
}

func fromAuthHeader(r *http.Request) string {
	return r.Header.Get(headers.Authorization)
}

func fromAuthHeaderOrCookie(r *http.Request, cookieName string) string {
	h := fromAuthHeader(r)
	if h == "" {
		return fromCookie(r, cookieName)
	} else {
		return h
	}
}

// --- middleware validating the values and adding to context values ---

func keyFuncForKey(rsaPublicKey *rsa.PublicKey) func(token *jwt.Token) (interface{}, error) {
	return func(token *jwt.Token) (interface{}, error) {
		return rsaPublicKey, nil
	}
}

type GlobalClaims struct {
	Name  string   `json:"name"`
	EMail string   `json:"email"`
	Roles []string `json:"roles"`
}

type CustomClaims struct {
	Global GlobalClaims `json:"global"`
}

type AllClaims struct {
	jwt.RegisteredClaims
	CustomClaims
}

func TokenValidator(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// try bearer token from either cookie or Authorization header
		// (in this one service, do NOT fail even if authentication is invalid - we have endpoints that are there to remedy exactly that situation here)
		bearerTokenValue := fromAuthHeaderOrCookie(r, config.OidcIdTokenCookieName())
		if bearerTokenValue != "" {
			const bearerPrefix = "Bearer "
			errorMessage := ""
			if !strings.HasPrefix(bearerTokenValue, bearerPrefix) {
				errorMessage = "value of Authorization header did not start with 'Bearer '"
			} else {
				tokenString := strings.TrimSpace(strings.TrimPrefix(bearerTokenValue, bearerPrefix))

				for _, key := range config.OidcKeySet() {
					claims := AllClaims{}
					token, err := jwt.ParseWithClaims(tokenString, &claims, keyFuncForKey(key), jwt.WithValidMethods([]string{"RS256", "RS512"}))
					if err == nil && token.Valid {
						parsedClaims, ok := token.Claims.(*AllClaims)
						if ok {
							ctxvalues.SetBearerIdToken(ctx, bearerTokenValue)
							ctxvalues.SetEmail(ctx, parsedClaims.Global.EMail)
							ctxvalues.SetName(ctx, parsedClaims.Global.Name)
							ctxvalues.SetSubject(ctx, parsedClaims.Subject)
							for _, role := range parsedClaims.Global.Roles {
								ctxvalues.SetAuthorizedAsRole(ctx, role)
							}

							if config.OidcAccessTokenCookieName() != "" {
								authTokenValue := fromCookie(r, config.OidcAccessTokenCookieName())
								if authTokenValue != "" {
									ctxvalues.SetBearerAccessToken(ctx, authTokenValue)
								} else {
									aulogging.Logger.Ctx(ctx).Warn().Printf("got id token, but no auth token for subject %s - continuing, but userinfo will fail", parsedClaims.Subject)
								}
							}

							next.ServeHTTP(w, r)
							return
						} else {
							errorMessage = "empty claims substructure"
						}
					} else if err != nil {
						errorMessage = "token parse error: " + err.Error()
					} else {
						errorMessage = "token parsed but invalid"
					}
				}
			}

			if errorMessage != "" {
				// log a warning, but still continue
				aulogging.Logger.Ctx(ctx).Warn().Print(errorMessage)
			}
		}

		// not supplying either is a valid use case, there are endpoints that allow anonymous access
		next.ServeHTTP(w, r)
		return
	}
	return http.HandlerFunc(fn)
}

// --- accessors see ctxvalues ---
