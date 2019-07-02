package star

import (
	"time"

	"github.com/libp2p/go-libp2p-core/mux"
)

type stream struct {}

var _ mux.MuxedStream = new(stream)

func (s *stream) Read(p []byte) (n int, err error) {
	panic("implement me")
}

func (s *stream) Write(p []byte) (n int, err error) {
	panic("implement me")
}

func (s *stream) Close() error {
	panic("implement me")
}

func (s *stream) Reset() error {
	panic("implement me")
}

func (s *stream) SetDeadline(time.Time) error {
	panic("implement me")
}

func (s *stream) SetReadDeadline(time.Time) error {
	panic("implement me")
}

func (s *stream) SetWriteDeadline(time.Time) error {
	panic("implement me")
}