package star

import (
	"time"

	"github.com/libp2p/go-libp2p-core/mux"
)

type Stream struct {}

var _ mux.MuxedStream = new(Stream)

func (s *Stream) Read(p []byte) (n int, err error) {
	panic("implement me")
}

func (s *Stream) Write(p []byte) (n int, err error) {
	panic("implement me")
}

func (s *Stream) Close() error {
	panic("implement me")
}

func (s *Stream) Reset() error {
	panic("implement me")
}

func (s *Stream) SetDeadline(time.Time) error {
	panic("implement me")
}

func (s *Stream) SetReadDeadline(time.Time) error {
	panic("implement me")
}

func (s *Stream) SetWriteDeadline(time.Time) error {
	panic("implement me")
}