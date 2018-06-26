// autopilot_test.go

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
	"testing"
	"time"
)

func TestGotoHeight(t *testing.T) {
	drone := new(Tello)

	err := drone.ControlConnectDefault()
	if err != nil {
		log.Fatalf("CCD failed with error %v", err)
	}
	log.Println("Connected to Tello control channel")
	time.Sleep(2 * time.Second)

	drone.TakeOff()
	time.Sleep(5 * time.Second)

	if _, err = drone.GotoHeight(5); err != nil { // should go down to .5m
		t.Errorf("Error %v calling GotoHeight(5)", err)
	}
	time.Sleep(5 * time.Second)

	done, err := drone.GotoHeight(15)
	if err != nil { // should go up to 1.5m
		t.Errorf("Error %v calling GotoHeight(15)", err)
	}
	<-done
	log.Println("Navigation completion notified")

	drone.Land()

	drone.ControlDisconnect()
	log.Println("Disconnected normally from Tello")
}
