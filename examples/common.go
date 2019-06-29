package examples

import (
	"github.com/multiformats/go-multiaddr"
	"log"
)

const (
	protocolID = "/p2p-webrtc-star/1.0.0"
	starSignalAddr = "/dns4/wrtc-star.discovery.libp2p.io/tcp/443/wss/p2p-webrtc-star"
)

var starMultiaddr multiaddr.Multiaddr

func init() {
	starMultiaddr = mustCreateSignalAddr()
}

func mustCreateSignalAddr() multiaddr.Multiaddr {
	var err error
	starMultiaddr, err = multiaddr.NewMultiaddr(starSignalAddr)
	if err != nil {
		log.Fatal(err)
	}
	return starMultiaddr
}