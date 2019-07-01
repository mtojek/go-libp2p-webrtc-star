package star

import (
	ma "github.com/multiformats/go-multiaddr"
	"github.com/multiformats/go-multiaddr-fmt"
)

const (
	protocolCode = 499
	protocolName = "p2p-webrtc-star"
)

var (
	protocol = ma.Protocol{
		Code:  protocolCode,
		Name:  protocolName,
		VCode: ma.CodeToVarint(protocolCode),
	}
	format = mafmt.And(wssFormat, mafmt.Base(protocolCode))

	wssProtocol = ma.Protocol{
		Code:  498,
		Name:  "wss",
		VCode: ma.CodeToVarint(498),
	}
	wssFormat = mafmt.And(mafmt.TCP, mafmt.Base(wssProtocol.Code))
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
}
