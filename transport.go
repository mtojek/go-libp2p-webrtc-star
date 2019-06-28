package star

import (
	"github.com/multiformats/go-multiaddr"
)

const (
	protocolCode = 499
	protocolName = "p2p-webrtc-star"
)

var protocol = &multiaddr.Protocol{
	Code:  protocolCode,
	Name:  protocolName,
	VCode: multiaddr.CodeToVarint(protocolCode),
}

func init() {
	err := multiaddr.AddProtocol(*protocol)
	if err != nil {
		logger.Fatal(err)
	}
}

func Protocol() *multiaddr.Protocol {
	return protocol
}

func Transport() *WebRTCStar {
	return new(WebRTCStar)
}