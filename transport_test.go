package star

import (
	golog "github.com/ipfs/go-log"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/transport"
	"github.com/libp2p/go-libp2p-peerstore/pstoremem"
	"github.com/libp2p/go-libp2p-testing/suites/transport"
	yamux "github.com/libp2p/go-libp2p-yamux"
	"github.com/mtojek/go-libp2p-webrtc-star/testutils"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/pion/webrtc"
	"testing"
)

const starAddress = "/dns4/localhost/tcp/9090/ws/p2p-webrtc-star"

func init() {
	golog.SetDebugLogging()

	wsProtocol := testutils.MustCreateProtocol(wsProtocolCode, "ws")
	testutils.MustAddProtocol(wsProtocol)
}

func TestBasic(t *testing.T) {
	starTransportA, starTransportB, mAddr, identityA := testParameters(t)
	ttransport.SubtestBasic(t, starTransportA, starTransportB, mAddr, identityA)
}

func TestCancel(t *testing.T) {
	starTransportA, starTransportB, mAddr, identityA := testParameters(t)
	ttransport.SubtestCancel(t, starTransportA, starTransportB, mAddr, identityA)
}

func TestPingPong(t *testing.T) {
	starTransportA, starTransportB, mAddr, identityA := testParameters(t)
	ttransport.SubtestCancel(t, starTransportA, starTransportB, mAddr, identityA)
}

func TestStress1Conn1Stream1Msg(t *testing.T) {
	starTransportA, starTransportB, mAddr, identityA := testParameters(t)
	ttransport.SubtestStress1Conn1Stream1Msg(t, starTransportA, starTransportB, mAddr, identityA)
}

func TestStress1Conn1Stream100Msg(t *testing.T) {
	starTransportA, starTransportB, mAddr, identityA := testParameters(t)
	ttransport.SubtestStress1Conn1Stream100Msg(t, starTransportA, starTransportB, mAddr, identityA)
}

func TestStress1Conn100Stream100Msg(t *testing.T) {
	starTransportA, starTransportB, mAddr, identityA := testParameters(t)
	ttransport.SubtestStress1Conn100Stream100Msg(t, starTransportA, starTransportB, mAddr, identityA)
}

func TestStress1Conn1000Stream10Msg(t *testing.T) {
	starTransportA, starTransportB, mAddr, identityA := testParameters(t)
	ttransport.SubtestStress1Conn1000Stream10Msg(t, starTransportA, starTransportB, mAddr, identityA)
}

func TestStress1Conn100Stream100Msg10MB(t *testing.T) {
	starTransportA, starTransportB, mAddr, identityA := testParameters(t)
	ttransport.SubtestStress1Conn100Stream100Msg10MB(t, starTransportA, starTransportB, mAddr, identityA)
}

func TestStreamOpenStress(t *testing.T) {
	starTransportA, starTransportB, mAddr, identityA := testParameters(t)
	ttransport.SubtestStreamOpenStress(t, starTransportA, starTransportB, mAddr, identityA)
}

func TestStreamReset(t *testing.T) {
	starTransportA, starTransportB, mAddr, identityA := testParameters(t)
	ttransport.SubtestStreamReset(t, starTransportA, starTransportB, mAddr, identityA)
}

func testParameters(t *testing.T) (transport.Transport, transport.Transport, ma.Multiaddr, peer.ID) {
	starTransportA, identityA := mustCreateStarTransport(t)
	starTransportB, _ := mustCreateStarTransport(t)

	mAddr, err := ma.NewMultiaddr(starAddress)
	if err != nil {
		t.Fatal(err)
	}
	return starTransportA, starTransportB, mAddr, identityA
}

func mustCreateStarTransport(t *testing.T) (transport.Transport, peer.ID) {
	privKey := testutils.MustCreatePrivateKey(t)
	identity := testutils.MustCreatePeerIdentity(t, privKey)
	peerstore := pstoremem.NewPeerstore()
	muxer := yamux.DefaultTransport
	return New(identity, peerstore).
		WithSignalConfiguration(SignalConfiguration{
			URLPath: "/socket.io/?EIO=3&transport=websocket",
		}).
		WithWebRTCConfiguration(webrtc.Configuration{
			ICEServers: []webrtc.ICEServer{
				{
					URLs: []string{
						"stun:stun.l.google.com:19302",
						"stun:stun1.l.google.com:19302",
						"stun:stun2.l.google.com:19302",
						"stun:stun3.l.google.com:19302",
						"stun:stun4.l.google.com:19302",
					},
				},
			},
		}).
		WithMuxer(muxer), identity
}
