package star

import (
	"fmt"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/mux"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/transport"
	"github.com/multiformats/go-multiaddr"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/pion/datachannel"
	"github.com/pion/webrtc"
	"net"
	"sync"
)

type connection struct {
	id             string
	peerConnection *webrtc.PeerConnection
	configuration  connectionConfiguration

	muxedConnection mux.MuxedConn
	m               sync.Mutex

	closed bool
}

var _ transport.CapableConn = new(connection)

var fakeNetAddress net.Addr

func init() {
	var err error
	fakeNetAddress, err = net.ResolveTCPAddr("tcp4", "0.0.0.0:1")
	if err != nil {
		logger.Fatalf("can't resolve fake TCP address")
	}
}

type connectionConfiguration struct {
	remotePeerID        peer.ID
	remotePeerMultiaddr ma.Multiaddr

	localPeerID        peer.ID
	localPeerMultiaddr ma.Multiaddr

	transport   transport.Transport
	multiplexer mux.Multiplexer
	isServer    bool
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
	err := c.ensureMuxedConnectionEstablished()
	if err != nil {
		return nil, err
	}
	return c.muxedConnection.OpenStream()
}

func (c *connection) AcceptStream() (mux.MuxedStream, error) {
	logger.Debugf("%s: Accept stream", c.id)
	err := c.ensureMuxedConnectionEstablished()
	if err != nil {
		return nil, err
	}
	return c.muxedConnection.AcceptStream()
}

func (c *connection) ensureMuxedConnectionEstablished() error {
	if c.muxedConnection == nil {
		c.m.Lock()
		if c.muxedConnection == nil {
			var err error
			c.muxedConnection, err = c.createMuxedConnection()
			c.m.Unlock()

			if err != nil {
				return err
			}
		} else {
			c.m.Unlock()
		}
	}
	return nil
}

func (c *connection) createMuxedConnection() (mux.MuxedConn, error) {
	dataChannel, err := c.peerConnection.CreateDataChannel(createRandomID("datachannel"), nil)
	if err != nil {
		logger.Warningf("Can't create data channel: %v", err)
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
		logger.Errorf("Detaching data channel failed: %v", err)
		return nil, r.err
	}
	return c.configuration.multiplexer.NewConn(newStream(r.dataChannel, fakeNetAddress), c.configuration.isServer)
}

func (c *connection) IsClosed() bool {
	return c.closed
}

func (c *connection) Close() error {
	logger.Debugf("%s: Close connection (no actions)", c.id)
	err := c.peerConnection.Close()
	c.closed = true
	if err != nil {
		return err
	}
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
