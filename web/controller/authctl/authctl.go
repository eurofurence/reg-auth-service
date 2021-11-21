package authctl

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/eurofurence/reg-auth-service/internal/entity"
	"github.com/eurofurence/reg-auth-service/internal/repository/config"
	"github.com/eurofurence/reg-auth-service/internal/repository/database"
	"github.com/eurofurence/reg-auth-service/internal/repository/logging"
	"github.com/go-chi/chi"
)

const responseType = "code"
const codeChallengeMethod = "S256"

func Create(server chi.Router) {
	server.Get("/v1/auth", authHandler)
}

/* Handle /auth requests.
 *
 * Required parameters are:
 *  * reg_app_name  - the name of the application that the user wants to be authenticated for
 *
 * Optional parameters are:
 *  * redirect_url  - where to redirect the user after a successfull authentication flow.
 *                    This URL must match the pattern of allowed URLs in the config file.
 *
 * All additional query parameters are appended to the app's redirect_url after a successfull
 * authentication.
 */
func authHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	query := r.URL.Query()
	regAppName := query.Get("reg_app_name")
	if regAppName == "" {
		authErrorHandler(ctx, w, http.StatusBadRequest, "reg_app_name parameter is missing")
		return
	}
	applicationConfig, err := config.GetApplicationConfig(regAppName)
	if err != nil {
		authErrorHandler(ctx, w, http.StatusNotFound, "reg_app_name is unknown")
		return
	}

	dropOffUrl := query.Get("redirect_url")
	if dropOffUrl == "" {
		dropOffUrl = applicationConfig.DefaultRedirectUrl
	} else {
		if !validateDropOffURL(ctx, w, applicationConfig.RedirectUrlPattern, dropOffUrl) {
			authErrorHandler(ctx, w, http.StatusForbidden, "the specified redirect_url is not allowed")
			return
		}
	}

	state, err := generateState()
	if err != nil {
		authErrorHandler(ctx, w, http.StatusInternalServerError, "state could not be generated")
		return
	}
	codeVerifier, err := generateCodeVerifier()
	if err != nil {
		authErrorHandler(ctx, w, http.StatusInternalServerError, "verifier could not be generated")
		return
	}
	codeChallenge := generateCodeChallenge(codeVerifier)

	err = storeFlowState(ctx, state, codeVerifier, dropOffUrl)
	if err != nil {
		authErrorHandler(ctx, w, http.StatusInternalServerError, "could not store flow state")
		return
	}

	err = redirectToOpenIDProvider(ctx, w, applicationConfig, state, codeChallenge)
	if err != nil {
		authErrorHandler(ctx, w, http.StatusInternalServerError, err.Error())
	}
}

func authErrorHandler(ctx context.Context, w http.ResponseWriter, status int, msg string) {
	logging.Ctx(ctx).Error(msg)
	// TODO: here we should display some information to the user
	w.WriteHeader(status)
}

func validateDropOffURL(ctx context.Context, w http.ResponseWriter, exp string, dropOffUrl string) bool {
	match, err := regexp.MatchString(exp, dropOffUrl)
	if err != nil {
		logging.Ctx(ctx).Error("could not match regular expression: " + err.Error())
		return false
	}
	return match
}

/* according to RFC 6749, "state" is defined as one or more characters within
 * the range of US ASCII  %20 - %7E (printable ASCII characters). See here:
 *
 *    https://datatracker.ietf.org/doc/html/rfc6749#appendix-A.5
 */
func generateState() (string, error) {
	const length = 40
	state := make([]byte, length)
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(0x7e-0x20)))
		if err != nil {
			return "", err
		}
		state[i] = byte(num.Int64())
	}
	return string(state), nil
}

/* according to RFC 7636, the "code verifier" is defined as between 43 and 12
 * characters within the range of US ASCII a-zA-Z0-9 and any of "-", ".", "_" or "~".
 * See here:
 *
 *    https://datatracker.ietf.org/doc/html/rfc7636#section-4.1
 */
func generateCodeVerifier() (string, error) {
	const letters string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~"
	const length = 128
	verifier := make([]byte, length)
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		verifier[i] = letters[num.Int64()]
	}
	return string(verifier), nil
}

/* according to RFC 7636, the code challende is derived from the varifier
 * this way:   code_challenge = base64_encode(sha256(verifier))
 * See here:
 *
 *    https://datatracker.ietf.org/doc/html/rfc7636#section-4.2
 */
func generateCodeChallenge(verifier string) string {
	h := sha256.New()
	h.Write([]byte(verifier))
	hash := h.Sum(nil)
	return base64.StdEncoding.EncodeToString(hash)
}

func storeFlowState(ctx context.Context, state string, codeVerifier string, dropOffUrl string) error {
	return database.GetRepository().AddAuthRequest(ctx, &entity.AuthRequest{
		State:            state,
		PkceCodeVerifier: codeVerifier,
		DropOffUrl:       dropOffUrl,
		ExpiresAt:        time.Now().Add(config.AuthRequestTimeout()),
	})
}

func redirectToOpenIDProvider(ctx context.Context, w http.ResponseWriter, applicationConfig config.ApplicationConfig, state string, codeChallenge string) error {
	u, err := url.Parse(config.AuthorizationEndpoint())
	if err != nil {
		return fmt.Errorf("could not parse auth endpoint url")
	}
	q := u.Query()
	q.Set("response_type", responseType)
	q.Set("client_id", applicationConfig.ClientId)
	q.Set("scope", applicationConfig.Scope)
	q.Set("state", state)
	q.Set("code_challenge", codeChallenge)
	q.Set("code_challenge_method", codeChallengeMethod)
	q.Set("redirect_url", config.ServerAddr()+"/send_off")
	u.RawQuery = q.Encode()
	w.Header().Set("Location", u.String())
	w.WriteHeader(http.StatusFound)
	return nil
}
