package main

import (
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/mtojek/go-libp2p-webrtc-star/testutils"
	"log"
	"os"
	"sync"
	"testing"
	"time"
)

const (
	protocolID = "/examples-echo/1.0.0"
	signalAddr = "/dns4/wrtc-star.discovery.libp2p.io/tcp/443/wss/p2p-webrtc-star"

	waitForStreamTimeout = 5 * time.Minute
)

func main() {
	if len(os.Args) != 2 {
		log.Println("usage: <app> server-id")
		os.Exit(1)
	}

	serverPeerID, err := peer.IDB58Decode(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)

	t := new(testing.T)
	host := testutils.MustCreateHost(t, ctx, signalAddr)
	log.Printf("Client host ID: %s\n", host.ID())

	hostStream := testutils.WaitForStream(t, func() (network.Stream, error) {
		return host.NewStream(ctx, serverPeerID, protocolID)
	}, waitForStreamTimeout)

	buffer := make([]byte, 1024)

	i := 1
	for {
		message := fmt.Sprintf("Simon says - number %d", i)

		_, err := hostStream.Write([]byte(message))
		if err != nil {
			log.Fatal(err)
		}

		_, err = hostStream.Read(buffer)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Echo: %s", string(buffer))

		i++
	}
}
