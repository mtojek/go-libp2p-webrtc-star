package star

import (
	"github.com/pion/datachannel"
	"io"
	"math"
	"net"
	"time"
)

const wrapperBufferSize = math.MaxUint16

type stream struct {
	id          string
	dataChannel datachannel.ReadWriteCloser
	address     net.Addr

	buffer      []byte
	bufferStart int
	bufferEnd   int
}

var _ net.Conn = new(stream)

func newStream(dataChannel datachannel.ReadWriteCloser, address net.Addr) *stream {
	return &stream{
		id:          createRandomID("stream"),
		dataChannel: dataChannel,
		address:     address,

		buffer: make([]byte, wrapperBufferSize),
	}
}

func (s *stream) Read(p []byte) (int, error) {
	var err error

	if s.bufferEnd == 0 {
		n := 0
		n, err = s.dataChannel.Read(s.buffer)
		if err != nil {
			logger.Debugf("Error occurred while reading from data channel: %v", err)
			err = io.EOF
		}
		s.bufferEnd = n
	}

	n := 0
	if s.bufferEnd-s.bufferStart > 0 {
		n = copy(p, s.buffer[s.bufferStart:s.bufferEnd])
		s.bufferStart += n

		if s.bufferStart >= s.bufferEnd {
			s.bufferStart = 0
			s.bufferEnd = 0
		}
	}
	return n, err
}

func (s *stream) Write(p []byte) (int, error) {
	if len(p) > wrapperBufferSize {
		return s.dataChannel.Write(p[:wrapperBufferSize])
	}
	return s.dataChannel.Write(p)
}

func (s *stream) Close() error {
	logger.Warningf("%s: Close stream", s.id)
	return s.dataChannel.Close()
}

func (s *stream) LocalAddr() net.Addr {
	return s.address
}

func (s *stream) RemoteAddr() net.Addr {
	return s.address
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
