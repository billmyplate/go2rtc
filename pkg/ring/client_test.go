package ring

import "testing"

func TestAnswerIsH264(t *testing.T) {
	h264 := "m=video 9 UDP/TLS/RTP/SAVPF 98\r\na=rtpmap:98 H264/90000\r\n"
	if err := answerIsH264(h264); err != nil {
		t.Fatal(err)
	}

	hevc := "m=video 9 UDP/TLS/RTP/SAVPF 100\r\na=rtpmap:100 H265/90000\r\n"
	if err := answerIsH264(hevc); err == nil {
		t.Fatal("expected H265 answer to be rejected")
	}

	both := h264 + "a=rtpmap:100 H265/90000\r\n"
	if err := answerIsH264(both); err == nil {
		t.Fatal("expected mixed answer to be rejected")
	}
}
