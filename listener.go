package star

import (
	"github.com/libp2p/go-libp2p-core/transport"
	"github.com/multiformats/go-multiaddr"
	"github.com/multiformats/go-multiaddr-net"
	"net"
)

type listener struct {
	address multiaddr.Multiaddr
}

var _ transport.Listener = new(listener)

func newListener(address multiaddr.Multiaddr) *listener {
	return &listener{
		address: address,
	}
}

func (l *listener) Accept() (transport.CapableConn, error) {
	panic("implement me: Accept")
}

func (l *listener) Close() error {
	return nil // TODO
}

func (l *listener) Addr() net.Addr {
	nAddr, err := manet.ToNetAddr(l.address)
	if err != nil {
		logger.Fatal(err)
	}
	return nAddr
}

func (l *listener) Multiaddr() multiaddr.Multiaddr {
	return l.address
}