package downstreamcall

import (
	"context"
	"github.com/eurofurence/reg-auth-service/web/util/media"
	"github.com/go-http-utils/headers"
	"io/ioutil"
	"net/http"
	"strings"
)

// PerformPOST performs a http POST, returning the response body and status and passing on the request id if present in the context
func PerformPOST(ctx context.Context, httpClient *http.Client, url string, requestBody string, contentType string) (string, int, error) {
	return performWithBody(ctx, http.MethodPost, httpClient, url, requestBody, contentType)
}

// PerformPUT performs a http PUT, returning the response body and status and passing on the request id if present in the context
func PerformPUT(ctx context.Context, httpClient *http.Client, url string, requestBody string, contentType string) (string, int, error) {
	return performWithBody(ctx, http.MethodPut, httpClient, url, requestBody, contentType)
}

// PerformGET performs a http GET, returning the response body and status and passing on the request id if present in the context
func PerformGET(ctx context.Context, httpClient *http.Client, url string) (string, int, error) {
	return performNoBody(ctx, http.MethodGet, httpClient, url)
}

// --- internal helper functions ---

func performNoBody(ctx context.Context, method string, httpClient *http.Client, url string) (string, int, error) {
	return performWithBody(ctx, method, httpClient, url, "", "")
}

func performWithBody(ctx context.Context, method string, httpClient *http.Client, url string, requestBody string, contentType string) (string, int, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, strings.NewReader(requestBody))
	if err != nil {
		return "", 0, err
	}

	if requestBody != "" {
		if contentType == "" {
			req.Header.Set(headers.ContentType, media.ContentTypeApplicationJson)
		} else {
			req.Header.Set(headers.ContentType, contentType)
		}
	}

	response, err := httpClient.Do(req)
	if err != nil {
		return "", 0, err
	}

	responseBody, err := responseBodyString(response)
	return responseBody, response.StatusCode, err
}

func responseBodyString(response *http.Response) (string, error) {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	err = response.Body.Close()
	if err != nil {
		return "", err
	}

	return string(body), nil
}
