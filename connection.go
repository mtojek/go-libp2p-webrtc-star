package star

import (
	"errors"
	"fmt"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/mux"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/transport"
	"github.com/multiformats/go-multiaddr"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/pion/datachannel"
	"github.com/pion/webrtc"
)

type connection struct {
	id             string
	peerConnection *webrtc.PeerConnection
	configuration  connectionConfiguration

	closed bool
}

var _ transport.CapableConn = new(connection)

type connectionConfiguration struct {
	remotePeerID        peer.ID
	remotePeerMultiaddr ma.Multiaddr

	localPeerID        peer.ID
	localPeerMultiaddr ma.Multiaddr

	transport transport.Transport
}

type detachResult struct {
	dataChannel datachannel.ReadWriteCloser
	err         error
}

func newConnection(configuration connectionConfiguration, peerConnection *webrtc.PeerConnection) *connection {
	return &connection{
		id:             createRandomID("connection"),
		peerConnection: peerConnection,
		configuration:  configuration,
	}
}

func (c *connection) OpenStream() (mux.MuxedStream, error) {
	logger.Debugf("%s: Open stream", c.id)
	if c.closed {
		return nil, errors.New("connection already closed")
	}

	dataChannel, err := c.peerConnection.CreateDataChannel("data", nil)
	if err != nil {
		return nil, err
	}

	onOpenResult := make(chan detachResult)
	dataChannel.OnOpen(func() {
		detached, err := dataChannel.Detach()
		onOpenResult <- detachResult{
			dataChannel: detached,
			err:         err,
		}
	})
	r := <-onOpenResult
	if r.err != nil {
		return nil, r.err
	}
	return newStream(r.dataChannel), nil
}

func (c *connection) AcceptStream() (mux.MuxedStream, error) {
	logger.Debugf("%s: Accept stream", c.id)
	return c.foo()
}

func (c *connection) foo() (mux.MuxedStream, error) {
	if c.closed {
		logger.Debug("foo closed")
		return nil, errors.New("connection already closed")
	}

	dataChannel, err := c.peerConnection.CreateDataChannel("data", nil)
	if err != nil {
		logger.Debug("CreateDataChannel err", err)
		return nil, err
	}

	onOpenResult := make(chan detachResult)
	dataChannel.OnOpen(func() {
		detached, err := dataChannel.Detach()
		onOpenResult <- detachResult{
			dataChannel: detached,
			err:         err,
		}
	})
	r := <-onOpenResult
	if r.err != nil {
		logger.Debug("onOpenResult err", err)
		return nil, r.err
	}
	return newStream(r.dataChannel), nil
}

func (c *connection) IsClosed() bool {
	return c.closed
}

func (c *connection) Close() error {
	logger.Debugf("%s: Close connection (no actions)", c.id)
	err := c.peerConnection.Close()
	if err != nil {
		return err
	}
	c.closed = true
	return nil
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
	return fmt.Sprintf("WebRTC connection (ID: %s, localPeerID: %v, localPeerMultiaddr: %v, remotePeerID: %v, remotePeerMultiaddr: %v",
		c.id, c.configuration.localPeerID, c.configuration.localPeerMultiaddr,
		c.configuration.remotePeerID, c.configuration.remotePeerMultiaddr)
}

func (c *connection) LocalPrivateKey() crypto.PrivKey {
	logger.Warningf("%s: Local private key undefined", c.id)
	return nil
}

func (c *connection) RemotePublicKey() crypto.PubKey {
	logger.Warningf("%s: Remote public key undefined", c.id)
	return nil
}
