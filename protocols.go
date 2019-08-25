package star

import (
	ma "github.com/multiformats/go-multiaddr"
	"github.com/multiformats/go-multiaddr-fmt"
)

const (
	ipfsProtocolName = "ipfs"

	webRTCStarProtocolCode = 499
	webRTCStarProtocolName = "p2p-webrtc-star"

	wsProtocolCode = 477

	wssProtocolCode = 498
	wssProtocolName = "wss"
)

var (
	protocol = mustCreateProtocol(webRTCStarProtocolCode, webRTCStarProtocolName)
	protocolMultiaddr ma.Multiaddr

	wssProtocol = mustCreateProtocol(wssProtocolCode, wssProtocolName)

	format = mafmt.And(mafmt.TCP,
		mafmt.Or(mafmt.Base(wssProtocol.Code), mafmt.Base(wsProtocolCode)),
		mafmt.Base(webRTCStarProtocolCode))
)

func init() {
	mustAddProtocol(wssProtocol)
	mustAddProtocol(protocol)
	mustCreateStarMultiaddr()
}

func mustCreateStarMultiaddr() {
	var err error
	protocolMultiaddr, err = ma.NewMultiaddr("/" + webRTCStarProtocolName)
	if err != nil {
		logger.Fatal(err)
	}
}

func mustAddProtocol(protocol ma.Protocol) {
	err := ma.AddProtocol(protocol)
	if err != nil {
		logger.Fatal(err)
	}
}

func mustCreateProtocol(code int, name string) ma.Protocol {
	return ma.Protocol{
		Code:  code,
		Name:  name,
		VCode: ma.CodeToVarint(code),
	}
}
