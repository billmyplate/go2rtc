package webrtc

import (
	"testing"

	"github.com/AlexxIT/go2rtc/pkg/core"
	"github.com/pion/webrtc/v4"
	"github.com/stretchr/testify/require"
)

func TestGetTrackReusesReceiverForSameCodec(t *testing.T) {
	api, err := NewAPI()
	require.Nil(t, err)

	pc, err := api.NewPeerConnection(webrtc.Configuration{})
	require.Nil(t, err)

	conn := NewConn(pc)
	conn.Mode = core.ModeActiveProducer

	media := &core.Media{
		Kind:      core.KindVideo,
		Direction: core.DirectionRecvonly,
		Codecs: []*core.Codec{
			{Name: core.CodecH264, ClockRate: 90000, PayloadType: 98, FmtpLine: "profile-level-id=640032"},
			{Name: core.CodecH264, ClockRate: 90000, PayloadType: 96, FmtpLine: "profile-level-id=42001f"},
			{Name: core.CodecH265, ClockRate: 90000, PayloadType: 100},
		},
	}

	bound, err := conn.GetTrack(media, media.Codecs[0])
	require.Nil(t, err)

	otherPayload, err := conn.GetTrack(media, media.Codecs[1])
	require.Nil(t, err)
	require.Same(t, bound, otherPayload)
	require.Len(t, conn.Receivers, 1)

	hevc, err := conn.GetTrack(media, media.Codecs[2])
	require.Nil(t, err)
	require.NotSame(t, bound, hevc)
	require.Len(t, conn.Receivers, 2)
}
