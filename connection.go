package star

import (
	"fmt"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/mux"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/transport"
	"github.com/multiformats/go-multiaddr"
	ma "github.com/multiformats/go-multiaddr"
	"time"
)

type connection struct {
	id string

	configuration connectionConfiguration
}

type connectionConfiguration struct {
	remotePeerID        peer.ID
	remotePeerMultiaddr ma.Multiaddr

	localPeerID        peer.ID
	localPeerMultiaddr ma.Multiaddr

	transport transport.Transport
}

var _ transport.CapableConn = new(connection)

func newConnection(configuration connectionConfiguration) *connection {
	return &connection{
		id:            createRandomID("connection"),
		configuration: configuration,
	}
}

func (c *connection) Close() error {
	logger.Debugf("%s: Close stream", c.id)
	return nil
	//time.Sleep(20 * time.Minute)
	//panic("implement me")
}

func (c *connection) IsClosed() bool {
	panic("implement me")
}

func (c *connection) OpenStream() (mux.MuxedStream, error) {
	logger.Debugf("%s: Open stream", c.id)

	time.Sleep(20 * time.Minute)
	panic("implement me")
}

func (c *connection) AcceptStream() (mux.MuxedStream, error) {
	logger.Debugf("%s: Accept stream", c.id)

	time.Sleep(20 * time.Minute)
	panic("implement me")
}

func (c *connection) LocalPeer() peer.ID {
	return c.configuration.localPeerID
}

func (c *connection) RemotePeer() peer.ID {
	return c.configuration.remotePeerID
}

func (c *connection) LocalMultiaddr() multiaddr.Multiaddr {
	return c.configuration.localPeerMultiaddr
}

func (c *connection) RemoteMultiaddr() multiaddr.Multiaddr {
	return c.configuration.remotePeerMultiaddr
}

func (c *connection) Transport() transport.Transport {
	return c.configuration.transport
}

func (c *connection) String() string {
	return fmt.Sprintf("WebRTC connection (localPeerID: %v, localPeerMultiaddr: %v, remotePeerID: %v, remotePeerMultiaddr: %v",
		c.configuration.localPeerID, c.configuration.localPeerMultiaddr,
		c.configuration.remotePeerID, c.configuration.remotePeerMultiaddr)
}

func (c *connection) LocalPrivateKey() crypto.PrivKey {
	return nil // not supported yet
}

func (c *connection) RemotePublicKey() crypto.PubKey {
	return nil // not supported yet
}
