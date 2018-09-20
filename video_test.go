// tello project video_test.go

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

func TestVideoBitrate(t *testing.T) {
	drone := new(Tello)
	log.Printf("Testing version: %s\n", TelloPackageVersion)
	err := drone.ControlConnectDefault()
	if err != nil {
		log.Fatalf("ControlConnect failed with error %v", err)
	}
	log.Println("Connected to Tello control channel")

	log.Println("Testing without video connection...")

	drone.GetVideoBitrate()
	time.Sleep(3 * time.Second)
	fd := drone.GetFlightData()
	log.Printf("Initial bitrate: %d\n", fd.VideoBitrate)

	drone.SetVideoBitrate(Vbr4M)
	log.Println("Attempt to set to 4mbps")

	time.Sleep(3 * time.Second)
	fd = drone.GetFlightData()
	log.Printf("Video bitrate now: %d\n", fd.VideoBitrate)

	log.Println("Testing with video connection...")

	_, err = drone.VideoConnectDefault()
	if err != nil {
		log.Fatalf("VideoConnect failed with error %v", err)
	}

	drone.GetVideoBitrate()
	time.Sleep(3 * time.Second)
	fd = drone.GetFlightData()
	log.Printf("Video bitrate now: %d\n", fd.VideoBitrate)

	drone.SetVideoBitrate(Vbr4M)
	log.Println("Attempt to set to 4mbps")

	time.Sleep(3 * time.Second)
	fd = drone.GetFlightData()
	log.Printf("Video bitrate now: %d\n", fd.VideoBitrate)

	drone.ControlDisconnect()
	log.Println("Disconnected normally from Tello")
}
