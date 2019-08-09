package star

import (
	ma "github.com/multiformats/go-multiaddr"
	"github.com/multiformats/go-multiaddr-fmt"
)

const (
	protocolCode = 499
	protocolName = "p2p-webrtc-star"

	wsProtocolCode = 477
	wssProtocolCode = 498
	wssProtocolName = "wss"
)

var (
	protocol = ma.Protocol{
		Code:  protocolCode,
		Name:  protocolName,
		VCode: ma.CodeToVarint(protocolCode),
	}
	protocolMultiaddr ma.Multiaddr
	format = mafmt.And(mafmt.TCP, mafmt.Or(mafmt.Base(wssProtocol.Code), mafmt.Base(wsProtocolCode)), mafmt.Base(protocolCode))

	wssProtocol = ma.Protocol{
		Code:  wssProtocolCode,
		Name:  wssProtocolName,
		VCode: ma.CodeToVarint(wssProtocolCode),
	}
)

func init() {
	err := ma.AddProtocol(wssProtocol)
	if err != nil {
		logger.Fatal(err)
	}

	err = ma.AddProtocol(protocol)
	if err != nil {
		logger.Fatal(err)
	}

	protocolMultiaddr, err = ma.NewMultiaddr("/" + protocolName)
	if err != nil {
		logger.Fatal(err)
	}
}
