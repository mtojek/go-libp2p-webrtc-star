package star

import (
	"errors"

	"github.com/multiformats/go-multiaddr"
)

type signaling struct {
	address string
}

func newSignaling(maddr multiaddr.Multiaddr) (*signaling, error) {
	signalingAddress, err := decapsulate(maddr)
	if err != nil {
		return nil, err
	}

	return &signaling{
		address: signalingAddress,
	}, nil
}

func decapsulate(addr multiaddr.Multiaddr) (string, error) {
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
		return "", errors.New("no protocols defined")
	} else if len(signalAddr.Protocols()) != 1 {
		return "", errors.New("single signaling server is supported")
	}

	sa, err := signalAddr.ValueForProtocol(signalAddr.Protocols()[0].Code)
	if err != nil {
		return "", err
	}
	return sa, nil
}