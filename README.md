# go-libp2p-webrtc-star
libp2p WebRTC transport in Go that includes a discovery mechanism provided by the signalling-star.

Status: **beta**

## Getting started

```bash
$ go get github.com/mtojek/go-libp2p-webrtc-star
```

## Basic example

(see: [examples/common.go](https://github.com/mtojek/go-libp2p-webrtc-star/blob/master/examples/common.go))

```go
func mustCreateHost(t *testing.T, ctx context.Context) host.Host {
	signalMultiaddr := testutils.MustCreateSignalAddr(t, firstSignalAddr)

	privKey := testutils.MustCreatePrivateKey(t)
	identity := testutils.MustCreatePeerIdentity(t, privKey)
	peerstore := pstoremem.NewPeerstore()

	muxer := yamux.DefaultTransport

	starTransport := star.New(identity, peerstore, muxer).
		WithSignalConfiguration(star.SignalConfiguration{
			URLPath: "/socket.io/?EIO=3&transport=websocket",
		}).
		WithWebRTCConfiguration(webrtc.Configuration{
			ICEServers: []webrtc.ICEServer{
				{
					URLs: []string{
						"stun:stun.l.google.com:19302",
						"stun:stun1.l.google.com:19302",
						"stun:stun2.l.google.com:19302",
						"stun:stun3.l.google.com:19302",
						"stun:stun4.l.google.com:19302",
					},
				},
			},
		})

	h, err := libp2p.New(ctx,
		libp2p.Identity(privKey),
		libp2p.ListenAddrs(signalMultiaddr),
		libp2p.Peerstore(peerstore),
		libp2p.Transport(starTransport),
		libp2p.Muxer("/yamux/1.0.0", muxer))
	require.NoError(t, err)
	return h
}
```

## Sample output

*TestSendSingleMessage:*

```bash
=== RUN   TestSendSingleMessage
23:12:42.500 DEBUG p2p-webrtc: Listen on address: /dns4/localhost/tcp/9090/ws/p2p-webrtc-star transport.go:38
23:12:42.500 DEBUG p2p-webrtc: Use signal server: ws://localhost:9090/socket.io/?EIO=3&transport=websocket signal_client.go:14
23:12:42.500 DEBUG p2p-webrtc: Create new listener (address: /dns4/localhost/tcp/9090/ws/p2p-webrtc-star) listener.go:20
23:12:42.500 DEBUG p2p-webrtc: Open new connection: ws://localhost:9090/socket.io/?EIO=3&transport=websocket signal_client.go:181
23:12:42.500 DEBUG p2p-webrtc: Accept connection listener.go:29
23:12:42.500 DEBUG   addrutil: InterfaceAddresses: from manet: [/ip4/127.0.0.1 /ip6/::1 /ip6/fe80::1 /ip6/fe80::14d3:e585:708b:ee8c /ip4/192.168.0.15 /ip6/fe80::b889:42ff:fe1e:afd0 /ip6/fe80::128c:fc0a:b79:7585] addr.go:121
23:12:42.500 DEBUG   addrutil: InterfaceAddresses: usable: [/ip4/127.0.0.1 /ip6/::1 /ip4/192.168.0.15] addr.go:133
23:12:42.500 DEBUG   addrutil: ResolveUnspecifiedAddresses: [/p2p-circuit /dns4/localhost/tcp/9090/ws/p2p-webrtc-star] [/ip4/127.0.0.1 /ip6/::1 /ip4/192.168.0.15] [/p2p-circuit /dns4/localhost/tcp/9090/ws/p2p-webrtc-star] addr.go:109
23:12:42.510 DEBUG p2p-webrtc: Connection to signal server established signal_client.go:54
23:12:42.510 DEBUG p2p-webrtc: _hkOoZ86nDafmyXjAAAW: Ping interval: 25s, Ping timeout: 5s signal_client.go:97
23:12:42.511 DEBUG p2p-webrtc: _hkOoZ86nDafmyXjAAAW: Join peer network (peerID: /dns4/localhost/tcp/9090/ws/p2p-webrtc-star/ipfs/QmZeK8E6g5Ppxars6E8yhyi19aN2yaG2MQTPZwVGwyBnaJ) signal_client.go:141
23:12:42.926 DEBUG p2p-webrtc: Listen on address: /dns4/localhost/tcp/9090/ws/p2p-webrtc-star transport.go:38
23:12:42.926 DEBUG p2p-webrtc: Use signal server: ws://localhost:9090/socket.io/?EIO=3&transport=websocket signal_client.go:14
23:12:42.926 DEBUG p2p-webrtc: Create new listener (address: /dns4/localhost/tcp/9090/ws/p2p-webrtc-star) listener.go:20
23:12:42.926 DEBUG     swarm2: [QmZeK8E6g5Ppxars6E8yhyi19aN2yaG2MQTPZwVGwyBnaJ] opening stream to peer [QmToRGv85ZbYsLTAZokp3fDACg1MM7M15brGJzoabiuHPt] swarm.go:292
23:12:42.926 DEBUG p2p-webrtc: Open new connection: ws://localhost:9090/socket.io/?EIO=3&transport=websocket signal_client.go:181
23:12:42.926 DEBUG p2p-webrtc: Accept connection listener.go:29
23:12:42.926 DEBUG     swarm2: [QmZeK8E6g5Ppxars6E8yhyi19aN2yaG2MQTPZwVGwyBnaJ] swarm dialing peer [QmToRGv85ZbYsLTAZokp3fDACg1MM7M15brGJzoabiuHPt] swarm_dial.go:193
23:12:42.926 DEBUG   addrutil: InterfaceAddresses: from manet: [/ip4/127.0.0.1 /ip6/::1 /ip6/fe80::1 /ip6/fe80::14d3:e585:708b:ee8c /ip4/192.168.0.15 /ip6/fe80::b889:42ff:fe1e:afd0 /ip6/fe80::128c:fc0a:b79:7585] addr.go:121
23:12:42.927 DEBUG   addrutil: InterfaceAddresses: usable: [/ip4/127.0.0.1 /ip6/::1 /ip4/192.168.0.15] addr.go:133
23:12:42.927 DEBUG   addrutil: ResolveUnspecifiedAddresses: [/p2p-circuit /dns4/localhost/tcp/9090/ws/p2p-webrtc-star] [/ip4/127.0.0.1 /ip6/::1 /ip4/192.168.0.15] [/p2p-circuit /dns4/localhost/tcp/9090/ws/p2p-webrtc-star] addr.go:109
23:12:42.934 DEBUG p2p-webrtc: Connection to signal server established signal_client.go:54
23:12:42.934 DEBUG p2p-webrtc: aXccOCbOLI64xrDrAAAX: Ping interval: 25s, Ping timeout: 5s signal_client.go:97
23:12:42.934 DEBUG p2p-webrtc: aXccOCbOLI64xrDrAAAX: Join peer network (peerID: /dns4/localhost/tcp/9090/ws/p2p-webrtc-star/ipfs/QmToRGv85ZbYsLTAZokp3fDACg1MM7M15brGJzoabiuHPt) signal_client.go:141
23:12:42.935 DEBUG p2p-webrtc: _hkOoZ86nDafmyXjAAAW: Received message: ["ws-peer","/dns4/localhost/tcp/9090/ws/p2p-webrtc-star/ipfs/QmToRGv85ZbYsLTAZokp3fDACg1MM7M15brGJzoabiuHPt"] signal_client.go:71
23:12:47.927 DEBUG     swarm2: [QmZeK8E6g5Ppxars6E8yhyi19aN2yaG2MQTPZwVGwyBnaJ] opening stream to peer [QmToRGv85ZbYsLTAZokp3fDACg1MM7M15brGJzoabiuHPt] swarm.go:292
23:12:47.928 DEBUG     swarm2: [QmZeK8E6g5Ppxars6E8yhyi19aN2yaG2MQTPZwVGwyBnaJ] swarm dialing peer [QmToRGv85ZbYsLTAZokp3fDACg1MM7M15brGJzoabiuHPt] swarm_dial.go:193
23:12:47.928 DEBUG     swarm2: QmZeK8E6g5Ppxars6E8yhyi19aN2yaG2MQTPZwVGwyBnaJ swarm dialing QmToRGv85ZbYsLTAZokp3fDACg1MM7M15brGJzoabiuHPt swarm_dial.go:371
23:12:47.928 DEBUG     swarm2: [limiter] adding a dial job through limiter: /dns4/localhost/tcp/9090/ws/p2p-webrtc-star limiter.go:192
23:12:47.928 DEBUG     swarm2: [limiter] taking FD token: peer: QmToRGv85ZbYsLTAZokp3fDACg1MM7M15brGJzoabiuHPt; addr: /dns4/localhost/tcp/9090/ws/p2p-webrtc-star; prev consuming: 0 limiter.go:160
23:12:47.928 DEBUG     swarm2: [limiter] executing dial; peer: QmToRGv85ZbYsLTAZokp3fDACg1MM7M15brGJzoabiuHPt; addr: /dns4/localhost/tcp/9090/ws/p2p-webrtc-star; FD consuming: 1; waiting: 0 limiter.go:166
23:12:47.928 DEBUG     swarm2: QmZeK8E6g5Ppxars6E8yhyi19aN2yaG2MQTPZwVGwyBnaJ swarm dialing QmToRGv85ZbYsLTAZokp3fDACg1MM7M15brGJzoabiuHPt /dns4/localhost/tcp/9090/ws/p2p-webrtc-star swarm_dial.go:455
23:12:47.928 DEBUG p2p-webrtc: Dial peer (ID: QmToRGv85ZbYsLTAZokp3fDACg1MM7M15brGJzoabiuHPt, address: /dns4/localhost/tcp/9090/ws/p2p-webrtc-star) transport.go:29
23:12:48.101 DEBUG p2p-webrtc: Subscribe to the specific handshake (intentID: signal-8032080093292679638) signal_handshakes.go:101
23:12:48.101 DEBUG p2p-webrtc: Send handshake offer (intentID: signal-8032080093292679638) signal_handshakes.go:39
23:12:48.101 DEBUG p2p-webrtc: _hkOoZ86nDafmyXjAAAW: Send handshake message signal_client.go:154
23:12:48.102 DEBUG p2p-webrtc: aXccOCbOLI64xrDrAAAX: Received message: ["ws-handshake",{"intentId":"signal-8032080093292679638","srcMultiaddr":"/dns4/localhost/tcp/9090/ws/p2p-webrtc-star/ipfs/QmZeK8E6g5Ppxars6E8yhyi19aN2yaG2MQTPZwVGwyBnaJ","dstMultiaddr":"/dns4/localhost/tcp/9090/ws/p2p-webrtc-star/ipfs/QmToRGv85ZbYsLTAZokp3fDACg1MM7M15brGJzoabiuHPt","signal":{"type":"offer","sdp":"v=0\r\no=- 372911691 1567458768 IN IP4 0.0.0.0\r\ns=-\r\nt=0 0\r\na=fingerprint:sha-256 1A:18:4E:43:0F:9D:7C:1E:50:3F:6A:95:C9:71:CC:E3:D8:E0:D4:76:36:6E:35:1E:E0:86:1C:ED:A4:00:EF:F9\r\na=group:BUNDLE 0\r\nm=application 9 DTLS/SCTP 5000\r\nc=IN IP4 0.0.0.0\r\na=setup:actpass\r\na=mid:0\r\na=sendrecv\r\na=sctpmap:5000 webrtc-datachannel 1024\r\na=ice-ufrag:BpNcSTXOgFfDQkiI\r\na=ice-pwd:RxVahHZjZiRZvAaWWSjMVKIzLrXZXHau\r\na=candidate:foundation 1 udp 2130706431 192.168.0.15 62101 typ host generation 0\r\na=candidate:foundation 2 udp 2130706431 192.168.0.15 62101 typ host generation 0\r\na=candidate:foundation 1 udp 1694498815 95.160.156.1 37581 typ srflx raddr 0.0.0.0 rport 49439 generation 0\r\na=candidate:foundation 2 udp 1694498815 95.160.156.1 37581 typ srflx raddr 0.0.0.0 rport 49439 generation 0\r\na=candidate:foundation 1 udp 1694498815 95.160.156.1 48114 typ srflx raddr 0.0.0.0 rport 55665 generation 0\r\na=candidate:foundation 2 udp 1694498815 95.160.156.1 48114 typ srflx raddr 0.0.0.0 rport 55665 generation 0\r\na=candidate:foundation 1 udp 1694498815 95.160.156.1 22002 typ srflx raddr 0.0.0.0 rport 52497 generation 0\r\na=candidate:foundation 2 udp 1694498815 95.160.156.1 22002 typ srflx raddr 0.0.0.0 rport 52497 generation 0\r\na=candidate:foundation 1 udp 1694498815 95.160.156.1 58315 typ srflx raddr 0.0.0.0 rport 49240 generation 0\r\na=candidate:foundation 2 udp 1694498815 95.160.156.1 58315 typ srflx raddr 0.0.0.0 rport 49240 generation 0\r\na=candidate:foundation 1 udp 1694498815 95.160.156.1 45984 typ srflx raddr 0.0.0.0 rport 58081 generation 0\r\na=candidate:foundation 2 udp 1694498815 95.160.156.1 45984 typ srflx raddr 0.0.0.0 rport 58081 generation 0\r\na=end-of-candidates\r\n"}}] signal_client.go:71
23:12:48.102 DEBUG p2p-webrtc: Emit handshake data (intentID: signal-8032080093292679638) signal_handshakes.go:77
23:12:48.247 DEBUG p2p-webrtc: aXccOCbOLI64xrDrAAAX: Send handshake message signal_client.go:154
23:12:48.247 DEBUG     swarm2: swarm listener accepted connection: WebRTC connection (ID: connection-6902490575867380988, localPeerID: QmToRGv85ZbYsLTAZokp3fDACg1MM7M15brGJzoabiuHPt, localPeerMultiaddr: /dns4/localhost/tcp/9090/ws/p2p-webrtc-star/ipfs/QmToRGv85ZbYsLTAZokp3fDACg1MM7M15brGJzoabiuHPt, remotePeerID: QmZeK8E6g5Ppxars6E8yhyi19aN2yaG2MQTPZwVGwyBnaJ, remotePeerMultiaddr: /dns4/localhost/tcp/9090/ws/p2p-webrtc-star/ipfs/QmZeK8E6g5Ppxars6E8yhyi19aN2yaG2MQTPZwVGwyBnaJ swarm_listen.go:87
23:12:48.247 DEBUG p2p-webrtc: Accept connection listener.go:29
23:12:48.247 WARNI p2p-webrtc: connection-6902490575867380988: Remote public key undefined connection.go:244
23:12:48.247 DEBUG p2p-webrtc: connection-6902490575867380988: Open stream connection.go:85
23:12:48.248 DEBUG p2p-webrtc: connection-6902490575867380988: Accept stream connection.go:139
23:12:48.248 DEBUG p2p-webrtc: _hkOoZ86nDafmyXjAAAW: Received message: ["ws-handshake",{"intentId":"signal-8032080093292679638","srcMultiaddr":"/dns4/localhost/tcp/9090/ws/p2p-webrtc-star/ipfs/QmZeK8E6g5Ppxars6E8yhyi19aN2yaG2MQTPZwVGwyBnaJ","dstMultiaddr":"/dns4/localhost/tcp/9090/ws/p2p-webrtc-star/ipfs/QmToRGv85ZbYsLTAZokp3fDACg1MM7M15brGJzoabiuHPt","signal":{"type":"answer","sdp":"v=0\r\no=- 606053709 1567458768 IN IP4 0.0.0.0\r\ns=-\r\nt=0 0\r\na=fingerprint:sha-256 D8:E5:18:61:A9:42:8F:50:EB:87:DE:56:A3:F6:C5:DC:11:9A:A0:32:EB:F3:91:33:6B:BB:08:81:AA:BF:C5:D4\r\na=group:BUNDLE 0\r\nm=application 9 DTLS/SCTP 5000\r\nc=IN IP4 0.0.0.0\r\na=setup:active\r\na=mid:0\r\na=sendrecv\r\na=sctpmap:5000 webrtc-datachannel 1024\r\na=ice-ufrag:dvJVVezIDwaxAuyd\r\na=ice-pwd:cbxAaxUCkQmGaNOEhLOBFIxLvCFLHMlI\r\na=candidate:foundation 1 udp 2130706431 192.168.0.15 56753 typ host generation 0\r\na=candidate:foundation 2 udp 2130706431 192.168.0.15 56753 typ host generation 0\r\na=candidate:foundation 1 udp 1694498815 95.160.156.1 4144 typ srflx raddr 0.0.0.0 rport 64283 generation 0\r\na=candidate:foundation 2 udp 1694498815 95.160.156.1 4144 typ srflx raddr 0.0.0.0 rport 64283 generation 0\r\na=candidate:foundation 1 udp 1694498815 95.160.156.1 15361 typ srflx raddr 0.0.0.0 rport 62458 generation 0\r\na=candidate:foundation 2 udp 1694498815 95.160.156.1 15361 typ srflx raddr 0.0.0.0 rport 62458 generation 0\r\na=candidate:foundation 1 udp 1694498815 95.160.156.1 44717 typ srflx raddr 0.0.0.0 rport 55169 generation 0\r\na=candidate:foundation 2 udp 1694498815 95.160.156.1 44717 typ srflx raddr 0.0.0.0 rport 55169 generation 0\r\na=candidate:foundation 1 udp 1694498815 95.160.156.1 29005 typ srflx raddr 0.0.0.0 rport 50880 generation 0\r\na=candidate:foundation 2 udp 1694498815 95.160.156.1 29005 typ srflx raddr 0.0.0.0 rport 50880 generation 0\r\na=candidate:foundation 1 udp 1694498815 95.160.156.1 2391 typ srflx raddr 0.0.0.0 rport 53713 generation 0\r\na=candidate:foundation 2 udp 1694498815 95.160.156.1 2391 typ srflx raddr 0.0.0.0 rport 53713 generation 0\r\na=end-of-candidates\r\n"},"answer":true}] signal_client.go:71
23:12:48.249 DEBUG p2p-webrtc: Emit handshake data (intentID: signal-8032080093292679638) signal_handshakes.go:77
23:12:48.249 DEBUG p2p-webrtc: Handshake answer received (intentID: signal-8032080093292679638) signal_handshakes.go:45
23:12:50.258 DEBUG     swarm2: [limiter] freeing FD token; waiting: 0; consuming: 1 limiter.go:82
23:12:50.258 DEBUG     swarm2: [limiter] freeing peer token; peer QmToRGv85ZbYsLTAZokp3fDACg1MM7M15brGJzoabiuHPt; addr: /dns4/localhost/tcp/9090/ws/p2p-webrtc-star; active for peer: 1; waiting on peer limit: 0 limiter.go:109
23:12:50.258 DEBUG     swarm2: [limiter] clearing all peer dials: QmToRGv85ZbYsLTAZokp3fDACg1MM7M15brGJzoabiuHPt limiter.go:200
23:12:50.258 WARNI p2p-webrtc: connection-7941191111306079968: Remote public key undefined connection.go:244
23:12:50.258 DEBUG     swarm2: network for QmZeK8E6g5Ppxars6E8yhyi19aN2yaG2MQTPZwVGwyBnaJ finished dialing QmToRGv85ZbYsLTAZokp3fDACg1MM7M15brGJzoabiuHPt swarm_dial.go:228
23:12:50.258 DEBUG p2p-webrtc: connection-7941191111306079968: Open stream connection.go:85
23:12:50.258 DEBUG p2p-webrtc: connection-7941191111306079968: Accept stream connection.go:139
23:12:50.258 DEBUG p2p-webrtc: connection-7941191111306079968: Open stream connection.go:85
23:12:50.258 DEBUG p2p-webrtc: connection-7941191111306079968: Accept stream connection.go:139
23:12:50.258 DEBUG p2p-webrtc: connection-6902490575867380988: Accept stream connection.go:139
23:12:50.258 DEBUG  basichost: protocol negotiation took 96.281µs basic_host.go:289
23:12:50.258 DEBUG  basichost: protocol negotiation took 36.66µs basic_host.go:289
23:12:50.258 DEBUG net/identi: QmZeK8E6g5Ppxars6E8yhyi19aN2yaG2MQTPZwVGwyBnaJ sent listen addrs to QmToRGv85ZbYsLTAZokp3fDACg1MM7M15brGJzoabiuHPt: [/dns4/localhost/tcp/9090/ws/p2p-webrtc-star] id.go:317
23:12:50.258 DEBUG net/identi: QmToRGv85ZbYsLTAZokp3fDACg1MM7M15brGJzoabiuHPt sent listen addrs to QmZeK8E6g5Ppxars6E8yhyi19aN2yaG2MQTPZwVGwyBnaJ: [/dns4/localhost/tcp/9090/ws/p2p-webrtc-star] id.go:317
23:12:50.258 DEBUG  basichost: protocol negotiation took 80.861µs basic_host.go:289
23:12:50.258 DEBUG p2p-webrtc: connection-6902490575867380988: Accept stream connection.go:139
23:12:50.259 DEBUG net/identi: /ipfs/id/1.0.0 sent message to QmZeK8E6g5Ppxars6E8yhyi19aN2yaG2MQTPZwVGwyBnaJ /dns4/localhost/tcp/9090/ws/p2p-webrtc-star/ipfs/QmZeK8E6g5Ppxars6E8yhyi19aN2yaG2MQTPZwVGwyBnaJ id.go:214
23:12:50.259 DEBUG net/identi: /ipfs/id/1.0.0 sent message to QmToRGv85ZbYsLTAZokp3fDACg1MM7M15brGJzoabiuHPt /dns4/localhost/tcp/9090/ws/p2p-webrtc-star/ipfs/QmToRGv85ZbYsLTAZokp3fDACg1MM7M15brGJzoabiuHPt id.go:214
--- PASS: TestSendSingleMessage (8.09s)
23:12:50.259 DEBUG net/identi: /ipfs/id/1.0.0 received message from QmToRGv85ZbYsLTAZokp3fDACg1MM7M15brGJzoabiuHPt /dns4/localhost/tcp/9090/ws/p2p-webrtc-star/ipfs/QmToRGv85ZbYsLTAZokp3fDACg1MM7M15brGJzoabiuHPt id.go:230
23:12:50.259 DEBUG net/identi: /ipfs/id/1.0.0 received message from QmZeK8E6g5Ppxars6E8yhyi19aN2yaG2MQTPZwVGwyBnaJ /dns4/localhost/tcp/9090/ws/p2p-webrtc-star/ipfs/QmZeK8E6g5Ppxars6E8yhyi19aN2yaG2MQTPZwVGwyBnaJ id.go:230
PASS
23:12:50.259 DEBUG p2p-webrtc: connection-6902490575867380988: Close connection (no actions) connection.go:199
23:12:50.259 DEBUG net/identi: identify identifying observed multiaddr: /dns4/localhost/tcp/9090/ws/p2p-webrtc-star/ipfs/QmToRGv85ZbYsLTAZokp3fDACg1MM7M15brGJzoabiuHPt [/p2p-circuit /dns4/localhost/tcp/9090/ws/p2p-webrtc-star] id.go:549

Process finished with exit code 0
```

## Development

### Start rendezvous server

```bash
$ npm install --global libp2p-webrtc-star
$ star-signal --port=9090 --host=127.0.0.1
```

### Run unit tests

```bash
$ go test -parallel 1 -json ./... | jq -jr .Output
```
