package acceptance

import (
	"encoding/json"
	"github.com/eurofurence/reg-auth-service/internal/api/v1/errorapi"
	"github.com/go-http-utils/headers"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"testing"
	"time"
)

type tstWebResponse struct {
	status      int
	body        string
	contentType string
	location    string
}

func tstPerformGetNoRedirect(relativeUrlWithLeadingSlash string) http.Response {
	request, err := http.NewRequest(http.MethodGet, ts.URL+relativeUrlWithLeadingSlash, nil)
	if err != nil {
		log.Fatal(err)
	}
	// create a client that doesn't follow redirects
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	return *response
}

func tstPerformGetWithCookies(relativeUrlWithLeadingSlash string, idToken string, accToken string) tstWebResponse {
	request, err := http.NewRequest(http.MethodGet, ts.URL+relativeUrlWithLeadingSlash, nil)
	if err != nil {
		log.Fatal(err)
	}
	expire := time.Now().AddDate(0, 0, 1)
	idCookie := http.Cookie{
		Name:       "JWT",
		Value:      idToken,
		Path:       "/",
		Domain:     "localhost",
		Expires:    expire,
		RawExpires: expire.Format(time.UnixDate),
		MaxAge:     86400,
		Secure:     true,
		HttpOnly:   true,
		SameSite:   http.SameSiteStrictMode,
		Raw:        "test=tcookie",
		Unparsed:   []string{"test=tcookie"},
	}
	accCookie := http.Cookie{
		Name:       "AUTH",
		Value:      accToken,
		Path:       "/",
		Domain:     "localhost",
		Expires:    expire,
		RawExpires: expire.Format(time.UnixDate),
		MaxAge:     86400,
		Secure:     true,
		HttpOnly:   true,
		SameSite:   http.SameSiteStrictMode,
		Raw:        "test=tcookie",
		Unparsed:   []string{"test=tcookie"},
	}

	if idToken != "" {
		request.AddCookie(&idCookie)
	}
	if accToken != "" {
		request.AddCookie(&accCookie)
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	return tstWebResponseFromResponse(response)
}

func tstWebResponseFromResponse(response *http.Response) tstWebResponse {
	status := response.StatusCode
	ct := ""
	if val, ok := response.Header[headers.ContentType]; ok {
		ct = val[0]
	}
	loc := ""
	if val, ok := response.Header[headers.Location]; ok {
		loc = val[0]
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = response.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	return tstWebResponse{
		status:      status,
		body:        string(body),
		contentType: ct,
		location:    loc,
	}
}

func tstResponseBodyString(response *http.Response) string {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "error"
	}

	err = response.Body.Close()
	if err != nil {
		return "error"
	}

	return string(body)
}

func tstParseJson(body string, dto interface{}) {
	err := json.Unmarshal([]byte(body), dto)
	if err != nil {
		log.Fatal(err)
	}
}

func tstRequireErrorResponse(t *testing.T, response tstWebResponse, expectedStatus int, expectedMessage string, expectedDetails interface{}) {
	require.Equal(t, expectedStatus, response.status, "unexpected http response status")
	errorDto := errorapi.ErrorDto{}
	tstParseJson(response.body, &errorDto)
	require.Equal(t, expectedMessage, errorDto.Message, "unexpected error code")
	expectedDetailsStr, ok := expectedDetails.(string)
	if ok && expectedDetailsStr != "" {
		require.EqualValues(t, url.Values{"details": []string{expectedDetailsStr}}, errorDto.Details, "unexpected error details")
	}
	expectedDetailsUrlValues, ok := expectedDetails.(url.Values)
	if ok {
		require.EqualValues(t, expectedDetailsUrlValues, errorDto.Details, "unexpected error details")
	}
}
