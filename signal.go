package star

import (
	"time"

	"github.com/libp2p/go-libp2p-core/peer"
	ma "github.com/multiformats/go-multiaddr"
)

type signal struct {
	address ma.Multiaddr
	addressBook addressBook
}

type addressBook interface {
	AddAddr(p peer.ID, addr ma.Multiaddr, ttl time.Duration)
}

func newSignal(maddr ma.Multiaddr, addressBook addressBook) (*signal, error) {
	return &signal{
		address: maddr.Decapsulate(protocolMultiaddr),
		addressBook: addressBook,
	}, nil
}