package star

import (
	"encoding/json"
	"fmt"
	"github.com/pion/webrtc"
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
	messagePrefix  = "42"

	ssJoinMessageType = "ss-join"
)

type signal struct {
	acceptedCh           <-chan transport.CapableConn
	outgoingHandshakesCh chan outgoingHandshake
	stopCh               chan<- struct{}

	webRTCConfiguration webrtc.Configuration
}

type outgoingHandshake struct {
	destinationPeerID peer.ID
	offer             webrtc.SessionDescription
	answerCh          chan<- handshakeAnswer
}

type handshakeAnswer struct {
	sessionDescription webrtc.SessionDescription
	err                error
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
	outgoingHandshakesCh := make(chan outgoingHandshake)
	stopCh := make(chan struct{})
	acceptedCh := startClient(url, signalMultiaddr, peerMultiaddr, smartAddressBook, outgoingHandshakesCh, stopCh)
	return &signal{
		acceptedCh:           acceptedCh,
		outgoingHandshakesCh: outgoingHandshakesCh,
		stopCh:               stopCh,
		webRTCConfiguration:  webRTCConfiguration,
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

func startClient(url string, signalMultiaddr ma.Multiaddr, peerMultiaddr ma.Multiaddr, addressBook addressBook,
	outgoingHandshakesCh <-chan outgoingHandshake, stopCh <-chan struct{}) <-chan transport.CapableConn {
	logger.Debugf("Use signal server: %s", url)

	acceptedCh := make(chan transport.CapableConn)
	go func() {
		var connection *websocket.Conn
		var sp *sessionProperties
		var err error

		for {
			if stopSignalReceived(stopCh) {
				logger.Debugf("Stop signal received. Closing")
				return
			}

			if !isConnectionHealthy(connection) {
				connection, err = openConnection(url)
				if err != nil {
					logger.Errorf("Can't establish connection: %v", err)
					time.Sleep(3 * time.Second)
					continue
				}
				logger.Debugf("Connection to signal server established")

				sp, err = openSession(connection, peerMultiaddr)
				if err != nil {
					logger.Errorf("Can't open session: %v", err)
					connection = nil
					continue
				}
			}

			logger.Debugf("%s: Connection is healthy.", sp.SID)

			if handshake, ok := nextOutgoingHandshake(outgoingHandshakesCh); ok {
				fmt.Println(handshake, ok)
			}

			message, err := readMessage(connection)
			if err != nil {
				logger.Errorf("%s: Can't read message: %v", sp.SID, err)
				connection = nil
				continue
			}
			logger.Debugf("%s: Received message: %s", sp.SID, message)
			err = processMessage(addressBook, message)
			if err != nil {
				logger.Errorf("%s: Can't process message: %v", sp.SID, err)
				continue
			}
		}
	}()
	return acceptedCh
}

func openSession(connection *websocket.Conn, peerMultiaddr ma.Multiaddr) (*sessionProperties, error) {
	message, err := readMessage(connection)
	if err != nil {
		return nil, err
	}

	var sp sessionProperties
	err = json.Unmarshal(message, &sp)
	if err != nil {
		return nil, err
	}

	pingInterval := time.Duration(sp.PingIntervalMillis * int64(time.Millisecond))
	pingTimeout := time.Duration(sp.PingTimeoutMillis * int64(time.Millisecond))
	logger.Debugf("%s: Ping interval: %v, Ping timeout: %v", sp.SID, pingInterval, pingTimeout)

	connection.SetReadLimit(maxMessageSize)
	connection.SetPongHandler(func(string) error {
		logger.Debugf("%s: Pong message received", sp.SID)
		connection.SetReadDeadline(time.Time{})
		return nil
	})

	err = readEmptyMessage(connection)
	if err != nil {
		return nil, err
	}

	go func() {
		pingTicker := time.NewTicker(pingInterval)
		for range pingTicker.C {
			logger.Debugf("%s: Send ping message", sp.SID)
			connection.SetReadDeadline(time.Now().Add(pingTimeout))
			err := sendMessage(connection, "ping", nil) // Application layer ping?
			if err != nil {
				logger.Errorf("%s: Can't send ping message: %v", sp.SID, err)
				pingTicker.Stop()
				return
			}

			err = connection.WriteControl(websocket.PingMessage, []byte("ping"), time.Time{})
			if err != nil {
				logger.Errorf("%s: Can't send ping message: %v", sp.SID, err)
				pingTicker.Stop()
				return
			}
		}
	}()

	logger.Debugf("%s: Join peer network (peerID: %s)", sp.SID, peerMultiaddr.String())
	err = sendMessage(connection, ssJoinMessageType, peerMultiaddr.String())
	if err != nil {
		return nil, err
	}
	return &sp, nil
}

func stopSignalReceived(stopCh <-chan struct{}) bool {
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

func nextOutgoingHandshake(outgoingHandshakeCh <-chan outgoingHandshake) (outgoingHandshake, bool) {
	fmt.Println(outgoingHandshakeCh)
	return outgoingHandshake{}, false
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

	answerCh := make(chan handshakeAnswer)

	logger.Debugf("WebRTC offer description: %v", offerDescription.SDP)
	s.outgoingHandshakesCh <- outgoingHandshake{
		destinationPeerID: peerID,
		offer:             offerDescription,
		answerCh:          answerCh,
	}
	answer := <-answerCh
	if answer.err != nil {
		return nil, err
	}

	logger.Debugf("WebRTC answer description: %v", answer.sessionDescription.SDP)
	err = peerConnection.SetRemoteDescription(answer.sessionDescription)
	if err != nil {
		return nil, err
	}

	panic("implement me")
}

func (s *signal) accept() (transport.CapableConn, error) {
	return <-s.acceptedCh, nil
}

func (s *signal) close() error {
	return s.stopClient()
}

func (s *signal) stopClient() error {
	s.stopCh <- struct{}{}
	return nil
}
