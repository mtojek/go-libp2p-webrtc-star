package peermessaging

import (
	"context"
	"io"
	"sync"
	"testing"
	"time"

	golog "github.com/ipfs/go-log"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/mtojek/go-libp2p-webrtc-star/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	peerMessagingSendSingleMessageProtocolID = protocol.ID("/peer-messaging-send-single-message/1.0.0")
	waitForStreamTimeout                     = 5 * time.Minute

	localSignalAddr  = "/dns4/localhost/tcp/9090/ws/p2p-webrtc-star"
	//remoteSignalAddr = "/dns4/wrtc-star.discovery.libp2p.io/tcp/443/wss/p2p-webrtc-star"
)

var (
	helloWorldMessage     = []byte("Hello world!")
	helloWorldMessageSize = len(helloWorldMessage)
)

func init() {
	golog.SetDebugLogging()
}

func TestSendSingleMessage(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	protocolID := peerMessagingSendSingleMessageProtocolID

	firstHost := testutils.MustCreateHost(t, ctx, localSignalAddr)
	secondHost := testutils.MustCreateHost(t, ctx, localSignalAddr)

	var wg sync.WaitGroup
	wg.Add(1)
	secondHost.SetStreamHandler(peerMessagingSendSingleMessageProtocolID, func(stream network.Stream) {
		message := make([]byte, helloWorldMessageSize)

		n, err := io.ReadFull(stream, message)
		require.NoError(t, err)
		require.NotZero(t, n, "no data read")

		// then
		assert.Equal(t, helloWorldMessage, message)
		wg.Done()
	})

	firstHostStream := testutils.WaitForStream(t, func() (network.Stream, error) {
		return firstHost.NewStream(ctx, secondHost.ID(), protocolID)
	}, waitForStreamTimeout)

	// when
	n, err := firstHostStream.Write(helloWorldMessage)
	require.NoError(t, err)

	wg.Wait()

	require.NotZero(t, n, "no data written")
}
