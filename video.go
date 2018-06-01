// video.go

// Copyright (C) 2018  Steve Merrony

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package tello

import (
	"log"
	"net"
	"strconv"
)

const (
	defaultTelloVideoPort = 6038
	defaultLocalVideoPort = 8801
)

// VideoConnect attempts to connect to a Tello video channel at the provided addr and starts a listener
func (tello *Tello) VideoConnect(udpAddr string, droneUDPPort int, localUDPPort int) (err error) {
	droneAddr, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(droneUDPPort))
	if err != nil {
		return err
	}
	// localAddr, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(localUDPPort))
	// if err != nil {
	// 	return err
	// }
	tello.videoConn, err = net.ListenUDP("udp", droneAddr)
	if err != nil {
		return err
	}
	tello.videoStopChan = make(chan bool, 2)
	tello.VideoChan = make(chan []byte, 100)
	go tello.videoResponseListener()
	log.Println("Video connection setup complete")
	return nil
}

// VideoConnectDefault attempts to connect to a Tello video channel using default addresses, then starts a listener
func (tello *Tello) VideoConnectDefault() (err error) {
	return tello.VideoConnect(defaultTelloAddr, defaultTelloVideoPort, defaultLocalVideoPort)
}

// VideoDisconnect closes the connecttion to the video channel
func (tello *Tello) VideoDisconnect() {
	// TODO Should we tell the Tello we are stopping video listening?
	tello.videoStopChan <- true
	tello.videoConn.Close()
}

func (tello *Tello) videoResponseListener() {
	for {
		vbuf := make([]byte, 2048)
		n, _, err := tello.videoConn.ReadFromUDP(vbuf)
		if err != nil {
			log.Printf("Error reading from video channel - %v\n", err)
		}
		select {
		case tello.VideoChan <- vbuf[2:n]:
		default: // so we don't block
		}
	}
}

// SetVideoBitrate ask the Tello to use the specified bitrate (or auto) for video encoding
func (tello *Tello) SetVideoBitrate(vbr VBR) {
	tello.ctrlMu.Lock()
	defer tello.ctrlMu.Unlock()

	var pkt packet
	// populate the command packet fields we need
	pkt.header = msgHdr
	pkt.toDrone = true
	pkt.packetType = ptSet
	pkt.messageID = msgSetVideoBitrate
	tello.ctrlSeq++
	pkt.sequence = tello.ctrlSeq
	pkt.payload = make([]byte, 1)
	pkt.payload[0] = byte(vbr)
	// pack the packet into raw format and calculate CRCs etc.
	buff := packetToBuffer(pkt)

	// send the command packet
	tello.ctrlConn.Write(buff)
	log.Printf("Set Video Bitrate command sent to drone % x\n", buff)
}

// StartVideo asks the Tello to start sending video
func (tello *Tello) StartVideo() {
	tello.ctrlMu.Lock()
	defer tello.ctrlMu.Unlock()

	var pkt packet
	// populate the command packet fields we need
	pkt.header = msgHdr
	pkt.toDrone = true
	pkt.packetType = ptData2
	pkt.messageID = msgGetVideoSPSPPS
	pkt.sequence = 0 // always zero
	// pack the packet into raw format and calculate CRCs etc.
	buff := packetToBuffer(pkt)
	// send the command packet
	tello.ctrlConn.Write(buff)
}
