package examples

import (
	"context"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/mtojek/go-libp2p-webrtc-star"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/stretchr/testify/assert"
)

const protocolID = "/go-libp2p-webrtc-star/1.0.0"

var helloWorldMessage = []byte("Hello world!")

func TestSendSingleMessage(t *testing.T) {
	// given
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	firstHost := mustCreateDefaultHost(t, ctx)
	secondHost := mustCreateDefaultHost(t, ctx)

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

	firstHostStream, err := firstHost.NewStream(ctx, secondHost.ID(), protocolID)
	require.NoError(t, err)

	// when
	n, err := firstHostStream.Write(helloWorldMessage)
	require.NoError(t, err)

	require.NotZero(t, n, "no data written")
	wg.Wait()
}

func mustCreateDefaultHost(t *testing.T, ctx context.Context) host.Host {
	h, err := libp2p.New(ctx,
		libp2p.Transport(star.Transport()))
	require.NoError(t, err)
	return h
}
