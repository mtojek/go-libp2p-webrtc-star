package star

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/libp2p/go-libp2p-core/peer"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/pion/webrtc"
	"math"
	"math/rand"
	"sync"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

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
		logger.Debugf("Handshake answer received (intentID: %s)", intentID)
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
	k := rand.Intn(math.MaxInt64-1000000000000000000) + 1000000000000000000
	return fmt.Sprintf("pysio-%d", k)
}

type handshakeSubscription struct {
	m sync.RWMutex

	subscribers map[string]chan handshakeData
	sink        chan handshakeData
}

func newHandshakeSubscription() *handshakeSubscription {
	return &handshakeSubscription{
		subscribers: map[string]chan handshakeData{},
		sink:        make(chan handshakeData),
	}
}

func (hs *handshakeSubscription) emit(answer handshakeData) {
	logger.Debugf("Emit handshake answer (intentID: %s)", answer.IntentID)

	hs.m.RLock()
	_, ok := hs.subscribers[answer.IntentID]
	hs.m.RUnlock()

	if ok {
		hs.m.Lock()
		defer hs.m.Unlock()

		c, ok := hs.subscribers[answer.IntentID]
		if ok {
			c <- answer
			delete(hs.subscribers, answer.IntentID)
			close(c)
		}
		return
	}
	hs.sink <- answer
}

func (hs *handshakeSubscription) unsubscribed() <-chan handshakeData {
	return hs.sink
}

func (hs *handshakeSubscription) subscribe(intentID string) <-chan handshakeData {
	logger.Debugf("Subscribe to the specific handshake (intentID: %s)", intentID)

	hs.m.Lock()
	defer hs.m.Unlock()

	hs.subscribers[intentID] = make(chan handshakeData)
	return hs.subscribers[intentID]
}

func (hs *handshakeSubscription) cancel(intentID string) {
	logger.Debugf("Cancel handshake subscription (intentID: %s)", intentID)

	hs.m.Lock()
	defer hs.m.Unlock()

	c := hs.subscribers[intentID]
	delete(hs.subscribers, intentID)
	close(c)
}
