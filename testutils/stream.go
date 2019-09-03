package testutils

import (
	"github.com/libp2p/go-libp2p-core/network"
	"testing"
	"time"
)

func WaitForStream(t *testing.T, newStream func() (network.Stream, error), timeout time.Duration) network.Stream {
	startTime := time.Now()

	var s network.Stream
	var err error

	for time.Now().Before(startTime.Add(timeout)) {
		s, err = newStream()
		if err == nil {
			return s
		}
		time.Sleep(5 * time.Second)
	}

	t.Fatalf("Timeout occurred while waiting for the stream: %v", err)
	return nil
}
