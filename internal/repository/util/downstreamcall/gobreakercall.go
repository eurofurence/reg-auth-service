package downstreamcall

import (
	"context"
	"fmt"
	"github.com/eurofurence/reg-auth-service/internal/repository/logging"
	"github.com/sony/gobreaker"
	"net/http"
	"time"
)

// https://medium.com/@_jesus_rafael/making-http-client-more-resilient-in-go-d24c66a64bd1

// see also https://github.com/sony/gobreaker

var cb *gobreaker.CircuitBreaker

// ConfigureGobreakerCommand configures
func ConfigureGobreakerCommand(commandName string) {
	// timeout in Gobreaker config is not the request timeout -> not using the parameter here
	var maxNumberRequestsInHalfopenState uint32 = 100
	var counterClearingIntervalWhileClosed time.Duration = 5 * time.Minute
	var timeUntilHalfopenAfterOpen time.Duration = 60 * time.Second
	// default ReadyToTrip opens after 5 consecutive failures
	// default IsSuccessful returns false for all non-nil errors
	settings := gobreaker.Settings{
		Name:          commandName,
		MaxRequests:   maxNumberRequestsInHalfopenState,
		Interval:      counterClearingIntervalWhileClosed,
		Timeout:       timeUntilHalfopenAfterOpen,
		ReadyToTrip:   nil,
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			logging.NoCtx().Warn(fmt.Sprintf("Circuitbreaker %s state change %s -> %s", name, from.String(), to.String()))
		},
		IsSuccessful:  nil,
	}
	cb = gobreaker.NewCircuitBreaker(settings)
	logging.NoCtx().Info(fmt.Sprintf("Circuitbreaker %s set up", commandName))
}

// GobreakerPerformPOST performs a http POST, returning the response body and status and passing on the request id if present in the context.
//
// The request is wrapped with a circuit breaker and timeout.
//
// Note: you must make at least one call to ConfigureGobreakerCommand() before calling this.
func GobreakerPerformPOST(ctx context.Context, httpClient *http.Client, url string, requestBody string, contentType string) (string, int, error) {
	return gobreakerPerformWithBody(ctx, http.MethodPost, httpClient, url, requestBody, contentType)
}

// GobreakerPerformPUT performs a http PUT, returning the response body and status and passing on the request id if present in the context.
//
// The request is wrapped with a circuit breaker and timeout.
//
// Note: you must make at least one call to ConfigureGobreakerCommand() before calling this.
func GobreakerPerformPUT(ctx context.Context, httpClient *http.Client, url string, requestBody string, contentType string) (string, int, error) {
	return gobreakerPerformWithBody(ctx, http.MethodPut, httpClient, url, requestBody, contentType)
}

// GobreakerPerformGET performs a http GET, returning the response body and status and passing on the request id if present in the context.
//
// The request is wrapped with a circuit breaker and timeout.
//
// Note: you must make at least one call to ConfigureGobreakerCommand() before calling this.
func GobreakerPerformGET(ctx context.Context, httpClient *http.Client, url string) (string, int, error) {
	return gobreakerPerformNoBody(ctx, http.MethodGet, httpClient, url)
}

// --- internal helper functions ---

func gobreakerPerformNoBody(ctx context.Context, method string, httpClient *http.Client, url string) (string, int, error) {
	return gobreakerPerformWithBody(ctx, method, httpClient, url, "", "")
}

func gobreakerPerformWithBody(ctx context.Context, method string, httpClient *http.Client, url string, requestBody string, contentType string) (string, int, error) {
	responseUntyped, err := cb.Execute(func() (interface{}, error) {
		responseBody, httpStatus, innerErr := performWithBody(ctx, method, httpClient, url, requestBody, contentType)

		// if we return an error at this point, it will count towards opening the circuit breaker
		if innerErr != nil {
			// this may need some more attention
			return nil, innerErr
		}
		if httpStatus >= 500 {
			// so let's make sure any http status in the 500 range causes us to return an error
			innerErr = fmt.Errorf("got unexpected http status %d", httpStatus)
			return nil, innerErr
		}

		return responseInfo{
			body:   responseBody,
			status: httpStatus,
		}, nil
	})
	if err != nil {
		return "", 500, err
	}
	response, ok := responseUntyped.(responseInfo)
	if !ok {
		return "", 500, fmt.Errorf("got no response data structure despite no error")
	}
	return response.body, response.status, err
}
