package star

import (
	"errors"
	"fmt"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/transport"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/multiformats/go-multiaddr-net"
	"github.com/pion/webrtc"
	"strings"
	"time"
)

const (
	maxMessageSize = 8192
	messagePrefix  = "42"

	handshakeAnswerTimeout = 10 * time.Second
)

type signal struct {
	peerMultiaddr   ma.Multiaddr
	signalMultiaddr ma.Multiaddr

	acceptedCh      <-chan transport.CapableConn
	handshakeDataCh chan handshakeData
	stopCh          chan<- struct{}

	handshakeSubscription *handshakeSubscription

	webRTCConfiguration webrtc.Configuration
}

type SignalConfiguration struct {
	URLPath string
}

type sessionProperties struct {
	SID                string `json:"sid"`
	PingIntervalMillis int64  `json:"pingInterval"`
	PingTimeoutMillis  int64  `json:"pingTimeout"`
}

func newSignal(signalMultiaddr ma.Multiaddr, addressBook addressBook, peerID peer.ID, signalConfiguration SignalConfiguration,
	webRTCConfiguration webrtc.Configuration) (*signal, error) {
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

	handshakeDataCh := make(chan handshakeData)
	stopCh := make(chan struct{}, 2)

	acceptedCh := startClient(url, peerMultiaddr, smartAddressBook, handshakeDataCh, handshakeSubscription, stopCh)
	return &signal{
		peerMultiaddr:         peerMultiaddr,
		signalMultiaddr:       signalMultiaddr,
		acceptedCh:            acceptedCh,
		handshakeSubscription: handshakeSubscription,
		handshakeDataCh:       handshakeDataCh,
		stopCh:                stopCh,
		webRTCConfiguration:   webRTCConfiguration,
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

func (s *signal) dial(peerID peer.ID) (transport.CapableConn, error) {
	peerConnection, err := webrtc.NewPeerConnection(s.webRTCConfiguration)
	if err != nil {
		return nil, err
	}

	offerDescription, err := peerConnection.CreateOffer(nil)
	if err != nil {
		return nil, err
	}

	answerDescription, err := s.doHandshake(peerID, offerDescription)
	if err != nil {
		return nil, err
	}

	err = peerConnection.SetRemoteDescription(answerDescription)
	if err != nil {
		return nil, err
	}

	panic("implement me")
}

func (s *signal) accept() (transport.CapableConn, error) {
	offerDescription, ok := <-s.handshakeSubscription.unsubscribed()
	if !ok {
		return nil, errors.New("subscription channel has been closed")
	}

	peerConnection, err := webrtc.NewPeerConnection(s.webRTCConfiguration)
	if err != nil {
		return nil, err
	}

	err = peerConnection.SetRemoteDescription(offerDescription.Signal)
	if err != nil {
		return nil, err
	}

	answerDescription, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		return nil, err
	}

	s.answerHandshake(offerDescription.IntentID, offerDescription.SrcMultiaddr, answerDescription)

	panic("implement me")
}

func (s *signal) close() error {
	return s.stopClient()
}

func (s *signal) stopClient() error {
	s.stopCh <- struct{}{}
	close(s.handshakeDataCh)
	return nil
}
