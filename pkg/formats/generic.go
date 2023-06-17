package formats

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pion/rtp"
)

func findClockRate(payloadType uint8, rtpMap string) (int, error) {
	// get clock rate from payload type
	// https://en.wikipedia.org/wiki/RTP_payload_formats
	switch payloadType {
	case 0, 1, 2, 3, 4, 5, 7, 8, 9, 12, 13, 15, 18:
		return 8000, nil

	case 6:
		return 16000, nil

	case 10, 11:
		return 44100, nil

	case 14, 25, 26, 28, 31, 32, 33, 34:
		return 90000, nil

	case 16:
		return 11025, nil

	case 17:
		return 22050, nil
	}

	// get clock rate from rtpmap
	// https://tools.ietf.org/html/rfc4566
	// a=rtpmap:<payload type> <encoding name>/<clock rate> [/<encoding parameters>]
	if rtpMap == "" {
		return 0, fmt.Errorf("attribute 'rtpmap' not found")
	}

	tmp := strings.Split(rtpMap, "/")
	if len(tmp) != 2 && len(tmp) != 3 {
		return 0, fmt.Errorf("invalid rtpmap (%v)", rtpMap)
	}

	v, err := strconv.ParseUint(tmp[1], 10, 31)
	if err != nil {
		return 0, err
	}

	return int(v), nil
}

// Generic is a generic RTP format.
type Generic struct {
	PayloadTyp uint8
	RTPMa      string
	FMT        map[string]string

	// clock rate of the format. Filled when calling Init().
	ClockRat int
}

// Init computes the clock rate of the format. It is mandatory to call it.
func (f *Generic) Init() error {
	f.ClockRat, _ = findClockRate(f.PayloadTyp, f.RTPMa)
	return nil
}

func (f *Generic) unmarshal(payloadType uint8, _ string, _ string, rtpmap string, fmtp map[string]string) error {
	f.PayloadTyp = payloadType
	f.RTPMa = rtpmap
	f.FMT = fmtp

	return f.Init()
}

// Codec implements Format.
func (f *Generic) Codec() string {
	return "Generic"
}

// String implements Format.
//
// Deprecated: replaced by Codec().
func (f *Generic) String() string {
	return f.Codec()
}

// ClockRate implements Format.
func (f *Generic) ClockRate() int {
	return f.ClockRat
}

// PayloadType implements Format.
func (f *Generic) PayloadType() uint8 {
	return f.PayloadTyp
}

// RTPMap implements Format.
func (f *Generic) RTPMap() string {
	return f.RTPMa
}

// FMTP implements Format.
func (f *Generic) FMTP() map[string]string {
	return f.FMT
}

// PTSEqualsDTS implements Format.
func (f *Generic) PTSEqualsDTS(*rtp.Packet) bool {
	return true
}
