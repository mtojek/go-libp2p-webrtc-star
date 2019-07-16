package star

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/multiformats/go-multiaddr-net"
	"strings"
	"time"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/transport"
	ma "github.com/multiformats/go-multiaddr"
)

type signal struct {
	accepted <-chan transport.CapableConn
	stopCh chan<- struct{}
}

type SignalConfiguration struct {
	URLPath string
}

type addressBook interface {
	AddAddr(p peer.ID, addr ma.Multiaddr, ttl time.Duration)
}

func newSignal(maddr ma.Multiaddr, addressBook addressBook, configuration SignalConfiguration) (*signal, error) {
	url, err := createSignalURL(maddr.Decapsulate(protocolMultiaddr), configuration)
	if err != nil {
		return nil, err
	}
	stopCh := make(chan struct{})
	accepted := startClient(url, addressBook, stopCh)
	return &signal{
		accepted: accepted,
		stopCh: stopCh,
	}, nil
}

func createSignalURL(addr ma.Multiaddr, configuration SignalConfiguration) (string, error) {
	var buf strings.Builder
	buf.WriteString(readProtocolForSignalURL(addr))

	_, hostPort, err := manet.DialArgs(addr)
	if err != nil {
		return "", err
	}
	buf.WriteString(hostPort)
	buf.WriteString(configuration.URLPath)
	return buf.String(), nil
}

func readProtocolForSignalURL(maddr ma.Multiaddr) string {
	if _, err := maddr.ValueForProtocol(wssProtocolCode); err == nil {
		return "wss://"
	}
	return "ws://"
}

func startClient(url string, addressBook addressBook, stopCh chan struct{}) <-chan transport.CapableConn {
	logger.Debugf("Use signal server: %s", url)

	accepted := make(chan transport.CapableConn)
	go func() {
		for {

		}
	}()
	return accepted
}

func (s *signal) ensureConnectionEstablished() error {
	var err error

	if s.connection != nil {
		return nil
	}

	s.connection, _, err = websocket.DefaultDialer.Dial(s.url, nil)
	if err != nil {
		return err
	}
	return nil
}

func (s *signal) Accept() (transport.CapableConn, error) {
	return <- s.accepted, nil
}

func (s *signal) Close() error {
	return s.stopClient()
}

func (s *signal) stopClient() error {
	s.stopCh <- struct{}{}
	return nil
}