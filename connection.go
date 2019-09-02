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
	"net"
	"sync"
)

type connection struct {
	id             string
	peerConnection *webrtc.PeerConnection
	initChannel    datachannel.ReadWriteCloser
	configuration  connectionConfiguration

	accept chan chan detachResult

	muxedConnection mux.MuxedConn
	m               sync.RWMutex

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

func newConnection(configuration connectionConfiguration, peerConnection *webrtc.PeerConnection,
	initChannel datachannel.ReadWriteCloser) *connection {
	accept := make(chan chan detachResult)
	peerConnection.OnDataChannel(func(dc *webrtc.DataChannel) {
		detachRes := detachChannel(dc)
		accept <- detachRes
	})
	return &connection{
		id:             createRandomID("connection"),
		peerConnection: peerConnection,
		configuration:  configuration,

		accept:      accept,
		initChannel: initChannel,
	}
}

func detachChannel(dc *webrtc.DataChannel) chan detachResult {
	onOpenRes := make(chan detachResult)
	dc.OnOpen(func() {
		// Detach the data channel
		raw, err := dc.Detach()
		onOpenRes <- detachResult{raw, err}
	})
	return onOpenRes
}

func (c *connection) OpenStream() (mux.MuxedStream, error) {
	logger.Debugf("%s: Open stream", c.id)

	muxed, err := c.getMuxed()
	if err != nil {
		return nil, err
	}
	if muxed != nil {
		return muxed.OpenStream()
	}

	rawDC := c.checkInitChannel()
	if rawDC == nil {
		pc, err := c.getPC()
		if err != nil {
			return nil, err
		}
		dc, err := pc.CreateDataChannel("data", nil)
		if err != nil {
			return nil, err
		}

		detachRes := detachChannel(dc)

		res := <-detachRes
		if res.err != nil {
			return nil, res.err
		}
		rawDC = res.dataChannel
	}

	return c.muxedConnection.OpenStream()
}

func (c *connection) checkInitChannel() datachannel.ReadWriteCloser {
	c.m.Lock()
	defer c.m.Unlock()
	if c.initChannel != nil {
		ch := c.initChannel
		c.initChannel = nil
		return ch
	}
	return nil
}

func (c *connection) getPC() (*webrtc.PeerConnection, error) {
	c.m.RLock()
	pc := c.peerConnection
	c.m.RUnlock()

	if pc == nil {
		return nil, errors.New("connection closed")
	}

	return pc, nil
}

func (c *connection) AcceptStream() (mux.MuxedStream, error) {
	logger.Debugf("%s: Accept stream", c.id)
	muxed, err := c.getMuxed()
	if err != nil {
		return nil, err
	}
	if muxed != nil {
		return muxed.AcceptStream()
	}

	rawDC := c.checkInitChannel()
	if rawDC == nil {
		rawDC, err = c.awaitAccept()
	}
	return c.muxedConnection.AcceptStream()
}

func (c *connection) getMuxed() (mux.MuxedConn, error) {
	c.m.Lock()
	defer c.m.Unlock()

	if c.muxedConnection != nil {
		return c.muxedConnection, nil
	}

	rawDC := c.initChannel
	if rawDC == nil {
		var err error
		rawDC, err = c.awaitAccept()
		if err != nil {
			return nil, err
		}
	}

	var err error
	c.muxedConnection, err = c.configuration.multiplexer.NewConn(newStream(rawDC, fakeNetAddress), c.configuration.isServer)
	if err != nil {
		return nil, err
	}

	return c.muxedConnection, nil
}

func (c *connection) awaitAccept() (datachannel.ReadWriteCloser, error) {
	detachRes, ok := <-c.accept
	if !ok {
		return nil, errors.New("connection closed")
	}

	res := <-detachRes
	return res.dataChannel, res.err
}

func (c *connection) IsClosed() bool {
	c.m.RLock()
	pc := c.peerConnection
	c.m.RUnlock()
	return pc == nil
}

func (c *connection) Close() error {
	logger.Debugf("%s: Close connection (no actions)", c.id)
	c.m.Lock()
	defer c.m.Unlock()

	var err error
	if c.peerConnection != nil {
		err = c.peerConnection.Close()
	}
	c.peerConnection = nil
	close(c.accept)
	return err
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
