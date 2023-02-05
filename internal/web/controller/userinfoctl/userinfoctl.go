package userinfoctl

import (
	"context"
	"encoding/json"
	"errors"
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
	server.Get("/v1/frontend-userinfo", frontendUserinfoHandler)
}

// localUserinfoHelper is the common part of both handlers.
func localUserinfoHelper(ctx context.Context, w http.ResponseWriter, r *http.Request) (userinfo.UserInfoDto, error) {
	// ensure we have a valid id token
	if ctxvalues.IdToken(ctx) == "" {
		unauthenticatedError(ctx, w, r, "you did not provide a valid token - see log for details", "no valid token in context - check logs above for validation errors")
		return userinfo.UserInfoDto{}, errors.New("no id token")
	}

	// ensure we have an access token at all
	if ctxvalues.AccessToken(ctx) == "" {
		unauthenticatedError(ctx, w, r, "you did not provide a valid access token - see log for details", "no valid access token in context - check logs above for validation errors")
		return userinfo.UserInfoDto{}, errors.New("no access token")
	}

	response := userinfo.UserInfoDto{
		Email:         ctxvalues.Email(ctx),
		EmailVerified: ctxvalues.EmailVerified(ctx),
		Groups:        []string{},
		Name:          ctxvalues.Name(ctx),
		Subject:       ctxvalues.Subject(ctx),
	}

	for _, group := range config.RelevantGroups() {
		if ctxvalues.IsAuthorizedAsGroup(ctx, group) {
			response.Groups = append(response.Groups, group)
		}
	}

	return response, nil
}

func userinfoHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	response, err := localUserinfoHelper(ctx, w, r)
	if err != nil {
		// unauthenticatedError sent already
		return
	}

	if config.OidcUserInfoURL() == "" {
		// we must accept the token info, or local testing won't work
		aulogging.Logger.Ctx(ctx).Warn().Print("skipping token validation with IDP and taking info from token - this is not safe for production!")
		w.Header().Add(headers.ContentType, media.ContentTypeApplicationJson)
		w.WriteHeader(http.StatusOK)
		writeJson(ctx, w, response)
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
	if idpUserinfo.Subject != response.Subject {
		unauthenticatedError(ctx, w, r, "identity provider rejected your token - see log for details", fmt.Sprintf("idp returned different subject %s instead of %s", idpUserinfo.Subject, response.Subject))
		return
	}

	for _, group := range response.Groups {
		groupGivenInIdp := false
		for _, idpGroup := range idpUserinfo.Groups {
			if idpGroup == group {
				groupGivenInIdp = true
			}
		}
		if !groupGivenInIdp {
			unauthenticatedError(ctx, w, r, "identity provider rejected your token - see log for details", fmt.Sprintf("group %s not assigned to subject %s according to idp", group, idpUserinfo.Subject))
			return
		}
	}

	if idpUserinfo.Email != response.Email {
		unauthenticatedError(ctx, w, r, "identity provider rejected your token - email mismatch, please re-authenticate", fmt.Sprintf("subject %s email mismatch", idpUserinfo.Subject))
		return
	}

	w.Header().Add(headers.ContentType, media.ContentTypeApplicationJson)
	w.WriteHeader(http.StatusOK)
	writeJson(ctx, w, response)
}

func frontendUserinfoHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	response, err := localUserinfoHelper(ctx, w, r)
	if err != nil {
		// unauthenticatedError sent already
		return
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
