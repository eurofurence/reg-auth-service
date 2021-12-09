package consumer

import (
	"context"
	"github.com/eurofurence/reg-auth-service/internal/repository/idp"
	"github.com/eurofurence/reg-auth-service/internal/repository/idp/idpclient"
	"github.com/eurofurence/reg-auth-service/web/util/media"
	"github.com/go-http-utils/headers"
	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/require"
	"log"
	"net/http"
	"testing"
)

// contract test consumer side

func TestConsumer(t *testing.T) {
	// Create Pact connecting to local Daemon
	pact := &dsl.Pact{
		Consumer: "RegAuthService",
		Provider: "IdentityProvider",
		Host:     "localhost",
	}
	defer pact.Teardown()

	// types and values used in interaction
	tstAuthorizationCode := "PzL34r_rsNOR7pvdDVbGQRdETOSSU3Tya6Z6AbE0FHFY4AtC"
	tstPkceVerifier := "pCAqaUVKzzeSRyp5L_ydTk38E-4PwSzJ459Xq65rrVe809vd"
	tstRequestBody := "client_id=democlient" +
		"&client_secret=democlientsecret" +
		"&code=" + tstAuthorizationCode +
		"&code_verifier=" + tstPkceVerifier +
		"&grant_type=authorization_code" +
		"&redirect_uri=https%3A%2F%2Fexample.com%2Fapp%2F"

	tstExpectedResponse := idp.TokenResponseDto{
		AccessToken: "XYZ",
		ExpiresIn: 86400,
		IdToken: "abc",
		Scope: "example",
		TokenType: "Bearer",
	}

	// Pass in test case (consumer side)
	// This uses the repository on the consumer side to make the http call, should be as low level as possible
	var test = func() (err error) {
		// initialize test configuration so we will talk to pact
		_loadContractTestConfig(pact.Server.Port)

		ctx := context.Background()

		client := idpclient.New()
		actualResponse, httpstatus, err := client.TokenWithAuthenticationCodeAndPKCE(ctx, "example-service", tstAuthorizationCode, tstPkceVerifier)
		if err != nil {
			return err
		}

		require.Equal(t, http.StatusOK, httpstatus)
		require.EqualValues(t, tstExpectedResponse, *actualResponse, "token response did not match")
		return nil
	}

	// Set up our expected interactions.
	pact.
		AddInteraction().
		// contrived example, not really needed. This is the identifier of the state handler that will be called on the other side
		Given("a client for the given configuration exists in the IDP").
		Given("the user has completed the authorization_code flow with PKCE, obtaining an authorization code").
		UponReceiving("a request to the token endpoint using this authorization code with correct client credentials and correct code verifier").
		WithRequest(dsl.Request{
			Method: http.MethodPost,
			Headers: dsl.MapMatcher{
				headers.ContentType:   dsl.String(media.ContentTypeApplicationXWwwFormUrlencoded),
			},
			Path: dsl.String("/token"),
			Body: tstRequestBody,
		}).
		WillRespondWith(dsl.Response{
			Status:  200,
			Headers: dsl.MapMatcher{headers.ContentType: dsl.String(media.ContentTypeApplicationJson)},
			Body:    tstExpectedResponse,
		})

	// Run the test, verify it did what we expected and capture the contract (writes a test log to logs/pact.log)
	if err := pact.Verify(test); err != nil {
		log.Fatalf("Error on Verify: %v", err)
	}

	// now write out the contract json (by default it goes to subdirectory pacts)
	if err := pact.WritePact(); err != nil {
		log.Fatalf("Error on pact write: %v", err)
	}
}
