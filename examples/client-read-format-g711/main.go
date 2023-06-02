package main

import (
	"log"

	"github.com/bluenviron/gortsplib/v3"
	"github.com/bluenviron/gortsplib/v3/pkg/formats"
	"github.com/bluenviron/gortsplib/v3/pkg/url"
	"github.com/pion/rtp"
)

// This example shows how to
// 1. connect to a RTSP server
// 2. check if there's a G711 media
// 3. get G711 frames of that media

func main() {
	c := gortsplib.Client{}

	// parse URL
	u, err := url.Parse("rtsp://localhost:8554/mystream")
	if err != nil {
		panic(err)
	}

	// connect to the server
	err = c.Start(u.Scheme, u.Host)
	if err != nil {
		panic(err)
	}
	defer c.Close()

	// find published medias
	medias, baseURL, _, err := c.Describe(u)
	if err != nil {
		panic(err)
	}

	// find the G711 media and format
	var forma *formats.G711
	medi := medias.FindFormat(&forma)
	if medi == nil {
		panic("media not found")
	}

	// create decoder
	rtpDec, err := forma.CreateDecoder2()
	if err != nil {
		panic(err)
	}

	// setup a single media
	_, err = c.Setup(medi, baseURL, 0, 0)
	if err != nil {
		panic(err)
	}

	// called when a RTP packet arrives
	c.OnPacketRTP(medi, forma, func(pkt *rtp.Packet) {
		// extract G711 frames from RTP packets
		op, _, err := rtpDec.Decode(pkt)
		if err != nil {
			log.Printf("ERR: %v", err)
			return
		}

		// print
		log.Printf("received G711 frame of size %d\n", len(op))
	})

	// start playing
	_, err = c.Play(nil)
	if err != nil {
		panic(err)
	}

	// wait until a fatal error
	panic(c.Wait())
}
