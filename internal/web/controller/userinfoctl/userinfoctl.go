package userinfoctl

import (
	"context"
	"encoding/json"
	"fmt"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	"github.com/eurofurence/reg-auth-service/internal/api/v1/errorapi"
	"github.com/eurofurence/reg-auth-service/internal/api/v1/userinfo"
	"github.com/eurofurence/reg-auth-service/internal/repository/config"
	"github.com/eurofurence/reg-auth-service/internal/repository/idp"
	"github.com/eurofurence/reg-auth-service/internal/web/util/ctxvalues"
	"github.com/eurofurence/reg-auth-service/internal/web/util/media"
	"github.com/go-chi/chi/v5"
	"github.com/go-http-utils/headers"
	"net/http"
	"net/url"
	"time"
)

var IDPClient idp.IdentityProviderClient

func Create(server chi.Router, idpClient idp.IdentityProviderClient) {
	if IDPClient == nil {
		IDPClient = idpClient
	}
	server.Get("/v1/userinfo", userinfoHandler)
}

func userinfoHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// ensure we have a id valid token at all
	if ctxvalues.BearerIdToken(ctx) == "" {
		unauthenticatedError(ctx, w, r, "you did not provide a valid token - see log for details", "no valid token in context - check logs above for validation errors")
		return
	}

	roles := []string{}
	for _, role := range config.RelevantRoles() {
		if ctxvalues.IsAuthorizedAsRole(ctx, role) {
			roles = append(roles, role)
		}
	}

	subject := ctxvalues.Subject(ctx)
	response := userinfo.UserInfoDto{
		Email: ctxvalues.Email(ctx),
		Roles: roles,
	}

	if config.OidcUserInfoURL() != "" {
		// ensure we have an access token at all
		if ctxvalues.BearerAccessToken(ctx) == "" {
			unauthenticatedError(ctx, w, r, "you did not provide a valid access token - see log for details", "no valid access token in context - check logs above for validation errors")
			return
		}

		idpUserinfo, status, err := IDPClient.UserInfo(ctx)
		if err != nil {
			idpDownstreamError(ctx, w, r, "identity provider could not be reached - see log for details", err.Error())
			return
		}
		if status == http.StatusUnauthorized || status == http.StatusForbidden {
			unauthenticatedError(ctx, w, r, "identity provider rejected your token - see log for details", fmt.Sprintf("idp returned rejection status %d", status))
			return
		}
		if status != http.StatusOK {
			idpDownstreamError(ctx, w, r, "identity provider returned error status - see log for details", fmt.Sprintf("idp returned error status %d", status))
			return
		}
		if idpUserinfo.Subject != subject {
			unauthenticatedError(ctx, w, r, "identity provider rejected your token - see log for details", fmt.Sprintf("idp returned different subject %s instead of %s", idpUserinfo.Subject, subject))
			return
		}

		for _, role := range roles {
			roleGivenInIdp := false
			for _, idpRole := range idpUserinfo.Global.Roles {
				if idpRole == role {
					roleGivenInIdp = true
				}
			}
			if !roleGivenInIdp {
				unauthenticatedError(ctx, w, r, "identity provider rejected your token - see log for details", fmt.Sprintf("role %s not assigned to subject %s according to idp", role, idpUserinfo.Subject))
				return
			}
		}

		if idpUserinfo.Global.Email != response.Email {
			unauthenticatedError(ctx, w, r, "identity provider rejected your token - email mismatch, please re-authenticate", fmt.Sprintf("subject %s email mismatch", idpUserinfo.Subject))
			return
		}
	}

	w.Header().Add(headers.ContentType, media.ContentTypeApplicationJson)
	w.WriteHeader(http.StatusOK)
	writeJson(ctx, w, response)
}

func idpDownstreamError(ctx context.Context, w http.ResponseWriter, r *http.Request, details string, logMessage string) {
	aulogging.Logger.Ctx(ctx).Warn().Print(logMessage)
	errorHandler(ctx, w, r, "auth.idp.error", http.StatusBadGateway, url.Values{"details": []string{details}})
}

func unauthenticatedError(ctx context.Context, w http.ResponseWriter, r *http.Request, details string, logMessage string) {
	aulogging.Logger.Ctx(ctx).Warn().Print(logMessage)
	errorHandler(ctx, w, r, "auth.unauthorized", http.StatusUnauthorized, url.Values{"details": []string{details}})
}

func errorHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, msg string, status int, details url.Values) {
	timestamp := time.Now().Format(time.RFC3339)
	response := errorapi.ErrorDto{Message: msg, Timestamp: timestamp, Details: details, RequestId: ctxvalues.RequestId(ctx)}
	w.Header().Set(headers.ContentType, media.ContentTypeApplicationJson)
	w.WriteHeader(status)
	writeJson(ctx, w, response)
}

func writeJson(ctx context.Context, w http.ResponseWriter, v interface{}) {
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(v)
	if err != nil {
		aulogging.Logger.Ctx(ctx).Warn().WithErr(err).Printf("error while encoding json response: %s", err.Error())
	}
}
