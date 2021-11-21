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
	// TODO can the IDP issue tokens using JSON Content Type? Have to try this!
	tstAuthorizationCode := "PzL34r_rsNOR7pvdDVbGQRdETOSSU3Tya6Z6AbE0FHFY4AtC"
	tstPkceVerifier := "pCAqaUVKzzeSRyp5L_ydTk38E-4PwSzJ459Xq65rrVe809vd"
	tstRequest := idpclient.TokenRequestDto{
		GrantType:    "authorization_code",
		ClientId:     "democlient",
		ClientSecret: "democlientsecret",
		// RedirectUri:  "https://mydemoapp/landing",
		Code:         tstAuthorizationCode,
		CodeVerifier: tstPkceVerifier,
	}

	tstExpectedResponse := idp.TokenResponseDto{
		TokenType: "Bearer",
		ExpiresIn: 86400,
		AccessToken: "Qs9QaRHGEiSC8FwerVUSijduguq0ZlqMrOiX6Tbya7CMpyCkrQ7TK7ol9WVMis6Ul_6Nm5XV",
		Scope: "example",
		RefreshToken: "EXXVAtYSjg2ZX1q7u6wc8lXH",
	}

	// Pass in test case (consumer side)
	// This uses the repository on the consumer side to make the http call, should be as low level as possible
	var test = func() (err error) {
		// initialize test configuration so we will talk to pact
		_loadContractTestConfig(pact.Server.Port)

		ctx := context.TODO()

		client := idpclient.New()
		actualResponse, err := client.TokenWithAuthenticationCodeAndPKCE(ctx, "example-service", tstAuthorizationCode, tstPkceVerifier)
		if err != nil {
			return err
		}

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
				headers.ContentType:   dsl.String(media.ContentTypeApplicationJson),
			},
			Path: dsl.String("/token"),
			Body: tstRequest,
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
