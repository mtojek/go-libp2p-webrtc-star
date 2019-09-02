package star

import (
	"context"
	"errors"
	"fmt"
	"github.com/libp2p/go-libp2p-core/mux"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/transport"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/multiformats/go-multiaddr-net"
	"github.com/pion/datachannel"
	"github.com/pion/webrtc"
	"strings"
	"time"
)

const (
	maxMessageSize = 8192
	messagePrefix  = "42"

	handshakeAnswerTimeout = 10 * time.Minute
)

type signal struct {
	transport transport.Transport

	peerID          peer.ID
	peerMultiaddr   ma.Multiaddr
	signalMultiaddr ma.Multiaddr

	acceptedCh      <-chan transport.CapableConn
	handshakeDataCh chan<- handshakeData

	handshakeSubscription *handshakeSubscription
	webRTCConfiguration   webrtc.Configuration
	multiplexer           mux.Multiplexer

	stopCh chan<- struct{}
}

type SignalConfiguration struct {
	URLPath string
}

type sessionProperties struct {
	SID                string `json:"sid"`
	PingIntervalMillis int64  `json:"pingInterval"`
	PingTimeoutMillis  int64  `json:"pingTimeout"`
}

var webrtcapi *webrtc.API

func init() {
	settingEngine := webrtc.SettingEngine{}
	settingEngine.DetachDataChannels()
	webrtcapi = webrtc.NewAPI(webrtc.WithSettingEngine(settingEngine))
}

func newSignal(transport transport.Transport, signalMultiaddr ma.Multiaddr, addressBook addressBook, peerID peer.ID,
	signalConfiguration SignalConfiguration, webRTCConfiguration webrtc.Configuration, multiplexer mux.Multiplexer) (*signal, error) {
	url, err := createSignalURL(signalMultiaddr, signalConfiguration)
	if err != nil {
		return nil, err
	}

	peerMultiaddr, err := createPeerMultiaddr(signalMultiaddr, peerID)
	if err != nil {
		return nil, err
	}

	smartAddressBook := decorateSelfIgnoreAddressBook(addressBook, peerID)
	handshakeSubscription := newHandshakeSubscription()

	stopCh := make(chan struct{}, 2)

	acceptedCh, handshakeDataCh := startClient(url, peerMultiaddr, smartAddressBook, handshakeSubscription, stopCh)
	return &signal{
		transport:             transport,
		peerID:                peerID,
		peerMultiaddr:         peerMultiaddr,
		signalMultiaddr:       signalMultiaddr,
		acceptedCh:            acceptedCh,
		handshakeSubscription: handshakeSubscription,
		handshakeDataCh:       handshakeDataCh,
		stopCh:                stopCh,
		webRTCConfiguration:   webRTCConfiguration,
		multiplexer:           multiplexer,
	}, nil
}

func createSignalURL(addr ma.Multiaddr, configuration SignalConfiguration) (string, error) {
	websocketAddr := addr.Decapsulate(protocolMultiaddr)

	var buf strings.Builder
	buf.WriteString(readProtocolForSignalURL(websocketAddr))

	_, hostPort, err := manet.DialArgs(websocketAddr)
	if err != nil {
		return "", err
	}
	buf.WriteString(hostPort)
	buf.WriteString(configuration.URLPath)
	return buf.String(), nil
}

func createPeerMultiaddr(signalMultiaddr ma.Multiaddr, peerID peer.ID) (ma.Multiaddr, error) {
	ipfsMultiaddr, err := ma.NewMultiaddr(fmt.Sprintf("/%s/%s", ipfsProtocolName, peerID.String()))
	if err != nil {
		logger.Fatal(err)
	}
	return signalMultiaddr.Encapsulate(ipfsMultiaddr), nil
}

func readProtocolForSignalURL(maddr ma.Multiaddr) string {
	if _, err := maddr.ValueForProtocol(wssProtocolCode); err == nil {
		return "wss://"
	}
	return "ws://"
}

func (s *signal) dial(ctx context.Context, remotePeerID peer.ID) (transport.CapableConn, error) {
	peerConnection, err := webrtcapi.NewPeerConnection(s.webRTCConfiguration)
	if err != nil {
		return nil, err
	}

	offerDescription, err := peerConnection.CreateOffer(nil)
	if err != nil {
		return nil, err
	}

	err = peerConnection.SetLocalDescription(offerDescription)
	if err != nil {
		return nil, err
	}

	dstMultiaddr, err := ma.NewMultiaddr(fmt.Sprintf("/%s/%s", ipfsProtocolName, remotePeerID.String()))
	if err != nil {
		return nil, err
	}
	offer := handshakeData{
		IntentID:     createRandomIntentID(),
		SrcMultiaddr: s.peerMultiaddr.String(),
		DstMultiaddr: s.signalMultiaddr.Encapsulate(dstMultiaddr).String(),
		Signal:       offerDescription,
	}
	answer, err := s.doHandshake(ctx, offer)
	if err != nil {
		return nil, err
	}

	err = peerConnection.SetRemoteDescription(answer.Signal)
	if err != nil {
		return nil, err
	}

	return s.openConnection(offer.DstMultiaddr, peerConnection, false)
}

func (s *signal) accept() (transport.CapableConn, error) {
	offer, ok := <-s.handshakeSubscription.unsubscribed()
	if !ok {
		return nil, errors.New("subscription channel has been closed")
	}

	peerConnection, err := webrtcapi.NewPeerConnection(s.webRTCConfiguration)
	if err != nil {
		return nil, err
	}

	err = peerConnection.SetRemoteDescription(offer.Signal)
	if err != nil {
		return nil, err
	}

	answerDescription, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		return nil, err
	}

	err = peerConnection.SetLocalDescription(answerDescription)
	if err != nil {
		return nil, err
	}

	answer := handshakeData{
		IntentID:     offer.IntentID,
		SrcMultiaddr: offer.SrcMultiaddr,
		DstMultiaddr: s.peerMultiaddr.String(),
		Signal:       answerDescription,
		Answer:       true,
	}
	s.answerHandshake(answer)
	return s.openConnection(offer.SrcMultiaddr, peerConnection, true)
}

func (s *signal) openConnection(destination string, peerConnection *webrtc.PeerConnection, isServer bool) (transport.CapableConn, error) {
	dstMultiaddr, err := ma.NewMultiaddr(destination)
	if err != nil {
		return nil, err
	}

	value, err := dstMultiaddr.ValueForProtocol(ma.P_P2P)
	if err != nil {
		return nil, err
	}

	remotePeerID, err := peer.IDB58Decode(value)
	if err != nil {
		return nil, err
	}

	var detachedDataChannel datachannel.ReadWriteCloser
	if !isServer {
		dc, err := peerConnection.CreateDataChannel("data", nil)
		if err != nil {
			return nil, err
		}

		detachRes := detachChannel(dc)
		select {
		case res := <-detachRes:
			if res.err != nil {
				return nil, res.err
			}
			detachedDataChannel = res.dataChannel
		}
		// TODO ctx
	}

	return newConnection(connectionConfiguration{
		remotePeerID:        remotePeerID,
		remotePeerMultiaddr: dstMultiaddr,

		localPeerID:        s.peerID,
		localPeerMultiaddr: s.peerMultiaddr,

		transport:   s.transport,
		multiplexer: s.multiplexer,
		isServer:    isServer,
	}, peerConnection, detachedDataChannel), nil
}

func (s *signal) close() error {
	return s.stopClient()
}

func (s *signal) stopClient() error {
	s.stopCh <- struct{}{}
	return nil
}

func createRandomIntentID() string {
	return createRandomID("signal")
}
