package star

import (
	"github.com/gorilla/websocket"
	"github.com/multiformats/go-multiaddr-net"
	"strings"
	"time"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/transport"
	ma "github.com/multiformats/go-multiaddr"
)

type signal struct {
	addressBook addressBook
	url string

	connection *websocket.Conn
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
	logger.Debugf("Use signal server: %s", url)
	return &signal{
		addressBook: addressBook,
		url: url,
	}, nil
}

func (s *signal) Accept() (transport.CapableConn, error) {
	err := s.ensureConnectionEstablished()
	if err != nil {
		return nil, err
	}

	panic("implement me: Accept")
}

func (s *signal) ensureConnectionEstablished() error {
	var err error
	s.connection, _, err = websocket.DefaultDialer.Dial(s.url, nil)
	if err != nil {
		return err
	}
	return nil
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

func (s *signal) Close() error {
	if s.connection != nil {
		return s.connection.Close()
	}
	return nil // TODO close other connections
}