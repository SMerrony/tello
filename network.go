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
)

const (
	defaultAddr        = "192.168.10.1"
	defaultCommandPort = 8889
	defaultLocalPort   = 8888
	defaultVideoPort   = 6038
)

func ControlConnect(udpAddr string, droneUdpPort int, localUdpPort int) (ctrlConn *net.UDPConn, err error) {
	droneAddr, err := net.ResolveUDPAddr("udp", udpAddr+":"+strconv.Itoa(droneUdpPort))
	if err != nil {
		return nil, err
	}
	localAddr, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(localUdpPort))
	if err != nil {
		return nil, err
	}
	ctrlConn, err = net.DialUDP("udp", localAddr, droneAddr)
	if err != nil {
		return nil, err
	}
	return ctrlConn, nil
}

func ControlConnectDefault() (ctrlConn *net.UDPConn, err error) {
	return ControlConnect(defaultAddr, defaultCommandPort, defaultLocalPort)
}

func ControlDisconnect(ctrlConn *net.UDPConn) {
	ctrlConn.Close()
}

func ControlResponseListener(ctrlConn *net.UDPConn, stop <-chan bool) {
	buff := make([]byte, 4096)
	var msgType uint16

	for {
		_, err := ctrlConn.Read(buff)
		select {
		case <-stop:
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

				default:
					log.Printf("Unknown message type from Tello <%d>\n", msgType)
				}
			}
		}

	}
}
