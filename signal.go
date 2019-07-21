package star

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/transport"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/multiformats/go-multiaddr-net"
)

const (
	maxMessageSize = 2048
)

type signal struct {
	accepted <-chan transport.CapableConn
	stopCh chan<- struct{}
}

type SignalConfiguration struct {
	URLPath string
}

type sessionProperties struct {
	SID string `json:"sid"`
	PingIntervalMillis int64 `json:"pingInterval"`
	PingTimeoutMillis int64 `json:"pingTimeout"`
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

				err := openSession(connection)
				if err != nil {
					logger.Errorf("Can't open session: %v", err)
					connection = nil
					continue
				}
			}

			logger.Debugf("Connection is healthy.")

			message, err := readMessage(connection)
			if err != nil {
				logger.Errorf("Can't read message: %v", err)
				connection = nil
				continue
			}
			logger.Debugf("Message: %s", message)
		}
	}()
	return accepted
}

func openSession(connection *websocket.Conn) error {
	message, err := readMessage(connection)
	if err != nil {
		return err
	}

	var sp sessionProperties
	err = json.Unmarshal(message, &sp)
	if err != nil {
		return err
	}

	pongWait := time.Duration(sp.PingIntervalMillis * int64(time.Millisecond))
	logger.Debugf("Setting read deadline: %v", pongWait)

	connection.SetReadLimit(maxMessageSize)
	//connection.SetReadDeadline(time.Now().Add(pongWait))
	connection.SetPingHandler(func(string) error {
		logger.Debugf("Ping!")
		return nil
	})
	connection.SetPongHandler(func(string) error {
		logger.Debugf("Pong!")
		connection.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	err = readEmptyMessage(connection)
	if err != nil {
		return err
	}
	return nil
}

func readMessage(connection *websocket.Conn) ([]byte, error) {
	_, message, err := connection.ReadMessage()
	if err != nil {
		return nil, err
	}

	i := bytes.IndexByte(message, '{')
	if i < 0 {
		return nil, errors.New("message token not found")
	}
	return message[i:], nil
}

func readEmptyMessage(connection *websocket.Conn) error {
	_, message, err := connection.ReadMessage()
	if err != nil {
		return err
	}

	i := bytes.IndexByte(message, '{')
	if i > 0 {
		return errors.New("empty message expected")
	}
	return nil
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