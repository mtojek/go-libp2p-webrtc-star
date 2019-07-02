package examples

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/mtojek/go-libp2p-webrtc-star"
	"github.com/multiformats/go-multiaddr"
	"github.com/stretchr/testify/require"
)

const (
	protocolID     = "/p2p-webrtc-star/1.0.0"
	starSignalAddr = "/dns4/wrtc-star.discovery.libp2p.io/tcp/443/wss/p2p-webrtc-star"

	peerstoreAddressTTL = 1 * time.Hour
)

var starMultiaddr = mustCreateSignalAddr()

func mustCreateHost(t *testing.T, ctx context.Context) host.Host {
	h, err := libp2p.New(ctx,
		libp2p.Transport(star.New),
		libp2p.ListenAddrs(starMultiaddr))
	require.NoError(t, err)
	return h
}

func mustCreateSignalAddr() multiaddr.Multiaddr {
	starMultiaddr, err := multiaddr.NewMultiaddr(starSignalAddr)
	if err != nil {
		log.Fatal(err)
	}
	return starMultiaddr
}