package star

import (
	"time"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/transport"
	ma "github.com/multiformats/go-multiaddr"
)

type signal struct {
	address ma.Multiaddr
	addressBook addressBook
	configuration SignalConfiguration
}

type SignalConfiguration struct {
	URLPath string
}

type addressBook interface {
	AddAddr(p peer.ID, addr ma.Multiaddr, ttl time.Duration)
}

func newSignal(maddr ma.Multiaddr, addressBook addressBook, configuration SignalConfiguration) *signal {
	return &signal{
		address: maddr.Decapsulate(protocolMultiaddr),
		addressBook: addressBook,
		configuration: configuration,
	}
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