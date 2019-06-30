package star

import (
	wss "github.com/mtojek/go-wss-transport"
	"github.com/multiformats/go-multiaddr"
	"github.com/multiformats/go-multiaddr-fmt"
)

const (
	protocolCode = 499
	protocolName = "p2p-webrtc-star"
)

var (
	format = mafmt.And(wss.WssFmt, mafmt.Base(protocolCode))

	protocol = &multiaddr.Protocol{
		Code:  protocolCode,
		Name:  protocolName,
		VCode: multiaddr.CodeToVarint(protocolCode),
	}
)

func init() {
	err := multiaddr.AddProtocol(*protocol)
	if err != nil {
		logger.Fatal(err)
	}
}
