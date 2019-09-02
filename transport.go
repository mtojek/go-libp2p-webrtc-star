package star

import (
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p-core/mux"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/transport"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/pion/webrtc"
	"sync"
)

type Transport struct {
	signals map[string]*signal
	m       sync.Mutex

	addressBook addressBook
	peerID      peer.ID

	signalConfiguration SignalConfiguration
	webRTCConfiguration webrtc.Configuration
	multiplexer         mux.Multiplexer
}

var _ transport.Transport = new(Transport)

func (t *Transport) Dial(ctx context.Context, raddr ma.Multiaddr, p peer.ID) (transport.CapableConn, error) {
	logger.Debugf("Dial peer (ID: %s, address: %v)", p, raddr)
	signal, err := t.getOrRegisterSignal(raddr)
	if err != nil {
		return nil, err
	}
	return signal.dial(ctx, p)
}

func (t *Transport) Listen(laddr ma.Multiaddr) (transport.Listener, error) {
	logger.Debugf("Listen on address: %s", laddr)
	signal, err := t.getOrRegisterSignal(laddr)
	if err != nil {
		return nil, err
	}
	return newListener(laddr, signal, t.unregisterSignal)
}

func (t *Transport) getOrRegisterSignal(addr ma.Multiaddr) (*signal, error) {
	var err error

	sAddr := addr.String()

	t.m.Lock()
	defer t.m.Unlock()

	if signal, ok := t.signals[sAddr]; ok {
		return signal, nil
	}

	t.signals[sAddr], err = newSignal(t, addr, t.addressBook, t.peerID, t.signalConfiguration, t.webRTCConfiguration,
		t.multiplexer)
	if err != nil {
		return nil, err
	}
	return t.signals[sAddr], nil
}

func (t *Transport) unregisterSignal(addr ma.Multiaddr) error {
	sAddr := addr.String()

	t.m.Lock()
	defer t.m.Unlock()

	if signal, ok := t.signals[sAddr]; ok {
		err := signal.close()
		if err != nil {
			logger.Errorf("Error while closing signal: %v", err)
		}
		delete(t.signals, sAddr)
		return nil
	}
	return fmt.Errorf(`no signal registered for "%s"`, sAddr)
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

func New(peerID peer.ID, peerstore addressBook, multiplexer mux.Multiplexer) *Transport {
	return &Transport{
		signals:     map[string]*signal{},
		peerID:      peerID,
		addressBook: peerstore,
		multiplexer: multiplexer,
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
