package examples

import (
	"context"
	"log"
	"testing"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-peerstore/pstoremem"
	"github.com/mtojek/go-libp2p-webrtc-star"
	"github.com/multiformats/go-multiaddr"
	"github.com/stretchr/testify/require"
)

const (
	protocolID     = "/p2p-webrtc-star/1.0.0"
	signalAddr = "/dns4/wrtc-star.discovery.libp2p.io/tcp/443/wss/p2p-webrtc-star"
)

func mustCreateHost(t *testing.T, ctx context.Context) host.Host {
	signalMultiaddr := mustCreateSignalAddr()
	peerstore := pstoremem.NewPeerstore()

	h, err := libp2p.New(ctx,
		libp2p.ListenAddrs(signalMultiaddr),
		libp2p.Peerstore(peerstore),
		libp2p.Transport(star.New(peerstore)))
	require.NoError(t, err)
	return h
}

func mustCreateSignalAddr() multiaddr.Multiaddr {
	starMultiaddr, err := multiaddr.NewMultiaddr(signalAddr)
	if err != nil {
		log.Fatal(err)
	}
	return starMultiaddr
}
