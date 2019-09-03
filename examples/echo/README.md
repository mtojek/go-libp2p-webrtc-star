# Echo client-server

The example constists of two parts - client and server. The server exposes basic *echo* handler that echoes back client's messages.

Star server https://wrtc-star.discovery.libp2p.io/ is used for signalling purpose.

## Getting started

Run server:

```bash
$ go get github.com/mtojek/go-libp2p-webrtc-star/examples/echo/server
$ server
2019/09/03 12:39:24 Server host ID: QmZrEj3jgeF3vUAxPVRSyk1QnZQBbQDNMJqXauLwh6knVr
```

Please note the server host ID and spawn client:
```bash
$ go get github.com/mtojek/go-libp2p-webrtc-star/examples/echo/client
$ client QmZrEj3jgeF3vUAxPVRSyk1QnZQBbQDNMJqXauLwh6knVr
```

Give some time to both endpoints to setup and start handling network stream. You're good to go!

## Sample session

The session has been established between a server running behind a couple of corporate NATs and client deployed on a shared shell account (different locations, VPNs, firewalls):

**Server**:

```bash
$ server
2019/09/03 12:39:24 Server host ID: QmZrEj3jgeF3vUAxPVRSyk1QnZQBbQDNMJqXauLwh6knVr
2019/09/03 12:40:11 Read 21 bytes
2019/09/03 12:40:11 Read 21 bytes
2019/09/03 12:40:11 Read 21 bytes
2019/09/03 12:40:11 Read 21 bytes
2019/09/03 12:40:12 Read 21 bytes
2019/09/03 12:40:12 Read 21 bytes
2019/09/03 12:40:12 Read 21 bytes
2019/09/03 12:40:12 Read 21 bytes
2019/09/03 12:40:12 Read 21 bytes
2019/09/03 12:40:12 Read 22 bytes
2019/09/03 12:40:12 Read 22 bytes
2019/09/03 12:40:12 Read 22 bytes
2019/09/03 12:40:12 Read 22 bytes
2019/09/03 12:40:12 Read 22 bytes
2019/09/03 12:40:12 Read 22 bytes
2019/09/03 12:40:12 Read 22 bytes
2019/09/03 12:40:12 Read 22 bytes
2019/09/03 12:40:12 Read 22 bytes
2019/09/03 12:40:12 Read 22 bytes
2019/09/03 12:40:12 Read 22 bytes
2019/09/03 12:40:12 Read 22 bytes
2019/09/03 12:40:12 Read 22 bytes
...
```

**Client**:

```bash
$ client QmZrEj3jgeF3vUAxPVRSyk1QnZQBbQDNMJqXauLwh6knVr
2019/09/03 12:39:49 Client host ID: QmQhNydYDGYzhuZTwpCSdxznQkfVCGiW5WA6etUUVRXn2w
ice ERROR: 2019/09/03 12:39:49 Failed to enable mDNS, continuing in mDNS disabled mode: (listen udp4 224.0.0.0:5353: bind: operation not permitted)
ice ERROR: 2019/09/03 12:39:59 Failed to enable mDNS, continuing in mDNS disabled mode: (listen udp4 224.0.0.0:5353: bind: operation not permitted)
ice ERROR: 2019/09/03 12:40:04 Failed to enable mDNS, continuing in mDNS disabled mode: (listen udp4 224.0.0.0:5353: bind: operation not permitted)
ice ERROR: 2019/09/03 12:40:09 Failed to enable mDNS, continuing in mDNS disabled mode: (listen udp4 224.0.0.0:5353: bind: operation not permitted)
2019/09/03 12:40:11 Echo: Simon says - number 1
2019/09/03 12:40:11 Echo: Simon says - number 2
2019/09/03 12:40:11 Echo: Simon says - number 3
2019/09/03 12:40:11 Echo: Simon says - number 4
2019/09/03 12:40:11 Echo: Simon says - number 5
2019/09/03 12:40:11 Echo: Simon says - number 6
2019/09/03 12:40:11 Echo: Simon says - number 7
2019/09/03 12:40:11 Echo: Simon says - number 8
2019/09/03 12:40:12 Echo: Simon says - number 9
2019/09/03 12:40:12 Echo: Simon says - number 10
2019/09/03 12:40:12 Echo: Simon says - number 11
2019/09/03 12:40:12 Echo: Simon says - number 12
2019/09/03 12:40:12 Echo: Simon says - number 13
2019/09/03 12:40:12 Echo: Simon says - number 14
2019/09/03 12:40:12 Echo: Simon says - number 15
2019/09/03 12:40:12 Echo: Simon says - number 16
2019/09/03 12:40:12 Echo: Simon says - number 17
2019/09/03 12:40:12 Echo: Simon says - number 18
2019/09/03 12:40:12 Echo: Simon says - number 19
2019/09/03 12:40:12 Echo: Simon says - number 20
2019/09/03 12:40:12 Echo: Simon says - number 21
```
