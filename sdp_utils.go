package sipgox

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/emiago/sipgox/sdp"
)

func GetCurrentNTPTimestamp() uint64 {
	ntpEpochOffset := 2208988800 // Offset from Unix epoch (January 1, 1970) to NTP epoch (January 1, 1900)
	currentTime := time.Now().Unix() + int64(ntpEpochOffset)

	return uint64(currentTime)
}

type SDPMode string

const (
	// https://datatracker.ietf.org/doc/html/rfc4566#section-6
	SDPModeRecvonly SDPMode = "recvonly"
	SDPModeSendrecv SDPMode = "sendrecv"
	SDPModeSendonly SDPMode = "sendonly"
)

func SDPGenerateForAudio(originIP net.IP, connectionIP net.IP, rtpPort int, mode SDPMode, fmts sdp.Formats) []byte {
	ntpTime := GetCurrentNTPTimestamp()

	formatsMap := []string{}
	for _, f := range fmts {
		switch f {
		case "0":
			formatsMap = append(formatsMap, "a=rtpmap:0 PCMU/8000")
		case "8":
			formatsMap = append(formatsMap, "a=rtpmap:8 PCMA/8000")

			// TODO add more here
		}
	}
	// Support only ulaw and alaw
	s := []string{
		"v=0",
		fmt.Sprintf("o=user1 %d %d IN IP4 %s", ntpTime, ntpTime, originIP),
		"s=Sip Go Media",
		// "b=AS:84",
		fmt.Sprintf("c=IN IP4 %s", connectionIP),
		"t=0 0",
		fmt.Sprintf("m=audio %d RTP/AVP %s", rtpPort, strings.Join(fmts, " ")+" 101"),
		"a=" + string(mode),
		// "a=ssrc:111222 cname:user@example.com",
		// "a=rtpmap:0 PCMU/8000",
		// "a=rtpmap:8 PCMA/8000",
		// THIS is FOR DTM
		"a=rtpmap:101 telephone-event/8000",
		"a=fmtp:101 0-16",
		// "",
		// "a=rtpmap:120 telephone-event/16000",
		// "a=fmtp:120 0-16",
		// "a=rtpmap:121 telephone-event/8000",
		// "a=fmtp:121 0-16",
		// "a=rtpmap:122 telephone-event/32000",
		// "a=rtcp-mux",
		// fmt.Sprintf("a=rtcp:%d IN IP4 %s", rtpPort+1, connectionIP),
	}

	s = append(s, formatsMap...)

	// s := []string{
	// 	"v=0",
	// 	fmt.Sprintf("o=- %d %d IN IP4 %s", ntpTime, ntpTime, originIP),
	// 	"s=Sip Go Media",
	// 	// "b=AS:84",
	// 	fmt.Sprintf("c=IN IP4 %s", connectionIP),
	// 	"t=0 0",
	// 	fmt.Sprintf("m=audio %d RTP/AVP 96 97 98 99 3 0 8 9 120 121 122", rtpPort),
	// 	"a=" + string(mode),
	// 	"a=rtpmap:96 speex/16000",
	// 	"a=rtpmap:97 speex/8000",
	// 	"a=rtpmap:98 speex/32000",
	// 	"a=rtpmap:99 iLBC/8000",
	// 	"a=fmtp:99 mode=30",
	// 	"a=rtpmap:120 telephone-event/16000",
	// 	"a=fmtp:120 0-16",
	// 	"a=rtpmap:121 telephone-event/8000",
	// 	"a=fmtp:121 0-16",
	// 	"a=rtpmap:122 telephone-event/32000",
	// 	"a=rtcp-mux",
	// 	fmt.Sprintf("a=rtcp:%d IN IP4 %s", rtpPort+1, connectionIP),
	// }

	res := strings.Join(s, "\r\n")
	return []byte(res)
}
