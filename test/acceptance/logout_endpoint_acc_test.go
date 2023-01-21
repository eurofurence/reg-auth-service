package acceptance

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/eurofurence/reg-auth-service/docs"
	"github.com/stretchr/testify/require"
)

// ----------------------------------------
// acceptance tests for the logout endpoint
// ----------------------------------------

/* The /logout endpoint deletes the cookie and redirects to the app's default dropoff url.
 *
 * Required parameters are:
 *  * app_name  - the name of the application that the user wants to be authenticated for
 */

func TestLogout_Success(t *testing.T) {
	docs.Given("given the standard test configuration")
	tstSetup(tstDefaultConfigFile)
	defer tstShutdown()

	docs.When("when they call the logout endpoint with a valid app_name")
	testUrl := "/v1/logout?app_name=example-service"
	response := tstPerformGetNoRedirect(testUrl)

	docs.Then("then the user agent is redirected to the default drop off URL of the app and the cookie deleted")
	require.Equal(t, http.StatusFound, response.StatusCode, "unexpected http response status, must be HTTP 302 MOVED")
	location_url := response.Header.Get("Location")
	require.NotNil(t, location_url, "missing or invalid Location header")
	loc, err := url.Parse(location_url)
	require.Nil(t, err, "Location header could not be parsed as a URL")
	require.Equal(t, "https", loc.Scheme, "unexpected Location scheme, must match the configured default drop off URL")
	require.Equal(t, "example.com", loc.Host, "unexpected Location host, must match the configured default drop off URL")
	require.Equal(t, "/app/", loc.Path, "unexpected Location path, must match the configured default drop off URL")

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
	require.NotNil(t, ac, "Access token cookie must be present")
	require.Equal(t, "", id.Value)
	require.Equal(t, "", ac.Value)
	require.Equal(t, "example.com", ac.Domain)
}

func TestLogout_Failure_AppNameMissing(t *testing.T) {
	docs.Given("given the standard test configuration")
	tstSetup(tstDefaultConfigFile)
	defer tstShutdown()

	docs.When("when they call the logout endpoint, but do not specify an app_name")
	testUrl := "/v1/logout"
	response := tstPerformGetNoRedirect(testUrl)

	docs.Then("then the correct error is displayed")
	require.Equal(t, http.StatusBadRequest, response.StatusCode, "unexpected http response status, must be HTTP 400")
	responseBody := tstResponseBodyString(&response)
	require.Contains(t, responseBody, "<b>error:</b> invalid parameters")
}

func TestLogout_Failure_UnknownAppName(t *testing.T) {
	docs.Given("given the standard test configuration")
	tstSetup(tstDefaultConfigFile)
	defer tstShutdown()

	docs.When("when they call the logout endpoint, but specify an unknown app_name")
	testUrl := "/v1/logout?app_name=unknown-service"
	response := tstPerformGetNoRedirect(testUrl)

	docs.Then("then the correct error is displayed")
	require.Equal(t, http.StatusNotFound, response.StatusCode, "unexpected http response status, must be HTTP 404")
	responseBody := tstResponseBodyString(&response)
	require.Contains(t, responseBody, "<b>error:</b> invalid parameters")
}
