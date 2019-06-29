package examples

import (
	"context"
	"github.com/libp2p/go-tcp-transport"
	"sync"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/mtojek/go-libp2p-webrtc-star"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	wss "github.com/mtojek/go-wss-transport"
)

var helloWorldMessage = []byte("Hello world!")

func TestSendSingleMessage(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	firstHost := mustCreateNewHost(t, ctx)
	secondHost := mustCreateNewHost(t, ctx)

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

	firstHost.Peerstore().AddAddr(secondHost.ID(), starMultiaddr, 3600 * time.Second)

	firstHostStream, err := firstHost.NewStream(ctx, secondHost.ID(), protocolID)
	require.NoError(t, err)

	// when
	n, err := firstHostStream.Write(helloWorldMessage)
	require.NoError(t, err)

	require.NotZero(t, n, "no data written")
	wg.Wait()
}

func mustCreateNewHost(t *testing.T, ctx context.Context) host.Host {
	h, err := libp2p.New(ctx,
		libp2p.Transport(wss.New),
		libp2p.Transport(tcp.NewTCPTransport),
		libp2p.Transport(star.New),
		libp2p.DefaultMuxers,
		libp2p.DefaultSecurity)
	require.NoError(t, err)
	return h
}