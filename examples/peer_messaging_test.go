package examples

import (
	"context"
	"sync"
	"testing"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/mtojek/go-libp2p-webrtc-star"
	"github.com/multiformats/go-multiaddr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	protocolID = "/p2p-webrtc-star/1.0.0"
	starSignalAddr = "/dns4/star-signal.cloud.ipfs.team/http/p2p-webrtc-star"
)

var helloWorldMessage = []byte("Hello world!")

func TestSendSingleMessage(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalAddr := mustCreateSignalAddr(t)
	firstHost := mustCreateNewHost(t, ctx, signalAddr)
	secondHost := mustCreateNewHost(t, ctx, signalAddr)

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

func mustCreateSignalAddr(t *testing.T) multiaddr.Multiaddr {
	starSignal, err := multiaddr.NewMultiaddr(starSignalAddr)
	require.NoError(t, err)
	return starSignal
}

func mustCreateNewHost(t *testing.T, ctx context.Context, signalAddr multiaddr.Multiaddr) host.Host {
	h, err := libp2p.New(ctx,
		libp2p.Transport(star.Transport()),
		libp2p.ListenAddrs(signalAddr))
	require.NoError(t, err)
	return h
}