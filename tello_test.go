// tello project network_test.go

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
// THE SOFTWARE.package memory

package tello

import (
	"log"
	"testing"
	"time"
)

// use go test -count=1 to bypass test caching

func TestControlConnectDisconnect(t *testing.T) {

	drone := new(Tello)

	_, err := drone.ControlConnectDefault()
	if err != nil {
		log.Fatalf("CCD failed with error %v", err)
	}
	log.Println("Connected to Tello control channel")

	// err = drone.VideoConnectDefault()
	// if err != nil {
	// 	log.Fatalf("Video connect failed with %v", err)
	// }
	// log.Println("Connected to Tello video channel")

	time.Sleep(10 * time.Second)

	//drone.VideoDisconnect()
	drone.ControlDisconnect()
	log.Println("Disconnected normally from Tello")
}

func TestStreamingData(t *testing.T) {
	drone := new(Tello)

	_, err := drone.ControlConnectDefault()
	if err != nil {
		log.Fatalf("CCD failed with error %v", err)
	}
	log.Println("Connected to Tello control channel")

	fdc, err := drone.StreamFlightData(false, 1000)
	if err != nil {
		log.Fatalf("StreamFlighData failed with error %v", err)
	}

	for i := 1; i <= 10; i++ {
		myFD := <-fdc
		log.Printf("Got FlightData with WifiStrength: %d", myFD.WifiStrength)
	}

	drone.ControlDisconnect()
	log.Println("Disconnected normally from Tello")
}

func TestTakeoffLand(t *testing.T) {
	drone := new(Tello)

	_, err := drone.ControlConnectDefault()
	if err != nil {
		log.Fatalf("CCD failed with error %v", err)
	}
	log.Println("Connected to Tello control channel")

	drone.TakeOff()

	time.Sleep(10 * time.Second)

	drone.Land()

	drone.ControlDisconnect()
	log.Println("Disconnected normally from Tello")
}
