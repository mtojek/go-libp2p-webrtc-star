package star

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"github.com/libp2p/go-libp2p-core/peer"
	ma "github.com/multiformats/go-multiaddr"
	"time"
)

const wsPeerAliveTTL = 60 * time.Second

func processMessage(addressBook addressBook, message []byte) error {
	if bytes.Index(message, []byte(`["ws-peer",`)) > -1 {
		var m []string
		err := json.Unmarshal(message, &m)
		if err != nil {
			return err
		}
		return processWsPeerMessage(addressBook, m)
	}
	return errors.New("tried to processed unknown message")
}

func processWsPeerMessage(addressBook addressBook, wsPeerMessage []string) error {
	if len(wsPeerMessage) < 2 {
		return errors.New("missing peer information")
	}

	peerID, signalMultiaddr, err := extractPeerDestination(wsPeerMessage[1])
	if err != nil {
		return err
	}

	addressBook.AddAddr(peerID, signalMultiaddr, wsPeerAliveTTL)
	return nil
}

func extractPeerDestination(peerAddr string) (peer.ID, ma.Multiaddr, error) {
	peerMultiaddr, err := ma.NewMultiaddr(peerAddr)
	if err != nil {
		return "", nil, err
	}

	value, err := peerMultiaddr.ValueForProtocol(ma.P_IPFS)
	if err != nil {
		return "", nil, err
	}

	peerID, err := peer.IDB58Decode(value)
	if err != nil {
		return "", nil, err
	}

	ipfsMultiaddr, err := ma.NewMultiaddr("/ipfs/" + peerID.String())
	if err != nil {
		return "", nil, err
	}
	return peerID, peerMultiaddr.Decapsulate(ipfsMultiaddr), nil
}

func readMessage(connection *websocket.Conn) ([]byte, error) {
	_, message, err := connection.ReadMessage()
	if err != nil {
		return nil, err
	}

	i := bytes.IndexAny(message, "[{")
	if i < 0 {
		return nil, errors.New("message token not found")
	}
	return message[i:], nil
}

func sendMessage(connection *websocket.Conn, messageType string, messageBody interface{}) error {
	var buffer bytes.Buffer
	buffer.WriteString(messagePrefix)
	buffer.WriteString(`["`)
	buffer.WriteString(messageType)
	buffer.WriteByte('"')

	if messageBody != nil {
		b, err := json.Marshal(messageBody)
		if err != nil {
			return err
		}

		buffer.WriteByte(',')
		buffer.Write(b)
	}
	buffer.WriteByte(']')
	return connection.WriteMessage(websocket.TextMessage, buffer.Bytes())
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

