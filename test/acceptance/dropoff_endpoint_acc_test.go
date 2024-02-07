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
 * user agent initially intended to visit. (the dropoff url)
 *
 * Required parameters are:
 *  * state - random-string identifier of this flow
 *  * code  - temporary credential to obtain the access token from the OIDC provider
 *
 * Error parameters are:
 *  * state - random-string identifier of this flow
 *  * error - the error code
 *  * error_description - human-readable description text
 */

func TestDropoff_Success(t *testing.T) {
	docs.Given("given the standard test configuration")
	tstSetup(tstDefaultConfigFile)
	defer tstShutdown()

	docs.When("when they call the dropoff endpoint with valid state and valid authorization_code")
	test_url := "/v1/dropoff?state=" + tstAuthRequest.State
	test_url = test_url + "&code=" + tstAuthorizationCode
	response := tstPerformGetNoRedirect(test_url)

	docs.Then("then the user agent is redirected to the drop off URL")
	require.Equal(t, http.StatusFound, response.StatusCode, "unexpected http response status, must be HTTP 302 MOVED")
	location_url := response.Header.Get("Location")
	require.NotNil(t, location_url, "missing or invalid Location header")
	loc, err := url.Parse(location_url)
	require.Nil(t, err, "Location header could not be parsed as a URL")
	require.Equal(t, "https", loc.Scheme, "unexpected Location scheme, must match the drop off URL")
	require.Equal(t, "example.com", loc.Host, "unexpected Location host, must match the drop off URL")
	require.Equal(t, "/drop_off_url", loc.Path, "unexpected Location path, must match the drop off URL")

	values := loc.Query()
	require.Equal(t, "5", values.Get("dingbaz"), "query parameter from the dropOffUrl should still be there")

	cookies := response.Cookies()
	var ac *http.Cookie = nil
	var id *http.Cookie = nil
	for _, cookie := range cookies {
		if cookie.Name == "JWT" {
			id = cookie
		}
		if cookie.Name == "AUTH" {
			ac = cookie
		}
	}
	require.NotNil(t, id, "Id token cookie must be present")
	require.NotNil(t, ac, "Auth token cookie must be present")
	require.Equal(t, "dummy_mock_value", id.Value)
	require.Equal(t, "example.com", id.Domain)
	require.Equal(t, "access_mock_value", ac.Value)
	require.Equal(t, "example.com", ac.Domain)
}

func TestDropoff_Failure_IDPError(t *testing.T) {
	docs.Given("given the standard test configuration")
	tstSetup(tstDefaultConfigFile)
	defer tstShutdown()

	docs.When("when they call the dropoff endpoint with an error")
	test_url := "/v1/dropoff?state=" + tstAuthRequest.State + "&error=request_unauthorized&error_description=The+request+could+not+be+authorized"
	response := tstPerformGetNoRedirect(test_url)

	docs.Then("then the correct error is displayed")
	require.Equal(t, http.StatusBadRequest, response.StatusCode, "unexpected http response status, must be HTTP 400")
	responseBody := tstResponseBodyString(&response)
	require.Contains(t, responseBody, "<b>error:</b> request_unauthorized: The request could not be authorized")
	require.Contains(t, responseBody, `You can also <a href="https://example.com/app/">go back to try again</a>.`)
}

func TestDropoff_Failure_StateMissing(t *testing.T) {
	docs.Given("given the standard test configuration")
	tstSetup(tstDefaultConfigFile)
	defer tstShutdown()

	docs.When("when they call the dropoff endpoint without a state parameter")
	test_url := "/v1/dropoff?code=" + tstAuthorizationCode
	response := tstPerformGetNoRedirect(test_url)

	docs.Then("then the correct error is displayed")
	require.Equal(t, http.StatusBadRequest, response.StatusCode, "unexpected http response status, must be HTTP 400")
	responseBody := tstResponseBodyString(&response)
	require.Contains(t, responseBody, "<b>error:</b> invalid parameters")
}

func TestDropoff_Failure_CodeMissing(t *testing.T) {
	docs.Given("given the standard test configuration")
	tstSetup(tstDefaultConfigFile)
	defer tstShutdown()

	docs.When("when they call the dropoff endpoint without a code parameter")
	test_url := "/v1/dropoff?state=" + tstAuthRequest.State
	response := tstPerformGetNoRedirect(test_url)

	docs.Then("then the correct error is displayed")
	require.Equal(t, http.StatusBadRequest, response.StatusCode, "unexpected http response status, must be HTTP 400")
	responseBody := tstResponseBodyString(&response)
	require.Contains(t, responseBody, "<b>error:</b> invalid parameters")
}

func TestDropoff_Failure_StateNotFound(t *testing.T) {
	docs.Given("given the standard test configuration")
	tstSetup(tstDefaultConfigFile)
	defer tstShutdown()

	docs.When("when they call the dropoff endpoint with a state parameter that is not in the internal storage (possibly expired)")
	test_url := "/v1/dropoff?state=notthere"
	test_url = test_url + "&code=" + tstAuthorizationCode
	response := tstPerformGetNoRedirect(test_url)

	docs.Then("then the correct error is displayed")
	require.Equal(t, http.StatusNotFound, response.StatusCode, "unexpected http response status, must be HTTP 404")
	responseBody := tstResponseBodyString(&response)
	require.Contains(t, responseBody, "<b>error:</b> auth request not found or timed out")
}
