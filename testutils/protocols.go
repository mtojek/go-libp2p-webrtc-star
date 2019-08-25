package testutils

import (
	ma "github.com/multiformats/go-multiaddr"
	"github.com/stretchr/testify/require"
	"testing"
)

func MustAddProtocol(t *testing.T, protocol ma.Protocol) {
	err := ma.AddProtocol(protocol)
	require.NoError(t, err)
}

func MustCreateProtocol(code int, name string) ma.Protocol {
	return ma.Protocol{
		Code:  code,
		Name:  name,
		VCode: ma.CodeToVarint(code),
	}
}
