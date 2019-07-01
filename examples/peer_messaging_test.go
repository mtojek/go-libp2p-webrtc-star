package examples

import (
	"context"
	"sync"
	"testing"
	
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var helloWorldMessage = []byte("Hello world!")

func TestSendSingleMessage(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	firstHost := mustCreateHost(t, ctx)
	secondHost := mustCreateHost(t, ctx)

	var wg sync.WaitGroup
	wg.Add(1)
	secondHost.SetStreamHandler(protocolID, func(stream network.Stream) {
		var message []byte

		n, err := stream.Read(message)
		require.NoError(t, err)
		require.NotZero(t, n, "no data read")

		// then
		assert.Equal(t, helloWorldMessage, message)
		wg.Done()
	})

	firstHost.Peerstore().AddAddr(secondHost.ID(), starMultiaddr, peerstoreAddressTTL)

	firstHostStream, err := firstHost.NewStream(ctx, secondHost.ID(), protocolID)
	require.NoError(t, err)

	wg.Wait()

	// when
	n, err := firstHostStream.Write(helloWorldMessage)
	require.NoError(t, err)

	require.NotZero(t, n, "no data written")
}