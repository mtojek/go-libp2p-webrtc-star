package testutils

import (
	"github.com/multiformats/go-multiaddr"
	"github.com/stretchr/testify/require"
	"testing"
)

func MustCreateSignalAddr(t *testing.T, signalAddr string) multiaddr.Multiaddr {
	starMultiaddr, err := multiaddr.NewMultiaddr(signalAddr)
	require.NoError(t, err)
	return starMultiaddr
}