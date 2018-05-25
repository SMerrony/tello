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
	ctrlMu                      sync.Mutex
}

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
	tello.ctrlStopChan = make(chan bool, 2)
	go tello.ControlResponseListener()
	return nil
}

func (tello *Tello) ControlConnectDefaultTello() (err error) {
	return tello.ControlConnect(defaultTelloAddr, defaultTelloControlPort, defaultLocalControlPort)
}

func (tello *Tello) ControlDisconnect() {
	tello.ctrlStopChan <- true
	tello.ctrlConn.Close()
}

func (tello *Tello) VideoConnect(udpAddr string, droneUdpPort int, localUdpPort int) (err error) {
	droneAddr, err := net.ResolveUDPAddr("udp", udpAddr+":"+strconv.Itoa(droneUdpPort))
	if err != nil {
		return err
	}
	localAddr, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(localUdpPort))
	if err != nil {
		return err
	}
	tello.videoConn, err = net.DialUDP("udp", localAddr, droneAddr)
	if err != nil {
		return err
	}
	tello.videoStopChan = make(chan bool, 2)
	go tello.VideoResponseListener()
	return nil
}

func (tello *Tello) VideoConnectDefaultTello() (err error) {
	return tello.VideoConnect(defaultTelloAddr, defaultTelloVideoPort, defaultLocalVideoPort)
}

func (tello *Tello) VideoDisconnect() {
	tello.videoStopChan <- true
	tello.videoConn.Close()
}

func (tello *Tello) ControlResponseListener() {
	buff := make([]byte, 4096)
	var msgType uint16

	for {
		_, err := tello.ctrlConn.Read(buff)
		select {
		case <-tello.ctrlStopChan:
			log.Println("ControlResponseLister stopped")
			return
		defaultTello:
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
				defaultTello:
					log.Printf("Unknown message type from Tello <%d>\n", msgType)
				}
			}
		}

	}
}

func (tello *Tello) TakeOff() {
	tello.ctrlMu.Lock()
	defer tello.ctrlMu.Unlock()
	// create the command packet
	// populate the command packet
	// send the command packet
}
