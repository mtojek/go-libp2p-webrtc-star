package star

import (
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/mux"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/transport"
	"github.com/multiformats/go-multiaddr"
)

type connection struct{}

var _ transport.CapableConn = new(connection)

func (c *connection) Close() error {
	panic("implement me")
}

func (c *connection) IsClosed() bool {
	panic("implement me")
}

func (c *connection) OpenStream() (mux.MuxedStream, error) {
	panic("implement me")
}

func (c *connection) AcceptStream() (mux.MuxedStream, error) {
	panic("implement me")
}

func (c *connection) LocalPeer() peer.ID {
	panic("implement me")
}

func (c *connection) LocalPrivateKey() crypto.PrivKey {
	panic("implement me")
}

func (c *connection) RemotePeer() peer.ID {
	panic("implement me")
}

func (c *connection) RemotePublicKey() crypto.PubKey {
	panic("implement me")
}

func (c *connection) LocalMultiaddr() multiaddr.Multiaddr {
	panic("implement me")
}

func (connection) RemoteMultiaddr() multiaddr.Multiaddr {
	panic("implement me")
}

func (connection) Transport() transport.Transport {
	panic("implement me")
}
