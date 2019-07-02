package star

import (
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/mux"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/transport"
	"github.com/multiformats/go-multiaddr"
)

type Connection struct {}

var _ transport.CapableConn = new(Connection)

func (c *Connection) Close() error {
	panic("implement me")
}

func (c *Connection) IsClosed() bool {
	panic("implement me")
}

func (c *Connection) OpenStream() (mux.MuxedStream, error) {
	panic("implement me")
}

func (c *Connection) AcceptStream() (mux.MuxedStream, error) {
	panic("implement me")
}

func (c *Connection) LocalPeer() peer.ID {
	panic("implement me")
}

func (c *Connection) LocalPrivateKey() crypto.PrivKey {
	panic("implement me")
}

func (c *Connection) RemotePeer() peer.ID {
	panic("implement me")
}

func (c *Connection) RemotePublicKey() crypto.PubKey {
	panic("implement me")
}

func (c *Connection) LocalMultiaddr() multiaddr.Multiaddr {
	panic("implement me")
}

func (Connection) RemoteMultiaddr() multiaddr.Multiaddr {
	panic("implement me")
}

func (Connection) Transport() transport.Transport {
	panic("implement me")
}