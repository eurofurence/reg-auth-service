package acceptance

import (
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
 *  * reg_app_name  - the name of the application that the user wants to be authenticated for
 *
 * Optional parameters are:
 *  * redirect_url  - where to redirect the user after a successfull authentication flow.
 *                    This URL must match the pattern of allowed URLs in the config file.
 *
 * All additional query parameters are appended to the app's redirect_url after a successfull
 * authentication.
 */

func TestLogin(t *testing.T) {
	docs.Given("given the standard test configuration")
	tstSetup(tstDefaultConfigFile)
	defer tstShutdown()

	docs.When("when they start an auth flow with valid reg_app_name and valid redirect_url")
	test_url := "/v1/auth?app_name=example-service&"
	test_url = test_url + "&dropoff_url=" + url.QueryEscape("https://example.com/app/?foo=abc")
	response := tstPerformGet(test_url)

	docs.Then("then the user agent is redirected to the OpenID Connect auth URL")
	// The user agent must be redirected using an HTTP 302 MOVED response.
	require.Equal(t, http.StatusFound, response.StatusCode, "unexpected http response status")
	// The 'Location:' header must match the application config's authorization_endpoint
	location_url := response.Header.Get("Location")
	require.NotNil(t, location_url, "missing or invalid Location header")
	loc, err := url.Parse(location_url)
	require.Nil(t, err, "Location header could not be parsed as a URL")
	require.Equal(t, "https", loc.Scheme, "unexpected Location scheme")
	require.Equal(t, "auth.example.com", loc.Host, "unexpected Location host")
	require.Equal(t, "/auth", loc.Path, "unexpected Location path")
	// Parse query parameters
	values := loc.Query()
	// There must be a 'scope' parameter and it must match the application config's scope(s)
	require.Equal(t, "example", values.Get("scope"), "unexpected scope")
	// The 'state' (nonce) parameter must be present
	require.NotNil(t, values.Get("state"), "missing state parameter")
	// The 'response_type' parameter must be present and contain the string "code"
	require.Equal(t, "code", values.Get("response_type"), "unexpected scope")
	// The 'redirect_url' parameter must contain the URL of the /redirect endpoint of this service
	// Note: this is *NOT* the redirect_url that we might receive as an optional input parameter.
	require.Equal(t, "http://localhost:8081/v1/dropoff", values.Get("redirect_url"), "unexpected redirect_url")
}

