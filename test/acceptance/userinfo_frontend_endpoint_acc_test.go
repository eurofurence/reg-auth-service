package acceptance

import (
	"github.com/eurofurence/reg-auth-service/docs"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

// ---------------------------------------------------
// acceptance tests for the frontend-userinfo endpoint
// ---------------------------------------------------

func TestFrontendUserinfo_Success(t *testing.T) {
	docs.Given("given the standard test configuration")
	tstSetup(tstDefaultConfigFile)
	defer tstShutdown()

	docs.When("when a logged in user calls the frontend-userinfo endpoint with a valid token")
	response := tstPerformGetWithCookies("/v1/frontend-userinfo", valid_JWT_id_is_not_staff_sub101, "access_mock_value 101")

	docs.Then("then the request is successful and the response is as expected")
	tstRequireUserinfoResponse(t, response, expected_response_by_token[valid_JWT_id_is_not_staff_sub101])

	docs.Then("and no calls to the IDP have been made")
	require.Empty(t, idpMock.recording)
}

func TestFrontendUserinfo_Success_Admin(t *testing.T) {
	docs.Given("given the standard test configuration")
	tstSetup(tstDefaultConfigFile)
	defer tstShutdown()

	docs.When("when a logged in admin calls the frontend-userinfo endpoint with a valid token")
	response := tstPerformGetWithCookies("/v1/frontend-userinfo", valid_JWT_id_is_staff_admin_sub1234567890, "access_mock_value")

	docs.Then("then the request is successful and the response is as expected")
	tstRequireUserinfoResponse(t, response, expected_response_by_token[valid_JWT_id_is_staff_admin_sub1234567890])

	docs.Then("and no calls to the IDP have been made")
	require.Empty(t, idpMock.recording)
}

func TestFrontendUserinfo_WrongSubjectForGroup(t *testing.T) {
	docs.Given("given the standard test configuration")
	tstSetup(tstDefaultConfigFile)
	defer tstShutdown()

	docs.When("when a logged in user calls the frontend-userinfo endpoint with a valid token with an extra relevant group (admin) for which the user's subject isn't in the allowlist")
	response := tstPerformGetWithCookies("/v1/frontend-userinfo", valid_JWT_id_is_staff_false_admin_sub444, "access_mock_value 444")

	docs.Then("then the request is successful and the list of groups does NOT include the extra group")
	expected := expected_response_by_token[valid_JWT_id_is_staff_false_admin_sub444]
	expected.Name = "John Admin" // trusting ID token here
	tstRequireUserinfoResponse(t, response, expected)

	docs.Then("and no calls to the IDP have been made")
	require.Empty(t, idpMock.recording)
}

func TestFrontendUserinfo_NoAccessTokenCookie(t *testing.T) {
	docs.Given("given the standard test configuration")
	tstSetup(tstDefaultConfigFile)
	defer tstShutdown()

	docs.When("when a logged user calls the frontend-userinfo endpoint with a valid token but supplies no access token")
	response := tstPerformGetWithCookies("/v1/frontend-userinfo", valid_JWT_id_is_not_staff_sub101, "")

	docs.Then("then the request fails with the appropriate error")
	tstRequireErrorResponse(t, response, http.StatusUnauthorized, "auth.unauthorized", "authorization failed to check out during local validation - please see logs for details")

	docs.Then("and no calls to the IDP have been made")
	require.Empty(t, idpMock.recording)
}

func TestFrontendUserinfo_InvalidAccessTokenCookie(t *testing.T) {
	docs.Given("given the standard test configuration")
	tstSetup(tstDefaultConfigFile)
	defer tstShutdown()

	docs.When("when a logged in user calls the frontend-userinfo endpoint with a valid token but supplies an invalid access token")
	response := tstPerformGetWithCookies("/v1/frontend-userinfo", valid_JWT_id_is_not_staff_sub101, "wrong-value")

	docs.Then("then this is NOT detected, because the frontend-userinfo endpoint does not talk to the IDP for performance reasons")
	tstRequireUserinfoResponse(t, response, expected_response_by_token[valid_JWT_id_is_not_staff_sub101])

	docs.Then("and no calls to the IDP have been made")
	require.Empty(t, idpMock.recording)
}

func TestFrontendUserinfo_IDPDown(t *testing.T) {
	docs.Given("given the standard test configuration")
	tstSetup(tstDefaultConfigFile)
	defer tstShutdown()

	docs.When("when a logged in user calls the userinfo endpoint while the idp is down")
	response := tstPerformGetWithCookies("/v1/frontend-userinfo", valid_JWT_id_is_not_staff_sub101, "idp_is_down")

	docs.Then("then this does NOT affect the request, which is successful")
	tstRequireUserinfoResponse(t, response, expected_response_by_token[valid_JWT_id_is_not_staff_sub101])
}
