package star

import (
	"context"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/transport"
	"github.com/multiformats/go-multiaddr"
)

type WebRTCStar struct {}

var _ transport.Transport = new(WebRTCStar)

func (wrs *WebRTCStar) Dial(ctx context.Context, raddr multiaddr.Multiaddr, p peer.ID) (transport.CapableConn, error) {
	panic("implement me: Dial")
}

func (wrs *WebRTCStar) CanDial(addr multiaddr.Multiaddr) bool {
	panic("implement me: CanDial")
}

func (wrs *WebRTCStar) Listen(laddr multiaddr.Multiaddr) (transport.Listener, error) {
	panic("implement me: Listen")
}

func (wrs *WebRTCStar) Protocols() []int {
	return []int{WebRTCStarProtocol.Code}
}

func (wrs *WebRTCStar) Proxy() bool {
	return false
}

