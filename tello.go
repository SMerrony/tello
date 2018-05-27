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
	"errors"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)

const (
	defaultTelloAddr        = "192.168.10.1"
	defaultTelloControlPort = 8889
	defaultLocalControlPort = 8800
	defaultTelloVideoPort   = 6038
	defaultLocalVideoPort   = 8801
)

// Tello holds the current state of a connection to a Tello drone
type Tello struct {
	ctrlConn, videoConn         *net.UDPConn
	ctrlStopChan, videoStopChan chan bool
	connecting, connected       bool
	ctrlMu                      sync.Mutex
	dataMu                      sync.RWMutex
	fd                          FlightData // our private amalgamated store of the latest data
	streamingData               bool       // are we currently sending FlightData out?
}

// ControlConnect attempts to connect to a Tello at the provided network addr.
// It then starts listening for responses on the control channel and waits for the Tello to respond
func (tello *Tello) ControlConnect(udpAddr string, droneUDPPort int, localUDPPort int) (err error) {
	droneAddr, err := net.ResolveUDPAddr("udp", udpAddr+":"+strconv.Itoa(droneUDPPort))
	if err != nil {
		return err
	}
	localAddr, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(localUDPPort))
	if err != nil {
		return err
	}
	tello.ctrlConn, err = net.DialUDP("udp", localAddr, droneAddr)
	if err != nil {
		return err
	}

	// start the control listener Goroutine
	tello.ctrlStopChan = make(chan bool, 2)
	go tello.controlResponseListener()

	// say hello to the Tello
	tello.sendConnectRequest(defaultTelloVideoPort)

	// wait up to 3 seconds for the Tello to respond
	for t := 0; t < 10; t++ {
		if tello.connected {
			break
		}
		time.Sleep(333 * time.Millisecond)
	}
	if !tello.connected {
		return errors.New("Timeout waiting for response to connection request from Tello")
	}

	return nil
}

// ControlConnectDefault attempts to connect to a Tello on the default network addresses.
// It then starts listening for responses on the control channel and waits for the Tello to respond
func (tello *Tello) ControlConnectDefault() (err error) {
	return tello.ControlConnect(defaultTelloAddr, defaultTelloControlPort, defaultLocalControlPort)
}

// ControlDisconnect stops the control channel listener and closes the connection to a Tello
func (tello *Tello) ControlDisconnect() {
	// TODO should we tell the Tello we are disconnecting?
	tello.ctrlStopChan <- true
	tello.ctrlConn.Close()
}

// VideoConnect attempts to connect to a Tello video channel at the provided adrr and starts a listener
func (tello *Tello) VideoConnect(udpAddr string, droneUDPPort int, localUDPPort int) (err error) {
	droneAddr, err := net.ResolveUDPAddr("udp", udpAddr+":"+strconv.Itoa(droneUDPPort))
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
	go tello.videoResponseListener()
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

// GetFlightData returns the current known state of the Tello
func (tello *Tello) GetFlightData() FlightData {
	tello.dataMu.RLock()
	rfd := tello.fd
	tello.dataMu.RUnlock()
	return rfd
}

// StreamFlightData starts a Goroutine which sends FlightData to a channel
// If asAvailable is true then updates are sent whenever fresh data arrives from the Tello and periodMs is ignored
// If asAvailable is false then updates are send every periodMs
// This streamer does not block on the channel, so unconsumed updates are lost
func (tello *Tello) StreamFlightData(asAvailable bool, periodMs time.Duration) <-chan FlightData {
	if tello.streamingData {
		return nil
	}
	fdChan := make(chan FlightData, 2)
	if asAvailable {
		log.Fatal("asAvailable FlightData stream not yet implemented") // TODO
	} else {
		go func() {
			for {
				tello.dataMu.RLock()
				select {
				case fdChan <- tello.fd:
				default:
				}
				tello.dataMu.RUnlock()
				time.Sleep(periodMs * time.Millisecond)
			}
		}()
	}

	return fdChan
}

func (tello *Tello) controlResponseListener() {
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
				pkt := bufferToPacket(buff)
				switch pkt.messageID {
				case msgFlightStatus:
				case msgLightStrength:
					log.Printf("Light strength received - Size: %d, Type: %d\n", pkt.size13, pkt.packetType)
					tello.dataMu.Lock()
					tello.fd.LightStrength = uint8(pkt.payload[0])
					tello.dataMu.Unlock()
				case msgLogHeader:
					log.Printf("Log Header received - Size: %d, Type: %d\n", pkt.size13, pkt.packetType)
				case msgWifiStrength:
					log.Printf("Wifi strength received - Size: %d, Type: %d\n", pkt.size13, pkt.packetType)
					tello.dataMu.Lock()
					tello.fd.WifiStrength = uint8(pkt.payload[0])
					tello.fd.WifiInterference = uint8(pkt.payload[1])
					log.Printf("Parsed Wifi Strength: %d, Interference: %d\n", tello.fd.WifiStrength, tello.fd.WifiInterference)
					tello.dataMu.Unlock()
				default:
					log.Printf("Unknown message type from Tello <%d>\n", msgType)
				}
			}
		}

	}
}

func (tello *Tello) videoResponseListener() {

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

// TakeOff sends a normal takeoff request to the Tello
func (tello *Tello) TakeOff() {
	tello.ctrlMu.Lock()
	defer tello.ctrlMu.Unlock()
	// create the command packet
	// populate the command packet
	// send the command packet
}
