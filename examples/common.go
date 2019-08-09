package examples

import (
	"context"
	"crypto/rand"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"testing"
	"time"

	golog "github.com/ipfs/go-log"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-peerstore/pstoremem"
	"github.com/mtojek/go-libp2p-webrtc-star"
	"github.com/multiformats/go-multiaddr"
	"github.com/stretchr/testify/require"
)

const (
	protocolID     = "/p2p-webrtc-star/1.0.0"
	//firstSignalAddr = "/dns4/wrtc-star.discovery.libp2p.io/tcp/443/wss/p2p-webrtc-star"
	firstSignalAddr = "/dns4/localhost/tcp/9090/ws/p2p-webrtc-star"
	waitForStreamTimeout = 60 * time.Minute
)

func init() {
	golog.SetDebugLogging()
}

func mustCreateHost(t *testing.T, ctx context.Context) host.Host {
	signalMultiaddr := mustCreateSignalAddr(t, firstSignalAddr)

	privKey := mustCreatePrivateKey(t)
	identity := mustCreatePeerIdentity(t, privKey)
	peerstore := pstoremem.NewPeerstore()

	starTransport := star.New(identity, peerstore).
		WithSignalConfiguration(star.SignalConfiguration{
			URLPath: "/socket.io/?EIO=3&transport=websocket",
		})

	h, err := libp2p.New(ctx,
		libp2p.Identity(privKey),
		libp2p.ListenAddrs(signalMultiaddr),
		libp2p.Peerstore(peerstore),
		libp2p.Transport(starTransport))
	require.NoError(t, err)
	return h
}

func mustCreatePrivateKey(t *testing.T) crypto.PrivKey {
	priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, rand.Reader)
	require.NoError(t, err)
	return priv
}

func mustCreatePeerIdentity(t *testing.T, privKey crypto.PrivKey) peer.ID {
	pid, err := peer.IDFromPublicKey(privKey.GetPublic())
	require.NoError(t, err)
	return pid
}

func mustCreateSignalAddr(t *testing.T, signalAddr string) multiaddr.Multiaddr {
	starMultiaddr, err := multiaddr.NewMultiaddr(signalAddr)
	require.NoError(t, err)
	return starMultiaddr
}


func waitForStream(t *testing.T, newStream func() (network.Stream, error), timeout time.Duration) network.Stream {
	startTime := time.Now()

	var s network.Stream
	var err error

	for time.Now().Before(startTime.Add(timeout)) {
		s, err = newStream()
		if err == nil {
			return s
		}
		time.Sleep(5 * time.Second)
	}

	t.Fatalf("Timeout occurred while waiting for the stream: %v", err)
	return nil
}