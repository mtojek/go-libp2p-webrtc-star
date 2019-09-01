package star

import (
	"time"

	"github.com/libp2p/go-libp2p-core/mux"
)

type stream struct {
	id string
}

var _ mux.MuxedStream = new(stream)

func newStream() *stream {
	return &stream{
		id: createRandomID("stream"),
	}
}

func (s *stream) Read(p []byte) (n int, err error) {
	panic("implement me") // TODO
}

func (s *stream) Write(p []byte) (n int, err error) {
	panic("implement me") // TODO
}

func (s *stream) Reset() error {
	logger.Debugf("%s: Reset stream", s.id)

	panic("implement me") // TODO
}

func (s *stream) Close() error {
	logger.Warningf("%s: Close stream (no actions)", s.id)
	return nil
}

func (s *stream) SetDeadline(time.Time) error {
	logger.Warningf("%s: Can't set deadline (not implemented)", s.id)
	return nil
}

func (s *stream) SetReadDeadline(time.Time) error {
	logger.Warningf("%s: Can't set read deadline (not implemented)", s.id)
	return nil
}

func (s *stream) SetWriteDeadline(time.Time) error {
	logger.Warningf("%s: Can't set write deadline (not implemented)", s.id)
	return nil
}
