package authctl

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestVerifierEncode(t *testing.T) {
	knownVerifier := "pCAqaUVKzzeSRyp5L_ydTk38E-4PwSzJ459Xq65rrVe809vd"
	knownChallenge := "5658wVxOCRqpTL0htZr_j6Ch6c0THQWIhfBqrADyGiA"

	actualChallenge := generateCodeChallenge(knownVerifier)

	require.Equal(t, knownChallenge, actualChallenge)
}
