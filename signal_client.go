package star

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/libp2p/go-libp2p-core/transport"
	ma "github.com/multiformats/go-multiaddr"
	"time"
)

func startClient(url string, peerMultiaddr ma.Multiaddr, addressBook addressBook,
	handshakeDataCh <-chan handshakeData, handshakeSubscription *handshakeSubscription,
	stopCh <-chan struct{}) <-chan transport.CapableConn {
	logger.Debugf("Use signal server: %s", url)

	acceptedCh := make(chan transport.CapableConn)
	internalStopCh := make(chan struct{})
	threadsRunning := false

	stopSessionThreads := func() {
		logger.Debugf("Stop active session threads")

		internalStopCh <- struct{}{}
		internalStopCh <- struct{}{}
	}

	go func() {
		var connection *websocket.Conn
		var sp *sessionProperties
		var err error

		for {
			if stopSignalReceived(stopCh) {
				logger.Debugf("Stop signal received. Closing")
				stopSessionThreads()
				return
			}

			if !isConnectionHealthy(connection) {
				if threadsRunning {
					stopSessionThreads()
					threadsRunning = false
				}

				connection, err = openConnection(url)
				if err != nil {
					logger.Errorf("Can't establish connection: %v", err)
					time.Sleep(3 * time.Second)
					continue
				}
				logger.Debugf("Connection to signal server established")

				sp, err = openSession(connection, peerMultiaddr, handshakeDataCh, internalStopCh)
				if err != nil {
					logger.Errorf("Can't open session: %v", err)
					connection = nil
					continue
				}
				threadsRunning = true
			}

			logger.Debugf("%s: Connection is healthy.", sp.SID)

			message, err := readMessage(connection)
			if err != nil {
				logger.Errorf("%s: Can't read message: %v", sp.SID, err)
				connection = nil
				continue
			}
			logger.Debugf("%s: Received message: %s", sp.SID, message)
			err = processMessage(addressBook, handshakeSubscription, message)
			if err != nil {
				logger.Warningf("%s: Can't process message: %v", sp.SID, err)
				continue
			}
		}
	}()
	return acceptedCh
}

func openSession(connection *websocket.Conn, peerMultiaddr ma.Multiaddr,
	handshakeDataCh <-chan handshakeData, stopCh <-chan struct{}) (*sessionProperties, error) {
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
		return connection.SetReadDeadline(time.Time{})
	})

	err = readEmptyMessage(connection)
	if err != nil {
		return nil, err
	}

	go func() {
		pingTicker := time.NewTicker(pingInterval)
		for {
			select {
			case <-stopCh:
				logger.Debugf("%s: Stop signal received. Close ping ticker", sp.SID)
				pingTicker.Stop()
				return
			case <-pingTicker.C:
				logger.Debugf("%s: Send ping message", sp.SID)
				err := connection.SetReadDeadline(time.Now().Add(pingTimeout))
				if err != nil {
					logger.Errorf("%s: Can't set connection read deadline: %v", sp.SID, err)
					continue
				}

				err = sendMessage(connection, "ping", nil) // Application layer ping?
				if err != nil {
					logger.Errorf("%s: Can't send ping message: %v", sp.SID, err)
					continue
				}

				err = connection.WriteControl(websocket.PingMessage, []byte("ping"), time.Time{})
				if err != nil {
					logger.Errorf("%s: Can't send ping message: %v", sp.SID, err)
					continue
				}
			}
		}
	}()

	logger.Debugf("%s: Join peer network (peerID: %s)", sp.SID, peerMultiaddr.String())
	err = sendMessage(connection, "ss-join", peerMultiaddr.String())
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			select {
			case <-stopCh:
				logger.Debugf("%s: Stop signal received. Close handshake offer sender", sp.SID)
				return
			case offer := <-handshakeDataCh:
				logger.Debugf("%s: Send handshake message: %v", sp.SID, offer.Signal.SDP)
				err = sendMessage(connection, "ss-handshake", offer)
				if err != nil {
					logger.Errorf("%s: Can't send handshake offer: %v", sp.SID, err)
					continue
				}
			}
		}
	}()

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
