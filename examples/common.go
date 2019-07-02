package examples

import (
	"context"
	"github.com/libp2p/go-libp2p-core/network"
	"log"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-peerstore/pstoremem"
	"github.com/mtojek/go-libp2p-webrtc-star"
	"github.com/multiformats/go-multiaddr"
	"github.com/stretchr/testify/require"
)

const (
	protocolID     = "/p2p-webrtc-star/1.0.0"
	firstSignalAddr = "/dns4/wrtc-star.discovery.libp2p.io/tcp/443/wss/p2p-webrtc-star"
	secondSignalAddr = "/dns4/star-signal.cloud.ipfs.team/tcp/443/wss/p2p-webrtc-star"

	waitForStreamTimeout = 60 * time.Minute
)

func mustCreateHost(t *testing.T, ctx context.Context) host.Host {
	firstSignalMultiaddr := mustCreateSignalAddr(firstSignalAddr)
	secondSignalMultiaddr := mustCreateSignalAddr(secondSignalAddr)

	peerstore := pstoremem.NewPeerstore()

	h, err := libp2p.New(ctx,
		libp2p.ListenAddrs(firstSignalMultiaddr, secondSignalMultiaddr),
		libp2p.Peerstore(peerstore),
		libp2p.Transport(star.New(peerstore)))
	require.NoError(t, err)
	return h
}

func mustCreateSignalAddr(signalAddr string) multiaddr.Multiaddr {
	starMultiaddr, err := multiaddr.NewMultiaddr(signalAddr)
	if err != nil {
		log.Fatal(err)
	}
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