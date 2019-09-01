package star

import (
	"net"

	"github.com/libp2p/go-libp2p-core/transport"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/multiformats/go-multiaddr-net"
)

type listener struct {
	address              ma.Multiaddr
	signal               *signal
	unregisterSignalFunc func(multiaddr ma.Multiaddr) error
}

var _ transport.Listener = new(listener)

func newListener(address ma.Multiaddr, signal *signal, unregisterSignalFunc func(multiaddr ma.Multiaddr) error) (*listener, error) {
	logger.Debugf("Create new listener (address: %s)", address)
	return &listener{
		address:              address,
		signal:               signal,
		unregisterSignalFunc: unregisterSignalFunc,
	}, nil
}

func (l *listener) Accept() (transport.CapableConn, error) {
	logger.Debug("Accept connection")
	return l.signal.accept()
}

func (l *listener) Close() error {
	logger.Debug("Close listener")
	return l.unregisterSignalFunc(l.address)
}

func (l *listener) Addr() net.Addr {
	networkAddress, err := manet.ToNetAddr(l.address)
	if err != nil {
		logger.Fatal(err)
	}
	return networkAddress
}

func (l *listener) Multiaddr() ma.Multiaddr {
	return l.address
}
