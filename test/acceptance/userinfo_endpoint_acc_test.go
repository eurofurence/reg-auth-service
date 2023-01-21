package acceptance

import (
	"github.com/eurofurence/reg-auth-service/docs"
	"github.com/eurofurence/reg-auth-service/internal/api/v1/userinfo"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

const valid_JWT_sub1234567890 = `eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwiZ2xvYmFsIjp7Im5hbWUiOiJKb2huIERvZSIsImVtYWlsIjoianNxdWlycmVsX2dpdGh1Yl85YTZkQHBhY2tldGxvc3MuZGUiLCJyb2xlcyI6W119LCJpYXQiOjE1MTYyMzkwMjJ9.lQlx5rhWIwkISQr1KjymKFUU4pqB0JitpaJ2fverqQNSEjzCJUnmkxvmhFgXkTmq31LOmeJd1w5Ijhbx3D3q8R-yz4CHkz0y_OR9BKklp5a-A5oH0s2mm0a6nVQBpjZxhgLv7aYJADCdoGPckak5oXHnqx8_nwKnQkW-BLRbYabvqugTPuv_p_bMVUGMWqMvIuZ9ywKALhzW5Moq80PdIR732Kj9wWd5HkEzve1vhKlJIN77Qsv_KP6mfldafzkS4dsdQLmTfunqtbKAbVRET9ZYVL-gD38OuJjGuDuhTUjV6h9holnmhisEZmrFnjBxxcvxJz0X036WIM_plV4A7A`
const valid_JWT_sub101 = `eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMDEiLCJnbG9iYWwiOnsibmFtZSI6Ik5vcm1hbG8gMTAxIiwiZW1haWwiOiJqc3F1aXJyZWxfZ2l0aHViXzlhNmRAcGFja2V0bG9zcy5kZSIsInJvbGVzIjpbXX0sImlhdCI6MTUxNjIzOTAyMn0.Z-CQB8yy_2seaLGJOO-XBpy1RDOUxbVXczxJPoyylZkB0wTKvzfzS8W7RdCMb2bWPJZZy_2CQa_7mDSfkP_JHX-preW8JxKLWJTuhzFxGec6bLwI_Ri9NRxgVX_hSEKjpm2QhSCxxFe6rXkU8ylRC_B1spCZVadE3EJgnvz0MCRHJk9abufJrvHmp5s4RqCmpXffgNVGsFYrDyoYVaTE5wmFSdA_OFSBslT6wVm5RHjyDUGtoFwEHUMQxlaaF6DXqNGMFcLFVtB059ZQpKyY14a_RWI4PxnXkgRCnZaKv1hmw6wfSY49hSNt6IaSDB9NUUM2LqspVwybVSCTNNcAIg`

const valid_JWT_staff_sub1234567890 = `eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwiZ2xvYmFsIjp7Im5hbWUiOiJKb2huIERvZSIsImVtYWlsIjoianNxdWlycmVsX2dpdGh1Yl85YTZkQHBhY2tldGxvc3MuZGUiLCJyb2xlcyI6WyJzdGFmZiJdfSwiaWF0IjoxNTE2MjM5MDIyfQ.kFLqLgdoS28CRN_I2qCe_L3IkmoynjXGSxIXZE2BQKlZXiHNWjit8Ikz3Yj4j8cSfGscII3tR08aDl10fPRXH08aVEicQy-XGVZYQnn7rAOKaDYMmtVBV3ovELxE0RBFcCeIvuW5GEhBJJ5u3vumJyQVhbnO_PDPyXUbctD8A7g9IAKcPu9AQ5yda0df_tq2iCyUVYHHGrhds6DwnBoRnb1N8dCJ79VbSBoJf96SZfkWiA9IczlJ7CcprTYufbZqVha0tB8JKMDY16ZkziR7A8l8qDPsF-oheCpB_mbjJGRdeBxuviAiOGsS3Tj4aqcTg7QDAhMBlwRpTTt8VpBX3Q`
const valid_JWT_staff_sub202 = `eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIyMDIiLCJnbG9iYWwiOnsibmFtZSI6IkpvaG4gU3RhZmYiLCJlbWFpbCI6ImpzcXVpcnJlbF9naXRodWJfOWE2ZEBwYWNrZXRsb3NzLmRlIiwicm9sZXMiOlsic3RhZmYiXX0sImlhdCI6MTUxNjIzOTAyMn0.rhHe4rsUinOUzztBKTpSTAiZbXaH_EfKj87ZvSvjF1qGq4MV_JpSZ0f-8W8TUAQ3LMs6CUkVQ777b3sPVYFhvB5QYDhot12jGe6YYHTWGYVE5YixbY0EHBv9ztY0Zi2F40uUyVf5d5Nrqe4QZM7vSzu8PUROnQuQos6W3jwUxnTVLRBEjD7J9JvOjzxqXUiR2FDtW1r1P8H2ap22jB9SN_H7iuVtU1EEU9wpO_T8H8jFTo7kHLQdHNd6e9_TqDS7He2wEk0vXd46RdNdKUBDfGPjtHEctOqjCN2juTXi53d1ACqe2EkPr3ih_IP2UWBz5WLlLiGk200nVSZesJTyoQ`

const valid_JWT_admin_sub1234567890 = `eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwiZ2xvYmFsIjp7Im5hbWUiOiJKb2huIERvZSIsImVtYWlsIjoianNxdWlycmVsX2dpdGh1Yl85YTZkQHBhY2tldGxvc3MuZGUiLCJyb2xlcyI6WyJhZG1pbiJdfSwiaWF0IjoxNTE2MjM5MDIyfQ.Phm8d5mDpIX2X4zeajFMvmLKcaLiAAzRS8G3PvNXFyJJTHDrZGZaMMi0aiMGwxRNnBEUVGXCht526rAz9NL3RB222p44gTIdJMujybx37VszF5mVHcxLLxxzu31lcO0W4p1_cwd-ZJlubVMKZ5Yc03BdnYv9KBCVnv6K1xIt_52Izgwl0rESllkEFHflebS78izYpUGHsQjWGXbl_-3HD5fOLhiC5Ixiv2Gq57AHhvL4UwEkAJux65T9CM-ToVdbZutDNDgisCTFvqBY9PG00Js1rTu3BY_2H8-EyFg0Vun2GbjYPaT8PqinBJ7z4kD4os5v3VZh0-tmdR6cahksxg`

// -----------------------------------------
// acceptance tests for the userinfo endpoint
// -----------------------------------------

func TestUserinfo_Success(t *testing.T) {
	docs.Given("given the standard test configuration")
	tstSetup(tstDefaultConfigFile)
	defer tstShutdown()

	docs.When("when a logged in user calls the userinfo endpoint with a valid token")
	response := tstPerformGetWithCookies("/v1/userinfo", valid_JWT_sub1234567890, "access_mock_value")

	docs.Then("then the request is successful and the response is as expected")
	require.Equal(t, http.StatusOK, response.status, "unexpected http response status")
	expectedResponse := userinfo.UserInfoDto{
		Email: "jsquirrel_github_9a6d@packetloss.de",
		Roles: []string{},
	}
	actualResponse := userinfo.UserInfoDto{}
	tstParseJson(response.body, &actualResponse)
	require.EqualValues(t, expectedResponse, actualResponse, "response did not match")
}

func TestUserinfo_Success_Admin(t *testing.T) {
	docs.Given("given the standard test configuration")
	tstSetup(tstDefaultConfigFile)
	defer tstShutdown()

	docs.When("when a logged in admin calls the userinfo endpoint with a valid token")
	response := tstPerformGetWithCookies("/v1/userinfo", valid_JWT_admin_sub1234567890, "access_mock_value")

	docs.Then("then the request is successful and the response is as expected")
	require.Equal(t, http.StatusOK, response.status, "unexpected http response status")
	expectedResponse := userinfo.UserInfoDto{
		Email: "jsquirrel_github_9a6d@packetloss.de",
		Roles: []string{"admin"},
	}
	actualResponse := userinfo.UserInfoDto{}
	tstParseJson(response.body, &actualResponse)
	require.EqualValues(t, expectedResponse, actualResponse, "response did not match")
}

func TestUserinfo_NoAccessTokenCookie(t *testing.T) {
	docs.Given("given the standard test configuration")
	tstSetup(tstDefaultConfigFile)
	defer tstShutdown()

	docs.When("when a logged user calls the userinfo endpoint with a valid token but supplies no access token")
	response := tstPerformGetWithCookies("/v1/userinfo", valid_JWT_sub1234567890, "")

	docs.Then("then the request fails with the appropriate error")
	tstRequireErrorResponse(t, response, http.StatusUnauthorized, "auth.unauthorized", "you did not provide a valid access token - see log for details")
}

func TestUserinfo_InvalidAccessTokenCookie(t *testing.T) {
	docs.Given("given the standard test configuration")
	tstSetup(tstDefaultConfigFile)
	defer tstShutdown()

	docs.When("when a logged in user calls the userinfo endpoint with a valid token but supplies an invalid access token")
	response := tstPerformGetWithCookies("/v1/userinfo", valid_JWT_sub1234567890, "wrong-value")

	docs.Then("then the request fails with the appropriate error")
	tstRequireErrorResponse(t, response, http.StatusUnauthorized, "auth.unauthorized", "identity provider rejected your token - see log for details")
}

func TestUserinfo_SubjectMismatch(t *testing.T) {
	docs.Given("given the standard test configuration")
	tstSetup(tstDefaultConfigFile)
	defer tstShutdown()

	docs.When("when a logged in user calls the userinfo endpoint with a valid token with a different subject")
	response := tstPerformGetWithCookies("/v1/userinfo", valid_JWT_staff_sub202, "access_mock_value")

	docs.Then("then the request fails with the appropriate error")
	tstRequireErrorResponse(t, response, http.StatusUnauthorized, "auth.unauthorized", "identity provider rejected your token - see log for details")
}

func TestUserinfo_MissingRole(t *testing.T) {
	docs.Given("given the standard test configuration")
	tstSetup(tstDefaultConfigFile)
	defer tstShutdown()

	docs.When("when a logged in user calls the userinfo endpoint with a valid token with an extra relevant role (staff) that the idp does not confirm")
	response := tstPerformGetWithCookies("/v1/userinfo", valid_JWT_staff_sub1234567890, "access_mock_value")

	docs.Then("then the request fails with the appropriate error")
	tstRequireErrorResponse(t, response, http.StatusUnauthorized, "auth.unauthorized", "identity provider rejected your token - see log for details")
}

func TestUserinfo_IDPDown(t *testing.T) {
	docs.Given("given the standard test configuration")
	tstSetup(tstDefaultConfigFile)
	defer tstShutdown()

	docs.When("when a logged in user calls the userinfo endpoint while the idp is down")
	response := tstPerformGetWithCookies("/v1/userinfo", valid_JWT_sub1234567890, "idp_is_down")

	docs.Then("then the request fails with the appropriate error")
	tstRequireErrorResponse(t, response, http.StatusBadGateway, "auth.idp.error", "identity provider could not be reached - see log for details")
}
