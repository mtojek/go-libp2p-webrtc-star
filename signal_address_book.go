package star

import (
	"github.com/libp2p/go-libp2p-core/peer"
	ma "github.com/multiformats/go-multiaddr"
	"time"
)

type addressBook interface {
	AddAddr(p peer.ID, addr ma.Multiaddr, ttl time.Duration)
}

type selfIgnoreAddressBook struct {
	addressBook addressBook
	ownPeerID peer.ID
}

var _ addressBook = new(selfIgnoreAddressBook)

func (siab *selfIgnoreAddressBook) AddAddr(p peer.ID, addr ma.Multiaddr, ttl time.Duration) {
	if p == siab.ownPeerID {
		logger.Debugf("Do not add own peer ID to the address book (ID: %v)", p)
		return
	}
	siab.addressBook.AddAddr(p, addr, ttl)
}

func decorateSelfIgnoreAddressBook(addressBook addressBook, peerID peer.ID) addressBook {
	return &selfIgnoreAddressBook{
		addressBook: addressBook,
		ownPeerID: peerID,
	}
}