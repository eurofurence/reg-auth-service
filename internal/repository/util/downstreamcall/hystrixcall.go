package downstreamcall

import (
	"context"
	"fmt"
	"net/http"

	"github.com/afex/hystrix-go/hystrix"
)

// ConfigureHystrixCommand configures hystrix command timeout, concurrency and error threshold
func ConfigureHystrixCommand(hystrixCommandName string, timeoutMs int) {
	hystrix.ConfigureCommand(hystrixCommandName, hystrix.CommandConfig{
		Timeout:               timeoutMs,
		MaxConcurrentRequests: 100,
		ErrorPercentThreshold: 25,
	})
}

// HystrixPerformPOST performs a http POST, returning the response body and status and passing on the request id if present in the context.
//
// The request is wrapped with a hystrix circuit breaker and timeout.
//
// Note: you must make at least one call to ConfigureHystrixCommand() before calling this.
func HystrixPerformPOST(ctx context.Context, hystrixCommandName string, httpClient *http.Client, url string, requestBody string) (string, int, error) {
	return hystrixPerformWithBody(ctx, hystrixCommandName, http.MethodPost, httpClient, url, requestBody)
}

// HystrixPerformPUT performs a http PUT, returning the response body and status and passing on the request id if present in the context.
//
// The request is wrapped with a hystrix circuit breaker and timeout.
//
// Note: you must make at least one call to ConfigureHystrixCommand() before calling this.
func HystrixPerformPUT(ctx context.Context, hystrixCommandName string, httpClient *http.Client, url string, requestBody string) (string, int, error) {
	return hystrixPerformWithBody(ctx, hystrixCommandName, http.MethodPut, httpClient, url, requestBody)
}

// HystrixPerformGET performs a http GET, returning the response body and status and passing on the request id if present in the context.
//
// The request is wrapped with a hystrix circuit breaker and timeout.
//
// Note: you must make at least one call to ConfigureHystrixCommand() before calling this.
func HystrixPerformGET(ctx context.Context, hystrixCommandName string, httpClient *http.Client, url string) (string, int, error) {
	return hystrixPerformNoBody(ctx, hystrixCommandName, http.MethodGet, httpClient, url)
}

// --- internal helper functions ---

type responseInfo struct {
	body   string
	status int
}

func hystrixPerformNoBody(ctx context.Context, hystrixCommandName string, method string, httpClient *http.Client, url string) (string, int, error) {
	return hystrixPerformWithBody(ctx, hystrixCommandName, method, httpClient, url, "")
}

func hystrixPerformWithBody(ctx context.Context, hystrixCommandName string, method string, httpClient *http.Client, url string, requestBody string) (string, int, error) {
	output := make(chan responseInfo, 1)

	// hystrix.DoC blocks until either completed or error returned
	err := hystrix.DoC(ctx, hystrixCommandName, func(subctx context.Context) error {
		responseBody, httpstatus, innerErr := performWithBody(subctx, method, httpClient, url, requestBody)
		output <- responseInfo{
			body:   responseBody,
			status: httpstatus,
		}

		// if we return an error at this point, it will count towards opening the circuit breaker
		if httpstatus >= 500 && innerErr == nil {
			// so let's make sure any http status in the 500 range causes us to return an error
			// in a real world situation this may need some more attention
			innerErr = fmt.Errorf("got unexpected http status %d", httpstatus)
		}
		return innerErr
	}, nil)

	responseData := responseInfo{}

	// non-blocking receive for optional output
	select {
	case out := <-output:
		responseData = out
	default:
		// presence of default branch means select will not block even if none of the channels are ready to read from
	}

	return responseData.body, responseData.status, err
}
