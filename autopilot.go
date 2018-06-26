// autopilot.go

// This file contains Tello flight command API except for stick control.

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
	"errors"
	"log"
	"time"
)

const autopilotPeriodMs = 25 // how often the autopilot(s) monitor the drone

// CancelGotoHeight stops any in-flight GotoHeight navigation.
// The drone should stop moving vertically.
func (tello *Tello) CancelGotoHeight() {
	tello.autoHeightMu.Lock()
	tello.autoHeight = false
	tello.autoHeightMu.Unlock()
}

// GotoHeight starts vertical movement to the specified height in decimetres.
// The func returns immediately and a Goroutine handles the navigation.
// The caller may optionally listen on the 'done' channel for a signal that
// the navigation is complete (may have been cancelled).
func (tello *Tello) GotoHeight(dm int16) (done chan bool, err error) {
	//log.Printf("GotoHeight called with height: %d\n", dm)
	// are we already navigating?
	tello.autoHeightMu.RLock()
	if tello.autoHeight {
		tello.autoHeightMu.RUnlock()
		return nil, errors.New("Already navigating vertically")
	}
	tello.autoHeightMu.RUnlock()

	tello.autoHeightMu.Lock()
	tello.autoHeight = true
	tello.autoHeightMu.Unlock()

	done = make(chan bool, 1) // buffered so send doesn't block

	//log.Println("Autoheight set - starting goroutine")

	go func() {
		for {
			// has autoflight been cancelled?
			tello.autoHeightMu.RLock()
			cancelled := tello.autoHeight == false
			tello.autoHeightMu.RUnlock()
			if cancelled {
				log.Println("Cancelled")
				// stop vertical movement
				tello.ctrlMu.Lock()
				tello.ctrlLy = 0
				tello.ctrlMu.Unlock()
				tello.sendStickUpdate()
				done <- true
				return
			}

			tello.fdMu.RLock()
			delta := dm - tello.fd.Height // delta will be positive if we are too low
			//log.Printf("Target: %d, Height: %d, Delta: %d\n", dm, tello.fd.Height, delta)
			tello.fdMu.RUnlock()

			tello.ctrlMu.Lock()
			switch {
			case delta > 4:
				tello.ctrlLy = 32500 // full throttle if >40cm off target
			case delta > 0:
				tello.ctrlLy = 16250 // half throttle if <40cm off target
			case delta < -4:
				tello.ctrlLy = -32500
			case delta < 0:
				tello.ctrlLy = -16250
			case delta == 0: // might need some 'tolerance' here?
				// we're there! Cancel...
				tello.autoHeightMu.Lock()
				tello.autoHeight = false
				tello.autoHeightMu.Unlock()
			}
			tello.ctrlMu.Unlock()
			tello.sendStickUpdate()

			time.Sleep(autopilotPeriodMs * time.Millisecond)
		}
	}()

	return done, nil
}

// CancelGotoYaw stops any in-flight GotoYaw navigation.
// The drone should stop rotating.
func (tello *Tello) CancelGotoYaw() {
	tello.autoYawMu.Lock()
	tello.autoYaw = false
	tello.autoYawMu.Unlock()
}

// GotoYaw starts rotational movement to the specified yaw in degrees.
// The yaw should be between -180 and +180 degrees.
// The func returns immediately and a Goroutine handles the navigation.
func (tello *Tello) GotoYaw(targetYaw int16) (err error) {
	//log.Printf("GotoYaw called with height: %d\n", targetYaw)
	if targetYaw < -180 || targetYaw > 180 {
		return errors.New("Target yaw must be between -180 and +180")
	}
	// are we already navigating?
	tello.autoYawMu.RLock()
	if tello.autoYaw {
		tello.autoYawMu.RUnlock()
		return errors.New("Already navigating rotationally")
	}
	tello.autoYawMu.RUnlock()

	tello.autoYawMu.Lock()
	tello.autoYaw = true
	tello.autoYawMu.Unlock()

	// // should we restate the target yaw to avoid unneccessary rotation?
	// tello.fdMu.RLock()
	// currentYaw := tello.fd.IMU.Yaw
	// tello.dfMu.RUnlock()
	// if int16Abs(currentYaw - targetYaw) > 180

	//log.Println("autoYaw set - starting goroutine")

	go func() {
		for {
			// has autoflight been cancelled?
			tello.autoYawMu.RLock()
			cancelled := tello.autoYaw == false
			tello.autoYawMu.RUnlock()
			if cancelled {
				log.Println("Cancelled")
				// stop rotational movement
				tello.ctrlMu.Lock()
				tello.ctrlLx = 0
				tello.ctrlMu.Unlock()
				tello.sendStickUpdate()
				return
			}

			tello.fdMu.RLock()
			delta := targetYaw - tello.fd.IMU.Yaw // delta will be positive if we are too anticlockwise
			//log.Printf("Target: %d, Height: %d, Delta: %d\n", dm, tello.fd.Height, delta)
			tello.fdMu.RUnlock()

			tello.ctrlMu.Lock()
			switch {
			case delta > 4:
				tello.ctrlLx = 32500 // full throttle if >4deg off target
			case delta > 0:
				tello.ctrlLx = 16250 // half throttle if <4deg off target
			case delta < -4:
				tello.ctrlLx = -32500
			case delta < 0:
				tello.ctrlLx = -16250
			case delta == 0: // might need some 'tolerance' here?
				// we're there! Cancel...
				tello.autoYawMu.Lock()
				tello.autoYaw = false
				tello.autoYawMu.Unlock()
			}
			tello.ctrlMu.Unlock()
			tello.sendStickUpdate()

			time.Sleep(autopilotPeriodMs * time.Millisecond)
		}
	}()

	return nil
}

func int16Abs(x int16) int16 {
	if x < 0 {
		return -x
	}
	return x
}
