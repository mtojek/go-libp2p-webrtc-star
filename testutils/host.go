package testutils

import (
	"context"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-peerstore/pstoremem"
	yamux "github.com/libp2p/go-libp2p-yamux"
	"github.com/mtojek/go-libp2p-webrtc-star"
	"github.com/pion/webrtc"
	"github.com/stretchr/testify/require"
	"testing"
)

func MustCreateHost(t *testing.T, ctx context.Context, signalAddr string) host.Host {
	signalMultiaddr := MustCreateSignalAddr(t, signalAddr)

	privKey := MustCreatePrivateKey(t)
	identity := MustCreatePeerIdentity(t, privKey)
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
