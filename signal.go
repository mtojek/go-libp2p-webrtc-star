package star

import (
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/transport"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/multiformats/go-multiaddr-net"
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
		var connection *websocket.Conn
		var err error

		for {
			time.Sleep(5 * time.Second)

			if stopSignalReceived(stopCh) {
				logger.Debugf("Stop signal received. Closing.")
				return
			}

			if !isConnectionHealthy(connection) {
				logger.Debugf("Connection is not healthy.")

				connection, err = openConnection(url)
				if err != nil {
					logger.Errorf("Can't establish connection: %v", err)
					continue
				}
				logger.Debugf("Connection to signal server established.")
			}

			logger.Debugf("Connection is healthy.")
		}
	}()
	return accepted
}

func stopSignalReceived(stopCh chan struct{}) bool {
	select {
		case <-stopCh:
			return true
		default:
			return false
	}
}

func isConnectionHealthy(connection *websocket.Conn) bool {
	return connection != nil
}

func openConnection(url string) (*websocket.Conn, error) {
	logger.Debugf("Open new connection: %s", url)

	connection, _, err := websocket.DefaultDialer.Dial(url, nil)
	return connection, err
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