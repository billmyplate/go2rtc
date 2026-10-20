package webrtc

import (
	"strings"
	"testing"

	"github.com/AlexxIT/go2rtc/pkg/core"
	"github.com/pion/sdp/v3"
	"github.com/pion/webrtc/v4"
	"github.com/stretchr/testify/require"
)

func TestRingVideoOffer(t *testing.T) {
	api, err := NewRingAPI()
	require.Nil(t, err)

	pc, err := api.NewPeerConnection(webrtc.Configuration{})
	require.Nil(t, err)

	prod := NewConn(pc)
	prod.SetVideoCodecPreferences(RingVideoCodecPreferences())

	medias := []*core.Media{
		{
			Kind:      core.KindAudio,
			Direction: core.DirectionSendRecv,
			Codecs: []*core.Codec{{
				Name:      core.CodecOpus,
				ClockRate: 48000,
				Channels:  2,
			}},
		},
		{
			Kind:      core.KindVideo,
			Direction: core.DirectionRecvonly,
			Codecs: []*core.Codec{{
				Name:      core.CodecH264,
				ClockRate: 90000,
			}},
		},
	}

	offer, err := prod.CreateOffer(medias)
	require.Nil(t, err)
	require.Contains(t, offer, "a=ice-ufrag:")

	sd := &sdp.SessionDescription{}
	require.Nil(t, sd.Unmarshal([]byte(offer)))

	var video *sdp.MediaDescription
	for _, md := range sd.MediaDescriptions {
		if md.MediaName.Media == "video" {
			require.Nil(t, video)
			video = md
		}
	}
	require.NotNil(t, video)

	var h264PT, rtxPT string
	for _, attr := range video.Attributes {
		switch attr.Key {
		case "rtpmap":
			switch {
			case strings.Contains(strings.ToUpper(attr.Value), " H264/"):
				require.Empty(t, h264PT)
				h264PT, _, _ = strings.Cut(attr.Value, " ")
			case strings.Contains(strings.ToLower(attr.Value), " rtx/"):
				require.Empty(t, rtxPT)
				rtxPT, _, _ = strings.Cut(attr.Value, " ")
			default:
				t.Fatalf("unexpected video rtpmap %q", attr.Value)
			}
		case "fmtp":
			lower := strings.ToLower(attr.Value)
			require.NotContains(t, lower, "42001f")
			require.NotContains(t, lower, "42e01f")
		}
	}

	require.NotEmpty(t, h264PT)
	require.NotEmpty(t, rtxPT)
	require.Equal(t, []string{h264PT, rtxPT}, video.MediaName.Formats)

	var h264Fmtp, rtxFmtp string
	for _, attr := range video.Attributes {
		if attr.Key != "fmtp" {
			continue
		}
		pt, line, _ := strings.Cut(attr.Value, " ")
		switch pt {
		case h264PT:
			h264Fmtp = line
		case rtxPT:
			rtxFmtp = line
		}
	}

	require.Contains(t, h264Fmtp, "profile-level-id=640032")
	require.Contains(t, h264Fmtp, "packetization-mode=1")
	require.Contains(t, h264Fmtp, "level-asymmetry-allowed=1")
	require.Equal(t, "apt="+h264PT, rtxFmtp)
	require.NotContains(t, strings.ToUpper(offer), "H265")
	require.NotContains(t, offer, "42001f")
	require.NotContains(t, offer, "42e01f")
}

func TestDefaultVideoOfferKeepsEveryCodec(t *testing.T) {
	api, err := NewAPI()
	require.Nil(t, err)

	pc, err := api.NewPeerConnection(webrtc.Configuration{})
	require.Nil(t, err)

	prod := NewConn(pc)
	offer, err := prod.CreateOffer([]*core.Media{{
		Kind:      core.KindVideo,
		Direction: core.DirectionRecvonly,
	}})
	require.Nil(t, err)
	require.Contains(t, offer, "42001f")
	require.Contains(t, offer, "42e01f")
	require.Contains(t, offer, "640032")
	require.Contains(t, strings.ToUpper(offer), "H265")
	require.NotContains(t, offer, "apt=")
}
