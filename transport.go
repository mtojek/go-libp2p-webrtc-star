package star

import (
	"context"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/transport"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/pion/webrtc"
	"sync"
)

type Transport struct {
	signals map[string]*signal
	m       sync.RWMutex

	addressBook addressBook
	peerID      peer.ID

	signalConfiguration SignalConfiguration
	webRTCConfiguration webrtc.Configuration
}

var _ transport.Transport = new(Transport)

func (t *Transport) Dial(ctx context.Context, raddr ma.Multiaddr, p peer.ID) (transport.CapableConn, error) {
	logger.Debugf("Dial peer (ID: %s, address: %v)", p, raddr)
	signal, err := t.getOrCreateSignal(raddr)
	if err != nil {
		return nil, err
	}

	peerConnection, err := webrtc.NewPeerConnection(t.webRTCConfiguration)
	if err != nil {
		return nil, err
	}

	offerDescription, err := peerConnection.CreateOffer(nil)
	if err != nil {
		return nil, err
	}

	logger.Debugf("WebRTC offer description: %v", offerDescription.SDP)
	answerDescription, err := signal.handshake(offerDescription)
	if err != nil {
		return nil, err
	}

	logger.Debugf("WebRTC answer description: %v", answerDescription.SDP)
	err = peerConnection.SetRemoteDescription(answerDescription)
	if err != nil {
		return nil, err
	}

	panic("implement me: Dial")
}

func (t *Transport) Listen(laddr ma.Multiaddr) (transport.Listener, error) {
	logger.Debugf("Listen on address: %s", laddr)
	signal, err := t.getOrCreateSignal(laddr)
	if err != nil {
		return nil, err
	}
	return newListener(laddr, signal)
}

func (t *Transport) getOrCreateSignal(addr ma.Multiaddr) (*signal, error) {
	var signal *signal
	var err error
	var ok bool

	sAddr := addr.String()

	t.m.RLock()
	signal, ok = t.signals[sAddr]
	t.m.RUnlock()

	if !ok {
		signal, err = newSignal(addr, t.addressBook, t.peerID, t.signalConfiguration)
		if err != nil {
			return nil, err
		}

		t.m.Lock()
		t.signals[sAddr] = signal
		t.m.Unlock()
	}
	return signal, err
}

func (t *Transport) CanDial(addr ma.Multiaddr) bool {
	return format.Matches(addr)
}

func (t *Transport) Protocols() []int {
	return []int{protocol.Code}
}

func (t *Transport) Proxy() bool {
	return false
}

func New(peerID peer.ID, peerstore addressBook) *Transport {
	return &Transport{
		signals:     map[string]*signal{},
		peerID:      peerID,
		addressBook: peerstore,
	}
}

func (t *Transport) WithSignalConfiguration(c SignalConfiguration) *Transport {
	t.signalConfiguration = c
	return t
}

func (t *Transport) WithWebRTCConfiguration(c webrtc.Configuration) *Transport {
	t.webRTCConfiguration = c
	return t
}
