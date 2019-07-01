package poc

import (
	"errors"
	"github.com/ipfs/go-log"

	"github.com/multiformats/go-multiaddr"
)

const protocolName = "p2p-webrtc-star"

type signaling struct {
	address multiaddr.Multiaddr
}

var logger = log.Logger("p2p-webrtc-star-poc")

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
	starMultiaddr, err := multiaddr.NewMultiaddr("/" + protocolName)
	if err != nil {
		logger.Fatal(err)
	}

	httpMultiaddr, err := multiaddr.NewMultiaddr("/http")
	if err != nil {
		logger.Fatal(err)
	}

	addr = addr.Decapsulate(starMultiaddr)
	signalAddr := addr.Decapsulate(httpMultiaddr)

	if len(signalAddr.Protocols()) == 0 {
		return nil, errors.New("no protocols defined")
	} else if len(signalAddr.Protocols()) != 1 {
		return nil, errors.New("single signaling server is supported")
	}

	return &signalAddr, nil
}