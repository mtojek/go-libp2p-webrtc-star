package star

import (
	"context"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/transport"
	"github.com/multiformats/go-multiaddr"
)

type Transport struct {
	addressBook addressBook
	signalConfiguration SignalConfiguration
}

var _ transport.Transport = new(Transport)

func (t *Transport) Dial(ctx context.Context, raddr multiaddr.Multiaddr, p peer.ID) (transport.CapableConn, error) {
	panic("implement me: Dial")
}

func (t *Transport) CanDial(addr multiaddr.Multiaddr) bool {
	return format.Matches(addr)
}

func (t *Transport) Listen(laddr multiaddr.Multiaddr) (transport.Listener, error) {
	logger.Debugf("Listen on address: %s", laddr)
	return newListener(laddr, t.addressBook, t.signalConfiguration), nil
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

func (t *Transport) WithPeerstore(a addressBook) *Transport {
	t.addressBook = a
	return t
}

func (t *Transport) WithSignalConfiguration(c SignalConfiguration) *Transport {
	t.signalConfiguration = c
	return t
}