package examples

import (
	"context"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-tcp-transport"
	"github.com/mtojek/go-libp2p-webrtc-star"
	wss "github.com/mtojek/go-wss-transport"
	"github.com/multiformats/go-multiaddr"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
	"time"
)

const (
	protocolID     = "/p2p-webrtc-star/1.0.0"
	starSignalAddr = "/dns4/wrtc-star.discovery.libp2p.io/tcp/443/wss/p2p-webrtc-star"

	peerstoreAddressTTL = 1 * time.Hour
)

func mustCreateHost(t *testing.T, ctx context.Context) host.Host {
	h, err := libp2p.New(ctx,
		libp2p.Transport(tcp.NewTCPTransport),
		libp2p.Transport(wss.New),
		libp2p.Transport(star.New),
		libp2p.DefaultMuxers,
		libp2p.DefaultSecurity)
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
