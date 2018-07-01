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

func TestAutoFlyToHeight(t *testing.T) {
	drone := new(Tello)

	err := drone.ControlConnectDefault()
	if err != nil {
		log.Fatalf("CCD failed with error %v", err)
	}
	log.Println("Connected to Tello control channel")
	time.Sleep(2 * time.Second)

	drone.TakeOff()
	time.Sleep(5 * time.Second)

	if _, err = drone.AutoFlyToHeight(5); err != nil { // should go down to .5m
		t.Errorf("Error %v calling AutoFlyToHeight(5)", err)
	}
	time.Sleep(5 * time.Second)

	done, err := drone.AutoFlyToHeight(15)
	if err != nil { // should go up to 1.5m
		t.Errorf("Error %v calling AutoFlyToHeight(15)", err)
	}
	<-done
	log.Println("Navigation completion notified")

	drone.Land()

	drone.ControlDisconnect()
	log.Println("Disconnected normally from Tello")
}
func TestAutoTurnToYaw(t *testing.T) {
	drone := new(Tello)

	err := drone.ControlConnectDefault()
	if err != nil {
		log.Fatalf("CCD failed with error %v", err)
	}
	log.Println("Connected to Tello control channel")
	time.Sleep(2 * time.Second)

	drone.TakeOff()
	time.Sleep(5 * time.Second)

	if _, err = drone.AutoTurnToYaw(40); err != nil { // should rotate +40deg
		t.Errorf("Error %v calling AutoTurnToYaw(40)", err)
	}
	time.Sleep(5 * time.Second)

	done, err := drone.AutoTurnToYaw(-90)
	if err != nil { // should rotate back 90 deg
		t.Errorf("Error %v calling AutoTurnToYaw(-90)", err)
	}
	<-done
	log.Println("Navigation completion notified")

	drone.Land()

	drone.ControlDisconnect()
	log.Println("Disconnected normally from Tello")
}

func TestAutoTurnToYawAndHeightConcurrently(t *testing.T) {
	drone := new(Tello)

	err := drone.ControlConnectDefault()
	if err != nil {
		log.Fatalf("Control Connect failed with error %v", err)
	}
	log.Println("Connected to Tello control channel")
	time.Sleep(2 * time.Second)

	drone.TakeOff()
	time.Sleep(5 * time.Second)

	hDoneC, err := drone.AutoFlyToHeight(4)
	if err != nil {
		log.Fatalf("AutoFlyToHeight failed with error %v", err)
	}
	yDoneC, err := drone.AutoTurnToYaw(120)
	if err != nil {
		log.Fatalf("AutoTurnToYaw failed with error %v", err)
	}

	var hDone, yDone bool
	for !hDone && !yDone {
		select {
		case hDone = <-hDoneC:
		case yDone = <-yDoneC:
		}
	}

	log.Println("Both manoeuvres complete")

	drone.Land()
	drone.ControlDisconnect()
	log.Println("Disconnected normally from Tello")
}
func TestAutoTurnByDeg(t *testing.T) {
	drone := new(Tello)

	err := drone.ControlConnectDefault()
	if err != nil {
		log.Fatalf("CCD failed with error %v", err)
	}
	log.Println("Connected to Tello control channel")
	time.Sleep(2 * time.Second)

	drone.TakeOff()
	time.Sleep(5 * time.Second)

	done, err := drone.AutoTurnByDeg(40)
	if err != nil { // should rotate +40deg
		t.Errorf("Error %v calling AutoTurnByDeg(40)", err)
	}
	<-done

	done, err = drone.AutoTurnByDeg(-80)
	if err != nil { // should rotate back to -40
		t.Errorf("Error %v calling AutoTurnByDeg(-80)", err)
	}
	<-done
	log.Println("Navigation completion notified")

	drone.Land()

	drone.ControlDisconnect()
	log.Println("Disconnected normally from Tello")
}
