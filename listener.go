package star

import (
	"net"

	"github.com/libp2p/go-libp2p-core/transport"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/multiformats/go-multiaddr-net"
)

type listener struct {
	address ma.Multiaddr
	signal *signal
}

var _ transport.Listener = new(listener)

func newListener(address ma.Multiaddr, addressBook addressBook) (*listener, error) {
	signal, err := newSignal(address, addressBook)
	if err != nil {
		return nil, err
	}
	return &listener{
		address: address,
		signal: signal,
	}, nil
}

func (l *listener) Accept() (transport.CapableConn, error) {
	panic("implement me: Accept")
}

func (l *listener) Close() error {
	panic("implement me: Close")
}

func (l *listener) Addr() net.Addr {
	networkAddress, err := manet.ToNetAddr(l.address)
	if err != nil {
		logger.Fatal(err)
	}
	return networkAddress
}

func (l *listener) Multiaddr() ma.Multiaddr {
	return l.address
}