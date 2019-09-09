package star

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/pion/webrtc/v2"
	"math/rand"
	"sync"
	"time"
)

const handshakeAnswerTimeout = 5 * time.Minute

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

func (s *signal) doHandshake(ctx context.Context, offer handshakeData) (handshakeData, error) {
	subscription := s.handshakeSubscription.subscribe(offer.IntentID)

	logger.Debugf("Send handshake offer (intentID: %s)", offer.IntentID)
	s.handshakeDataCh <- offer

	timeout := time.After(handshakeAnswerTimeout)
	select {
	case answer := <-subscription:
		logger.Debugf("Handshake answer received (intentID: %s)", offer.IntentID)
		return answer, nil
	case <-ctx.Done():
		logger.Debugf("Cancel handshake (intentID: %s)", offer.IntentID)
		s.handshakeSubscription.cancel(offer.IntentID)
		return handshakeData{}, errors.New("handshake canceled")
	case <-timeout:
		logger.Debugf("Handshake timeout (intentID: %s)", offer.IntentID)
		s.handshakeSubscription.cancel(offer.IntentID)
		return handshakeData{}, errors.New("handshake answer timeout")
	}
}

func (s *signal) answerHandshake(answer handshakeData) {
	s.handshakeDataCh <- answer
}

type handshakeSubscription struct {
	m sync.Mutex

	subscribers map[string]chan handshakeData
	sink        chan handshakeData
}

func newHandshakeSubscription() *handshakeSubscription {
	return &handshakeSubscription{
		subscribers: map[string]chan handshakeData{},
		sink:        make(chan handshakeData),
	}
}

func (hs *handshakeSubscription) emit(data handshakeData) {
	logger.Debugf("Emit handshake data (intentID: %s)", data.IntentID)

	hs.m.Lock()
	defer hs.m.Unlock()

	if c, ok := hs.subscribers[data.IntentID]; ok {
		c <- data
		delete(hs.subscribers, data.IntentID)
		close(c)
		return
	}

	if !data.Answer {
		hs.sink <- data
	} else {
		logger.Debugf("Received answer to probably cancelled handshake (intentID: %s)", data.IntentID)
	}
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
