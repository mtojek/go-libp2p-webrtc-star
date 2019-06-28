package star

import (
	"context"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/transport"
	"github.com/multiformats/go-multiaddr"
)

type WebRTCStar struct {}

var _ transport.Transport = new(WebRTCStar)

func (t *WebRTCStar) Dial(ctx context.Context, raddr multiaddr.Multiaddr, p peer.ID) (transport.CapableConn, error) {
	panic("implement me: Dial")
}

func (t *WebRTCStar) CanDial(addr multiaddr.Multiaddr) bool {
	panic("implement me: CanDial")
}

func (t *WebRTCStar) Listen(laddr multiaddr.Multiaddr) (transport.Listener, error) {
	panic("implement me: Listen")
}

func (t *WebRTCStar) Protocols() []int {
	return []int{protocol.Code}
}

func (t *WebRTCStar) Proxy() bool {
	return false
}