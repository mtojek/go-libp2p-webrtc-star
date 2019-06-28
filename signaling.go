package star

import (
	"errors"

	"github.com/multiformats/go-multiaddr"
)

type signaling struct {
	address multiaddr.Multiaddr
}

func newSignaling(maddr multiaddr.Multiaddr) (*signaling, error) {
	address, err := decapsulate(maddr)
	if err != nil {
		return nil, err
	}

	return &signaling{
		address: *address,
	}, nil
}

func decapsulate(addr multiaddr.Multiaddr) (*multiaddr.Multiaddr, error) {
	starMultiAddr, err := multiaddr.NewMultiaddr("/" + protocolName)
	if err != nil {
		logger.Fatal(err)
	}

	httpMultiAddr, err := multiaddr.NewMultiaddr("/http")
	if err != nil {
		logger.Fatal(err)
	}

	addr = addr.Decapsulate(starMultiAddr)
	signalAddr := addr.Decapsulate(httpMultiAddr)

	if len(signalAddr.Protocols()) == 0 {
		return nil, errors.New("no protocols defined")
	} else if len(signalAddr.Protocols()) != 1 {
		return nil, errors.New("single signaling server is supported")
	}

	return &signalAddr, nil
}