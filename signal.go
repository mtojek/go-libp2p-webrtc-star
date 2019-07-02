package star

import (
	"github.com/libp2p/go-libp2p-core/transport"
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

func (s *signal) Accept() (transport.CapableConn, error) {
	for {
		time.Sleep(1 * time.Minute)
	}

	panic("implement me: Accept")
}

func (s *signal) Close() error {
	panic("implement me: Close")
}