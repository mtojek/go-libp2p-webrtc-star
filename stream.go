package star

import (
	"github.com/libp2p/go-libp2p-core/mux"
	"github.com/pion/datachannel"
	"io"
	"time"
)

type stream struct {
	id          string
	dataChannel datachannel.ReadWriteCloser
}

var _ mux.MuxedStream = new(stream)

func newStream(dataChannel datachannel.ReadWriteCloser) *stream {
	return &stream{
		id:          createRandomID("stream"),
		dataChannel: dataChannel,
	}
}

func (s *stream) Read(p []byte) (int, error) {
	i, err := s.dataChannel.Read(p)
	if err != nil {
		return i, io.EOF
	}
	return i, nil
}

func (s *stream) Write(p []byte) (n int, err error) {
	return s.dataChannel.Write(p)
}

func (s *stream) Reset() error {
	logger.Debugf("%s: Reset stream", s.id)
	return s.dataChannel.Close()
}

func (s *stream) Close() error {
	logger.Warningf("%s: Close stream (no actions)", s.id)
	return nil
}

func (s *stream) SetDeadline(time.Time) error {
	return nil
}

func (s *stream) SetReadDeadline(time.Time) error {
	return nil
}

func (s *stream) SetWriteDeadline(time.Time) error {
	return nil
}
