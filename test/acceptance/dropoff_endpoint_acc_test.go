package acceptance

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/eurofurence/reg-auth-service/docs"
	"github.com/stretchr/testify/require"
)

// -----------------------------------------
// acceptance tests for the dropoff endpoint
// -----------------------------------------

/* The /dropoff endpoint is part of the OpenID Connect authorization code flow. Once the OpenID
 * Connect provider agrees to provide an access token it redirects the user agent to this
 * endpoint. Here, the reg-auth-service obtains the access token from the OIDC provider,
 * stores it in a cookie, and then redirects the user agent once more to the URL the
 * user agent initially inteded to visit. (the redirect_url)
 * 
 * Required parameters are:
 *  * state              - random-string identifier of this flow
 *  * authorization_code - temporary credential to obtain the access token from the OIDC
 *
 */

func TestDropOff(t *testing.T) {
	docs.Given("given the standard test configuration")
	tstSetup(tstDefaultConfigFile)
	defer tstShutdown()

	docs.When("when they call the endpoint with valid state and valid authorization_code")
	test_url := "/v1/dropoff?state="+tstAuthRequest.State+"&"
	test_url = test_url + "&authorization_code="+tstAuthorizationCode

	response := tstPerformGet(test_url)
	docs.Then("then the user agent is redirected to the drop off URL")
	// The user agent must be redirected using an HTTP 302 MOVED response.
	require.Equal(t, http.StatusFound, response.StatusCode, "unexpected http response status")
	// The 'Location:' header must match the drop off URL
	location_url := response.Header.Get("Location")
	require.NotNil(t, location_url, "missing or invalid Location header")
	loc, err := url.Parse(location_url)
	require.Nil(t, err, "Location header could not be parsed as a URL")
	require.Equal(t, "https", loc.Scheme, "unexpected Location scheme")
	require.Equal(t, "example.com", loc.Host, "unexpected Location host")
	require.Equal(t, "/drop_off_url", loc.Path, "unexpected Location path")
	// The query parameter from the dropOffUrl should still be there
	values := loc.Query()
	require.Equal(t, "5", values.Get("dingbaz"), "unexpected scope")
	// The access code must be in the cookie
	cookies := response.Cookies()
	var ac *http.Cookie = nil
	for _, cookie := range cookies {
		if cookie.Name != "AccessCode" {
			continue
		}
		ac = cookie
	}
	require.NotNil(t, ac, "AccessCode cookie must be present")
	require.Equal(t, "dummy_mock_value", ac.Value)
	require.Equal(t, "example.com", ac.Domain)
}

