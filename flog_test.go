// tello project flog_test.go

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

func TestAckHeader(t *testing.T) {

	drone := new(Tello)
	log.Printf("Testing version: %s\n", TelloPackageVersion)
	err := drone.ControlConnectDefault()
	if err != nil {
		log.Fatalf("CCD failed with error %v", err)
	}
	log.Println("Connected to Tello control channel")

	time.Sleep(5 * time.Second)

	drone.ControlDisconnect()
	log.Println("Disconnected normally from Tello")
}

func TestQuatToEulerDeg(t *testing.T) {
	p, r, y := QuatToEulerDeg(0, 0, 0, 1) // straight ahead
	log.Printf("p: %f, r: %f, y: %f\n", p, r, y)
	if p != 0 || r != 0 || y != 0 {
		t.Errorf("p: %f, r: %f, y: %f\n", p, r, y)
	}

	p, r, y = QuatToEulerDeg(0, 0.7071, 0, 0.7071) // 90 deg about  y axis
	log.Printf("p: %f, r: %f, y: %f\n", p, r, y)
	if p != 90 || r != 0 || y != 0 {
		t.Errorf("p: %f, r: %f, y: %f\n", p, r, y)
	}

	p, r, y = QuatToEulerDeg(0, 0, 1, 1) // roll 117 deg
	log.Printf("p: %f, r: %f, y: %f\n", p, r, y)
	if p != 0 || r != 0 || y != 117 {
		t.Errorf("p: %f, r: %f, y: %f\n", p, r, y)
	}
}

func TestQuatToYawDeg(t *testing.T) {
	y := quatToYawDeg(0, 0, 0, 1) // straight ahead
	log.Printf("y: %f\n", y)
	if y != 0 {
		t.Errorf("y: %f\n", y)
	}

	y = quatToYawDeg(0, 0.7071, 0, 0.7071) // 90 deg about  y axis
	log.Printf("y: %f\n", y)
	if y != 0 {
		t.Errorf("y: %f\n", y)
	}

	y = quatToYawDeg(0, 0, 1, 1) // roll 117 deg
	log.Printf("y: %f\n", y)
	if y != 117 {
		t.Errorf("y: %f\n", y)
	}
}
