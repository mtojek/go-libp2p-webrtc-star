package star

import (
	golog "github.com/ipfs/go-log"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/transport"
	"github.com/libp2p/go-libp2p-peerstore/pstoremem"
	"github.com/libp2p/go-libp2p-testing/suites/transport"
	"github.com/mtojek/go-libp2p-webrtc-star/testutils"
	"testing"
)

func init() {
	golog.SetDebugLogging()
}

func TestTransport(t *testing.T) {
	wsProtocol := testutils.MustCreateProtocol(wsProtocolCode, "ws")
	testutils.MustAddProtocol(t, wsProtocol)

	starTransportA, identityA := mustCreateStarTransport(t)
	starTransportB, _ := mustCreateStarTransport(t)

	ttransport.SubtestTransport(t,
		starTransportA,
		starTransportB,
		"/dns4/localhost/tcp/9090/ws/p2p-webrtc-star",
		identityA)
}

func mustCreateStarTransport(t *testing.T) (transport.Transport, peer.ID) {
	privKey := testutils.MustCreatePrivateKey(t)
	identity := testutils.MustCreatePeerIdentity(t, privKey)
	peerstore := pstoremem.NewPeerstore()
	return New(identity, peerstore).
		WithSignalConfiguration(SignalConfiguration{
			URLPath: "/socket.io/?EIO=3&transport=websocket",
		}), identity
}
