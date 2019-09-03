package main

import (
	"context"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/mtojek/go-libp2p-webrtc-star/testutils"
	"log"
	"sync"
	"testing"
	"time"
)

const (
	protocolID = "/examples-echo/1.0.0"
	signalAddr = "/dns4/wrtc-star.discovery.libp2p.io/tcp/443/wss/p2p-webrtc-star"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)

	host := testutils.MustCreateHost(new(testing.T), ctx, signalAddr)
	log.Printf("Server host ID: %s\n", host.ID())

	host.SetStreamHandler(protocolID, func(stream network.Stream) {
		message := make([]byte, 1024)

		for {
			n, err := stream.Read(message)
			if err != nil {
				log.Fatal(err)
			} else if n < 1 {
				time.Sleep(1 * time.Second)
				continue
			}

			log.Printf("Read %d bytes", n)
			stream.Write(message)
		}

		wg.Done()
	})

	wg.Wait()
}
