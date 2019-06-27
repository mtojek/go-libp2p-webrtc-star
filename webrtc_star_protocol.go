package star

import (
	"github.com/multiformats/go-multiaddr"
	"log"
)

const protocolCode = 499

var WebRTCStarProtocol = multiaddr.Protocol{
	Code:  protocolCode,
	Name:  "p2p-webrtc-star",
	VCode: multiaddr.CodeToVarint(protocolCode),
}

func init() {
	err := multiaddr.AddProtocol(WebRTCStarProtocol)
	if err != nil {
		log.Fatal(err)
	}
}