// network.go

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
	"bytes"
	"log"
	"net"
	"strconv"
	"sync"
)

const (
	defaultTelloAddr        = "192.168.10.1"
	defaultTelloControlPort = 8889
	defaultLocalControlPort = 8800
	defaultTelloVideoPort   = 6038
	defaultLocalVideoPort   = 8801
)

type Tello struct {
	ctrlConn, videoConn         *net.UDPConn
	ctrlStopChan, videoStopChan chan bool
	connecting, connected       bool
	ctrlMu                      sync.Mutex
}

// ControlConnect attempts to connect to a Tello at the provided network addr.
// It then starts listening for responses on the control channel.
func (tello *Tello) ControlConnect(udpAddr string, droneUdpPort int, localUdpPort int) (err error) {
	droneAddr, err := net.ResolveUDPAddr("udp", udpAddr+":"+strconv.Itoa(droneUdpPort))
	if err != nil {
		return err
	}
	localAddr, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(localUdpPort))
	if err != nil {
		return err
	}
	tello.ctrlConn, err = net.DialUDP("udp", localAddr, droneAddr)
	if err != nil {
		return err
	}

	// start the control listener Goroutine
	tello.ctrlStopChan = make(chan bool, 2)
	go tello.ControlResponseListener()

	// say hello to the Tello
	tello.sendConnectRequest(defaultTelloVideoPort)

	return nil
}

// ControlConnectDefault attempts to connect to a Tello on the default network addresses.
// It then starts listening for responses on the control channel.
func (tello *Tello) ControlConnectDefault() (err error) {
	return tello.ControlConnect(defaultTelloAddr, defaultTelloControlPort, defaultLocalControlPort)
}

// ControlDisconnect stops the control channel listener and closes the connection to a Tello
func (tello *Tello) ControlDisconnect() {
	// TODO should we tell the Tello we are disconnecting?
	tello.ctrlStopChan <- true
	tello.ctrlConn.Close()
}

func (tello *Tello) VideoConnect(udpAddr string, droneUdpPort int, localUdpPort int) (err error) {
	droneAddr, err := net.ResolveUDPAddr("udp", udpAddr+":"+strconv.Itoa(droneUdpPort))
	if err != nil {
		return err
	}
	// localAddr, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(localUdpPort))
	// if err != nil {
	// 	return err
	// }
	tello.videoConn, err = net.ListenUDP("udp", droneAddr)
	if err != nil {
		return err
	}
	tello.videoStopChan = make(chan bool, 2)
	go tello.VideoResponseListener()
	return nil
}

func (tello *Tello) VideoConnectDefault() (err error) {
	return tello.VideoConnect(defaultTelloAddr, defaultTelloVideoPort, defaultLocalVideoPort)
}

func (tello *Tello) VideoDisconnect() {
	// TODO Shoul we tell the Tello we are stopping video listening?
	tello.videoStopChan <- true
	tello.videoConn.Close()
}

func (tello *Tello) ControlResponseListener() {
	buff := make([]byte, 4096)
	var msgType uint16

	for {
		n, err := tello.ctrlConn.Read(buff)

		// the initial connect response is different...
		if tello.connecting && n == 11 {
			if bytes.ContainsAny(buff, "conn_ack:") {
				// TODO handle returned video port?
				log.Printf("Debug: conn_ack received, buffer len: %d\n", n)
				tello.ctrlMu.Lock()
				tello.connecting = false
				tello.connected = true
				tello.ctrlMu.Unlock()
			} else {
				log.Printf("Unexpected response to connection request <%s>\n", string(buff))
			}
			continue
		}

		select {
		case <-tello.ctrlStopChan:
			log.Println("ControlResponseLister stopped")
			return
		default:
		}
		if err != nil {
			log.Printf("Network Read Error - %v\n", err)
		} else {
			if buff[0] != msgHdr {
				log.Printf("Unexpected network message from Tello <%d>\n", buff[0])
			} else {
				msgType = (uint16(buff[6]) << 8) | uint16(buff[5])
				switch msgType {
				case msgFlightStatus:
				case msgLightStrength:
					log.Println("Light strength received")
				case msgLogHeader:
					log.Println("Log header received")
				case msgWifiStrength:
					log.Println("Wifi strength received")
				default:
					log.Printf("Unknown message type from Tello <%d>\n", msgType)
				}
			}
		}

	}
}

func (tello *Tello) VideoResponseListener() {

}

func (tello *Tello) sendConnectRequest(videoPort uint16) {
	// the initial connect request is different to the usual packets...
	msgBuff := []byte("conn_req:lh")
	msgBuff[9] = byte(videoPort & 0xff)
	msgBuff[10] = byte(videoPort >> 8)
	tello.ctrlMu.Lock()
	tello.connecting = true
	tello.ctrlConn.Write(msgBuff)
	tello.ctrlMu.Unlock()
}

func (tello *Tello) TakeOff() {
	tello.ctrlMu.Lock()
	defer tello.ctrlMu.Unlock()
	// create the command packet
	// populate the command packet
	// send the command packet
}
