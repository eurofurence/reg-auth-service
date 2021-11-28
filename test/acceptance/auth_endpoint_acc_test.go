package acceptance

import (
	"context"
	"github.com/eurofurence/reg-auth-service/internal/repository/database"
	"net/http"
	"net/url"
	"testing"

	"github.com/eurofurence/reg-auth-service/docs"
	"github.com/stretchr/testify/require"
)

// --------------------------------------
// acceptance tests for the auth endpoint
// --------------------------------------

/* The /auth endpoint begins an OpenID Connect authentication code flow. Some (potentially unknown)
 * website or application has realized that an unauthenticated user tried to access their service,
 * and now the user is being directed to the /auth endpoint so they can log in.
 *
 * The config file contains a list of valid application configurations that this service can perform
 * authentication flows for. For each application configuration, a client_id, a client_secret, a
 * pattern for valid redirect_urls, a list of scopes and so on are configured. (see example config file)
 *
 * Required parameters are:
 *  * app_name  - the name of the application that the user wants to be authenticated for
 *
 * Optional parameters are:
 *  * dropoff_url   - where to redirect the user after a successfull authentication flow.
 *                    This URL must match the pattern of allowed URLs in the config file.
 *
 * All additional query parameters are appended to the app's dropoff_url after a successfull
 * authentication.
 */

func TestAuth_Success_DropoffUrlSpecified(t *testing.T) {
	docs.Given("given the standard test configuration")
	tstSetup(tstDefaultConfigFile)
	defer tstShutdown()

	docs.When("when they start an auth flow with valid app_name and valid dropoff_url")
	testUrl := "/v1/auth?app_name=example-service"
	testUrl = testUrl + "&dropoff_url=" + url.QueryEscape("https://example.com/app/?foo=abc")
	response := tstPerformGet(testUrl)

	docs.Then("then the user agent is redirected to the OpenID Connect auth URL with the correct parameters")
	require.Equal(t, http.StatusFound, response.StatusCode, "unexpected http response status, must be HTTP 302 MOVED")

	locationUrl := response.Header.Get("Location")
	require.NotNil(t, locationUrl, "missing or invalid Location header, must match the application config's authorization_endpoint")
	loc, err := url.Parse(locationUrl)
	require.Nil(t, err, "Location header could not be parsed as a URL")
	require.Equal(t, "https", loc.Scheme, "unexpected Location scheme, must match the application config's authorization_endpoint")
	require.Equal(t, "auth.example.com", loc.Host, "unexpected Location host, must match the application config's authorization_endpoint")
	require.Equal(t, "/auth", loc.Path, "unexpected Location path, must match the application config's authorization_endpoint")

	values := loc.Query()
	require.Equal(t, "IAmNotSoSecret.", values.Get("client_id"), "unexpected client_id, must match the application config")
	require.NotEmpty(t, values.Get("code_challenge"))
	require.Equal(t, "S256", values.Get("code_challenge_method"), "unexpected code_challenge_method, must be 'S256'")
	// Note: this is *NOT* the dropoff_url that we might receive as an optional input parameter.
	require.Equal(t, "http://localhost:8081/v1/dropoff", values.Get("redirect_url"), "unexpected redirect_url parameter, must be the URL of the /dropoff endpoint of this service")
	require.Equal(t, "code", values.Get("response_type"), "unexpected response_type parameter, must be 'code'")
	require.Equal(t, "example", values.Get("scope"), "unexpected scope parameter, must match the application config's scope(s)")
	state := values.Get("state")
	require.NotEmpty(t, state, "missing state (nonce) parameter")

	docs.Then("and the provided dropoff url is stored internally for this state")
	internalStateData, err := database.GetRepository().GetAuthRequestByState(context.TODO(), state)
	require.Nil(t, err)
	require.Equal(t, "https://example.com/app/?foo=abc", internalStateData.DropOffUrl)
}

func TestAuth_Success_DefaultDropoffUrl(t *testing.T) {
	docs.Given("given the standard test configuration")
	tstSetup(tstDefaultConfigFile)
	defer tstShutdown()

	docs.When("when they start an auth flow with valid app_name and do not specify a dropoff_url")
	testUrl := "/v1/auth?app_name=example-service"
	response := tstPerformGet(testUrl)

	docs.Then("then the user agent is redirected to the OpenID Connect auth URL with the correct parameters")
	require.Equal(t, http.StatusFound, response.StatusCode, "unexpected http response status, must be HTTP 302 MOVED")

	locationUrl := response.Header.Get("Location")
	require.NotNil(t, locationUrl, "missing or invalid Location header, must match the application config's authorization_endpoint")
	loc, err := url.Parse(locationUrl)
	require.Nil(t, err, "Location header could not be parsed as a URL")
	require.Equal(t, "https", loc.Scheme, "unexpected Location scheme, must match the application config's authorization_endpoint")
	require.Equal(t, "auth.example.com", loc.Host, "unexpected Location host, must match the application config's authorization_endpoint")
	require.Equal(t, "/auth", loc.Path, "unexpected Location path, must match the application config's authorization_endpoint")

	values := loc.Query()
	require.Equal(t, "IAmNotSoSecret.", values.Get("client_id"), "unexpected client_id, must match the application config")
	require.NotEmpty(t, values.Get("code_challenge"))
	require.Equal(t, "S256", values.Get("code_challenge_method"), "unexpected code_challenge_method, must be 'S256'")
	// Note: this is *NOT* the dropoff_url that we might receive as an optional input parameter.
	require.Equal(t, "http://localhost:8081/v1/dropoff", values.Get("redirect_url"), "unexpected redirect_url parameter, must be the URL of the /dropoff endpoint of this service")
	require.Equal(t, "code", values.Get("response_type"), "unexpected response_type parameter, must be 'code'")
	require.Equal(t, "example", values.Get("scope"), "unexpected scope parameter, must match the application config's scope(s)")
	state := values.Get("state")
	require.NotEmpty(t, state, "missing state (nonce) parameter")

	docs.Then("and the default dropoff url is stored internally for this state")
	internalStateData, err := database.GetRepository().GetAuthRequestByState(context.TODO(), state)
	require.Nil(t, err)
	require.Equal(t, "https://example.com/app/", internalStateData.DropOffUrl)
}

func TestAuth_Failure_AppNameMissing(t *testing.T) {
	docs.Given("given the standard test configuration")
	tstSetup(tstDefaultConfigFile)
	defer tstShutdown()

	docs.When("when they start an auth flow, but do not specify an app_name")
	testUrl := "/v1/auth"
	response := tstPerformGet(testUrl)

	docs.Then("then the correct error is displayed")
	require.Equal(t, http.StatusBadRequest, response.StatusCode, "unexpected http response status, must be HTTP 400")
	responseBody := tstResponseBodyString(&response)
	require.Contains(t, responseBody, "<b>error:</b> invalid parameters")
}

func TestAuth_Failure_UnknownAppName(t *testing.T) {
	docs.Given("given the standard test configuration")
	tstSetup(tstDefaultConfigFile)
	defer tstShutdown()

	docs.When("when they start an auth flow, but specify an unknown app_name")
	testUrl := "/v1/auth?app_name=unknown-service"
	response := tstPerformGet(testUrl)

	docs.Then("then the correct error is displayed")
	require.Equal(t, http.StatusNotFound, response.StatusCode, "unexpected http response status, must be HTTP 404")
	responseBody := tstResponseBodyString(&response)
	require.Contains(t, responseBody, "<b>error:</b> invalid parameters")
}

func TestAuth_Failure_InvalidDropoffUrl(t *testing.T) {
	docs.Given("given the standard test configuration")
	tstSetup(tstDefaultConfigFile)
	defer tstShutdown()

	docs.When("when they start an auth flow, but specify a dropoff url that does not match the configured pattern")
	testUrl := "/v1/auth?app_name=example-service"
	testUrl = testUrl + "&dropoff_url=" + url.QueryEscape("https://example.com/nomatch")
	response := tstPerformGet(testUrl)

	docs.Then("then the correct error is displayed")
	require.Equal(t, http.StatusForbidden, response.StatusCode, "unexpected http response status, must be HTTP 403")
	responseBody := tstResponseBodyString(&response)
	require.Contains(t, responseBody, "<b>error:</b> invalid parameters")
}
