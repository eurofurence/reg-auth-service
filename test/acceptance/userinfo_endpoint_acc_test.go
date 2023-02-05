package acceptance

import (
	"github.com/eurofurence/reg-auth-service/docs"
	"github.com/eurofurence/reg-auth-service/internal/api/v1/userinfo"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

const valid_JWT_id_is_not_staff_sub101 = `eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJhdF9oYXNoIjoidDdWYkV5NVQ3STdYSlh3VHZ4S3hLdyIsImF1ZCI6WyIxNGQ5ZjM3YS0xZWVjLTQ3YzktYTk0OS01ZjFlYmRmOWM4ZTUiXSwiYXV0aF90aW1lIjoxNTE2MjM5MDIyLCJlbWFpbCI6ImpzcXVpcnJlbF9naXRodWJfOWE2ZEBwYWNrZXRsb3NzLmRlIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsImV4cCI6MjA3NTEyMDgxNiwiZ3JvdXBzIjpbInNvbWVncm91cCJdLCJpYXQiOjE1MTYyMzkwMjIsImlzcyI6Imh0dHA6Ly9pZGVudGl0eS5sb2NhbGhvc3QvIiwianRpIjoiNDA2YmUzZTQtZjRlOS00N2I3LWFjNWYtMDZiOTI3NDMyODQ4IiwibmFtZSI6IkpvaG4gRG9lIiwibm9uY2UiOiIzMGM4M2MxM2M5MTc5ODA0YWEwZjliMzkzNDI1OWQ3NSIsInJhdCI6MTY3NTExNzE3Nywic2lkIjoiZDdiOGZlN2EtMDc5YS00NTk2LThlNTMtYTYwZjg2YTA4YWM2Iiwic3ViIjoiMTAxIn0.ntHz3G7LLtHC3pJ1PoWJoG3mnzg96IIcP3LAV4V1CcKYMFoKVQfh7MiOdRXpiB-_j4QFE7O-za3mynwFqRbF3_Tw_Sp7Zsgk9OUPo2Mk3VBSl9yPIU4pmc8v7nrmaAVOQLyjglVG7NLRWLpx0oIG8SSN0d75PBI5iLyQ0H7Zu0npEu6xekHeAYAg9DHQxqZInzom72aLmHdtG7tOqOgN0XphiK7zmIqm5aCg7R9_J9s0UU0g16_Phxm3DaynufGCjEPE2YrSL7hY9UVT2nfrHO7MvVOEKMG3RaKUDjzqOkLawz9TcUJlUTBc1J-91zYbdXLHYT_2b4EW_qa1C-P3Ow`

const valid_JWT_id_is_staff_sub202 = `eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJhdF9oYXNoIjoidDdWYkV5NVQ3STdYSlh3VHZ4S3hLdyIsImF1ZCI6WyIxNGQ5ZjM3YS0xZWVjLTQ3YzktYTk0OS01ZjFlYmRmOWM4ZTUiXSwiYXV0aF90aW1lIjoxNTE2MjM5MDIyLCJlbWFpbCI6ImpzcXVpcnJlbF9naXRodWJfOWE2ZEBwYWNrZXRsb3NzLmRlIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsImV4cCI6MjA3NTEyMDgxNiwiZ3JvdXBzIjpbInN0YWZmIl0sImlhdCI6MTUxNjIzOTAyMiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LmxvY2FsaG9zdC8iLCJqdGkiOiI0MDZiZTNlNC1mNGU5LTQ3YjctYWM1Zi0wNmI5Mjc0MzI4NDgiLCJuYW1lIjoiSm9obiBTdGFmZiIsIm5vbmNlIjoiMzBjODNjMTNjOTE3OTgwNGFhMGY5YjM5MzQyNTlkNzUiLCJyYXQiOjE2NzUxMTcxNzcsInNpZCI6ImQ3YjhmZTdhLTA3OWEtNDU5Ni04ZTUzLWE2MGY4NmEwOGFjNiIsInN1YiI6IjIwMiJ9.pM-jMGdjwNvHQMov8JQpRa79CBjHAUHpwElYRvUz_DvhkqcG4SrntVruAlJRS8D9CccflKeTjSEfOiS2l52p0qQ7ZeNPSRQ9nsr_EHDpB7UqcUszqVaBWtIhwkiwca_sxe-8L9A9hPSe_kH9dhDHVbhUsj9vp0HBIV89mtH3i3D6s3quRYtCe9puepkmyf5JC-2TSGoSiyURoFdqXSNRPEuv1FhlyVICqylfkINceCe8dt7lJCrhOc8R-11vA-SRsrBhdxBvYjT29hhQO3eHgJenPufjJPj6kYSWvN91U3KcsffoBmu-C1A7zBLu-zmWBHh4RkYWqbZpNr59TIpRSw`

const valid_JWT_id_is_staff_admin_sub1234567890 = `eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJhdF9oYXNoIjoidDdWYkV5NVQ3STdYSlh3VHZ4S3hLdyIsImF1ZCI6WyIxNGQ5ZjM3YS0xZWVjLTQ3YzktYTk0OS01ZjFlYmRmOWM4ZTUiXSwiYXV0aF90aW1lIjoxNTE2MjM5MDIyLCJlbWFpbCI6ImpzcXVpcnJlbF9naXRodWJfOWE2ZEBwYWNrZXRsb3NzLmRlIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsImV4cCI6MjA3NTEyMDgxNiwiZ3JvdXBzIjpbInN0YWZmIiwiYWRtaW4iXSwiaWF0IjoxNTE2MjM5MDIyLCJpc3MiOiJodHRwOi8vaWRlbnRpdHkubG9jYWxob3N0LyIsImp0aSI6IjQwNmJlM2U0LWY0ZTktNDdiNy1hYzVmLTA2YjkyNzQzMjg0OCIsIm5hbWUiOiJKb2huIEFkbWluIiwibm9uY2UiOiIzMGM4M2MxM2M5MTc5ODA0YWEwZjliMzkzNDI1OWQ3NSIsInJhdCI6MTY3NTExNzE3Nywic2lkIjoiZDdiOGZlN2EtMDc5YS00NTk2LThlNTMtYTYwZjg2YTA4YWM2Iiwic3ViIjoiMTIzNDU2Nzg5MCJ9.DRKPy0Rq-r5-Va6W5coT91JpDV2RkhYjniqIJmmPzOq3LphzRrlDKioDns4ilMxMEpfhFcmv87yOdPsPijUhEqy1a93BeJYMyU7DMdQBtD8R9oYU_-FmqS5dM9ZrBCZZUxTBeNBl2JGI-H1_IIqUH65PodoijO4N5ayw43q5xT1KP7PJKZ9YiMSsa4tUOp0R_Ay51DTIuti21TqqbSCC66sGH_1e1eeuhwBoU7Iws4oeepTRZ_XOdpn_YzTViPs7Byua-zohYgQYthDoCvLVfJOr77BV2vTUrQZfRca7prizXbVuQyxQJEpOBIuke29Gye6Qfbwpb4rMaza3fTLhZg`

var expected_response_by_token = map[string]userinfo.UserInfoDto{
	valid_JWT_id_is_not_staff_sub101: {
		Subject:       "101",
		Name:          "John Doe",
		Email:         "jsquirrel_github_9a6d@packetloss.de",
		EmailVerified: true,
		Groups:        []string{},
	},
	valid_JWT_id_is_staff_sub202: {
		Subject:       "202",
		Name:          "John Staff",
		Email:         "jsquirrel_github_9a6d@packetloss.de",
		EmailVerified: true,
		Groups:        []string{"staff"},
	},
	valid_JWT_id_is_staff_admin_sub1234567890: {
		Subject:       "1234567890",
		Name:          "John Admin",
		Email:         "jsquirrel_github_9a6d@packetloss.de",
		EmailVerified: true,
		Groups:        []string{"admin", "staff"},
	},
}

// -----------------------------------------
// acceptance tests for the userinfo endpoint
// -----------------------------------------

func TestUserinfo_Success(t *testing.T) {
	docs.Given("given the standard test configuration")
	tstSetup(tstDefaultConfigFile)
	defer tstShutdown()

	docs.When("when a logged in user calls the userinfo endpoint with a valid token")
	response := tstPerformGetWithCookies("/v1/userinfo", valid_JWT_id_is_not_staff_sub101, "access_mock_value 101")

	docs.Then("then the request is successful and the response is as expected")
	tstRequireUserinfoResponse(t, response, expected_response_by_token[valid_JWT_id_is_not_staff_sub101])

	docs.Then("and the expected calls to the IDP have been made")
	require.EqualValues(t, []string{"access_mock_value 101"}, idpMock.recording)
}

func TestUserinfo_Success_Admin(t *testing.T) {
	docs.Given("given the standard test configuration")
	tstSetup(tstDefaultConfigFile)
	defer tstShutdown()

	docs.When("when a logged in admin calls the userinfo endpoint with a valid token")
	response := tstPerformGetWithCookies("/v1/userinfo", valid_JWT_id_is_staff_admin_sub1234567890, "access_mock_value")

	docs.Then("then the request is successful and the response is as expected")
	tstRequireUserinfoResponse(t, response, expected_response_by_token[valid_JWT_id_is_staff_admin_sub1234567890])

	docs.Then("and the expected calls to the IDP have been made")
	require.EqualValues(t, []string{"access_mock_value"}, idpMock.recording)
}

func TestUserinfo_NoAccessTokenCookie(t *testing.T) {
	docs.Given("given the standard test configuration")
	tstSetup(tstDefaultConfigFile)
	defer tstShutdown()

	docs.When("when a logged user calls the userinfo endpoint with a valid token but supplies no access token")
	response := tstPerformGetWithCookies("/v1/userinfo", valid_JWT_id_is_not_staff_sub101, "")

	docs.Then("then the request fails with the appropriate error")
	tstRequireErrorResponse(t, response, http.StatusUnauthorized, "auth.unauthorized", "authorization failed to check out during local validation - please see logs for details")

	docs.Then("and no calls to the IDP have been made")
	require.Empty(t, idpMock.recording)
}

func TestUserinfo_InvalidAccessTokenCookie(t *testing.T) {
	docs.Given("given the standard test configuration")
	tstSetup(tstDefaultConfigFile)
	defer tstShutdown()

	docs.When("when a logged in user calls the userinfo endpoint with a valid token but supplies an invalid access token")
	response := tstPerformGetWithCookies("/v1/userinfo", valid_JWT_id_is_not_staff_sub101, "wrong-value")

	docs.Then("then the request fails with the appropriate error")
	tstRequireErrorResponse(t, response, http.StatusUnauthorized, "auth.unauthorized", "identity provider rejected your token - see log for details")

	docs.Then("and the expected calls to the IDP have been made")
	require.EqualValues(t, []string{"wrong-value"}, idpMock.recording)
}

func TestUserinfo_SubjectMismatch(t *testing.T) {
	docs.Given("given the standard test configuration")
	tstSetup(tstDefaultConfigFile)
	defer tstShutdown()

	docs.When("when a logged in user calls the userinfo endpoint with a valid token with a different subject")
	response := tstPerformGetWithCookies("/v1/userinfo", valid_JWT_id_is_staff_sub202, "access_mock_value")

	docs.Then("then the request fails with the appropriate error")
	tstRequireErrorResponse(t, response, http.StatusUnauthorized, "auth.unauthorized", "identity provider rejected your token - see log for details")

	docs.Then("and the expected calls to the IDP have been made")
	require.EqualValues(t, []string{"access_mock_value"}, idpMock.recording)
}

func TestUserinfo_MissingGroup(t *testing.T) {
	docs.Given("given the standard test configuration")
	tstSetup(tstDefaultConfigFile)
	defer tstShutdown()

	docs.When("when a logged in user calls the userinfo endpoint with a valid token with an extra relevant group (staff) that the idp does not confirm")
	response := tstPerformGetWithCookies("/v1/userinfo", valid_JWT_id_is_staff_sub202, "access_mock_value 202") // does not return staff group

	docs.Then("then the request fails with the appropriate error")
	tstRequireErrorResponse(t, response, http.StatusUnauthorized, "auth.unauthorized", "identity provider rejected your token - see log for details")

	docs.Then("and the expected calls to the IDP have been made")
	require.EqualValues(t, []string{"access_mock_value 202"}, idpMock.recording)
}

func TestUserinfo_IDPDown(t *testing.T) {
	docs.Given("given the standard test configuration")
	tstSetup(tstDefaultConfigFile)
	defer tstShutdown()

	docs.When("when a logged in user calls the userinfo endpoint while the idp is down")
	response := tstPerformGetWithCookies("/v1/userinfo", valid_JWT_id_is_not_staff_sub101, "idp_is_down")

	docs.Then("then the request fails with the appropriate error")
	tstRequireErrorResponse(t, response, http.StatusBadGateway, "auth.idp.error", "identity provider could not be reached - see log for details")
}

// --- helpers

func tstRequireUserinfoResponse(t *testing.T, response tstWebResponse, expectedResponse userinfo.UserInfoDto) {
	require.Equal(t, http.StatusOK, response.status, "unexpected http response status")

	actualResponse := userinfo.UserInfoDto{}
	tstParseJson(response.body, &actualResponse)
	require.EqualValues(t, expectedResponse, actualResponse, "response did not match")
}
