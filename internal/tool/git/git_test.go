package git

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExplicitRefOrHead(t *testing.T) {
	require.Equal(t, "refs/heads/implicit", ExplicitRefOrHead("implicit"), "Branch name becomes explicit")
	require.Equal(t, "refs/heads/explicit", ExplicitRefOrHead("refs/heads/explicit"), "Explicit ref is left intact")
	require.Equal(t, "HEAD", ExplicitRefOrHead(""), "Empty ref becomes HEAD")
}
