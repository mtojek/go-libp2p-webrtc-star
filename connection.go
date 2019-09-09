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
	"github.com/pion/webrtc/v2"
	"net"
	"sync"
)

type connection struct {
	id             string
	peerConnection *webrtc.PeerConnection
	initChannel    datachannel.ReadWriteCloser
	configuration  connectionConfiguration

	dataChannelDetachedCh chan chan detachResult
	m                     sync.RWMutex
	muxedConnection       mux.MuxedConn
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

	isServer bool
}

type detachResult struct {
	dataChannel datachannel.ReadWriteCloser
	err         error
}

func newConnection(configuration connectionConfiguration, peerConnection *webrtc.PeerConnection,
	initChannel datachannel.ReadWriteCloser) *connection {
	dataChannelDetachedCh := make(chan chan detachResult)
	peerConnection.OnDataChannel(func(dc *webrtc.DataChannel) {
		dataChannelDetachedCh <- detachDataChannel(dc)
	})
	return &connection{
		id:             createRandomID("connection"),
		peerConnection: peerConnection,
		configuration:  configuration,

		dataChannelDetachedCh: dataChannelDetachedCh,
		initChannel:           initChannel,
	}
}

func detachDataChannel(dataChannel *webrtc.DataChannel) chan detachResult {
	detachedCh := make(chan detachResult)
	dataChannel.OnOpen(func() {
		channel, err := dataChannel.Detach()
		detachedCh <- detachResult{channel, err}
	})
	return detachedCh
}

func (c *connection) OpenStream() (mux.MuxedStream, error) {
	logger.Debugf("%s: Open stream", c.id)

	muxedConnection, err := c.getMuxedConnection()
	if err != nil {
		return nil, err
	}
	if muxedConnection != nil {
		return muxedConnection.OpenStream()
	}

	rawDataChannel := c.checkInitChannel()
	if rawDataChannel == nil {
		pc, err := c.getPeerConnection()
		if err != nil {
			return nil, err
		}
		dc, err := pc.CreateDataChannel("data", nil)
		if err != nil {
			return nil, err
		}

		detached := <-detachDataChannel(dc)
		if detached.err != nil {
			return nil, detached.err
		}
		rawDataChannel = detached.dataChannel
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

func (c *connection) getPeerConnection() (*webrtc.PeerConnection, error) {
	c.m.RLock()
	pc := c.peerConnection
	c.m.RUnlock()

	if pc == nil {
		return nil, errors.New("peer connection closed")
	}
	return pc, nil
}

func (c *connection) AcceptStream() (mux.MuxedStream, error) {
	logger.Debugf("%s: Accept stream", c.id)
	muxedConnection, err := c.getMuxedConnection()
	if err != nil {
		return nil, err
	}
	if muxedConnection != nil {
		return muxedConnection.AcceptStream()
	}

	rawDataChannel := c.checkInitChannel()
	if rawDataChannel == nil {
		rawDataChannel, err = c.awaitDataChannelDetached()
	}
	return c.muxedConnection.AcceptStream()
}

func (c *connection) getMuxedConnection() (mux.MuxedConn, error) {
	c.m.Lock()
	defer c.m.Unlock()

	if c.muxedConnection != nil {
		return c.muxedConnection, nil
	}

	rawDataChannel := c.initChannel
	if rawDataChannel == nil {
		var err error
		rawDataChannel, err = c.awaitDataChannelDetached()
		if err != nil {
			return nil, err
		}
	}

	var err error
	c.muxedConnection, err = c.configuration.multiplexer.NewConn(newStream(rawDataChannel, fakeNetAddress),
		c.configuration.isServer)
	if err != nil {
		return nil, err
	}
	return c.muxedConnection, nil
}

func (c *connection) awaitDataChannelDetached() (datachannel.ReadWriteCloser, error) {
	detachedCh, ok := <-c.dataChannelDetachedCh
	if !ok {
		return nil, errors.New("connection closed")
	}

	detached := <-detachedCh
	return detached.dataChannel, detached.err
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
	close(c.dataChannelDetachedCh)
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
