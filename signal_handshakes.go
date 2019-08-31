package star

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/libp2p/go-libp2p-core/peer"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/pion/webrtc"
	"time"
)

type handshakeData struct {
	IntentID     string                    `json:"intentId,omitempty"`
	SrcMultiaddr string                    `json:"srcMultiaddr"`
	DstMultiaddr string                    `json:"dstMultiaddr"`
	Signal       webrtc.SessionDescription `json:"signal"`
	Answer       bool                      `json:"answer,omitempty"`
}

func (hd *handshakeData) String() string {
	m, err := json.Marshal(hd)
	if err != nil {
		logger.Error("can't marshal handshake data")
		return ""
	}
	return string(m)
}

func (s *signal) doHandshake(destinationPeerID peer.ID, offerDescription webrtc.SessionDescription) (webrtc.SessionDescription, error) {
	dstMultiaddr, err := ma.NewMultiaddr(fmt.Sprintf("/%s/%s", ipfsProtocolName, destinationPeerID.String()))
	if err != nil {
		return webrtc.SessionDescription{}, err
	}
	intentID := createRandomIntentID()
	s.handshakeDataCh <- handshakeData{
		IntentID:     intentID,
		DstMultiaddr: s.signalMultiaddr.Encapsulate(dstMultiaddr).String(),
		SrcMultiaddr: s.peerMultiaddr.String(),
		Signal:       offerDescription,
	}

	timeout := time.After(handshakeAnswerTimeout)
	select {
	case answer := <-s.handshakeSubscription.subscribe(intentID):
		return answer.Signal, nil
	case <-timeout:
		s.handshakeSubscription.cancel(intentID)
		return webrtc.SessionDescription{}, errors.New("handshake answer timeout")
	}
}

func (s *signal) answerHandshake(intentID string, dstMultiaddr string, answerDescription webrtc.SessionDescription) {
	s.handshakeDataCh <- handshakeData{
		IntentID:     intentID,
		DstMultiaddr: dstMultiaddr,
		SrcMultiaddr: s.peerMultiaddr.String(),
		Signal:       answerDescription,
		Answer:       true,
	}
}

func createRandomIntentID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

type handshakeSubscription struct{}

func newHandshakeSubscription() *handshakeSubscription {
	return new(handshakeSubscription)
}

func (hs *handshakeSubscription) emit(answer handshakeData) {
	logger.Debugf("Emit handshake answer (intentID: %s)", answer.IntentID)
	// TODO
}

func (hs *handshakeSubscription) unsubscribed() <-chan handshakeData {
	// TODO
	return nil
}

func (hs *handshakeSubscription) subscribe(intentID string) <-chan handshakeData {
	logger.Debugf("Subscribe to the specific answer (intentID: %s)", intentID)
	// TODO
	return nil
}

func (hs *handshakeSubscription) cancel(intentID string) {
	logger.Debugf("Cancel handshake subscription (intentID: %s)", intentID)
	// TODO
}
