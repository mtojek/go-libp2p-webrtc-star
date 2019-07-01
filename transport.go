package star

import (
	"context"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/transport"
	"github.com/multiformats/go-multiaddr"
)

type Transport struct {}

var _ transport.Transport = new(Transport)

func (t *Transport) Dial(ctx context.Context, raddr multiaddr.Multiaddr, p peer.ID) (transport.CapableConn, error) {
	panic("implement me: Dial")
}

func (t *Transport) CanDial(addr multiaddr.Multiaddr) bool {
	return format.Matches(addr)
}

func (t *Transport) Listen(laddr multiaddr.Multiaddr) (transport.Listener, error) {
	return newListener(laddr)
}

func (t *Transport) Protocols() []int {
	return []int{protocol.Code}
}

func (t *Transport) Proxy() bool {
	return false
}

func New() *Transport {
	return new(Transport)
}