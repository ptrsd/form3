package form3

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewDefaultClient(t *testing.T) {
	client := NewDefaultClient(nil)

	require.NotNil(t, client)
	require.NotNil(t, client.httpClient)

	require.Equal(t, "http://localhost:8080/", client.BaseURL.String())
	require.Equal(t, "form3-client/v1", client.UserAgent)
}
