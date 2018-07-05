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
	"math"
	"time"
)

const (
	autopilotPeriodMs   = 25 // how often the autopilot(s) monitor/redirect the drone
	autoPilotSpeedFast  = 32767
	autoPilotSpeedSlow  = 16384
	autoPilotSpeedVSlow = 8192
	// AutoHeightLimitDm is the maximum vertical displacement allowed for AutoFlyToHeight() etc. in decimetres.
	AutoHeightLimitDm = 300
	// AutoXYLimitM is the maximum horizontal displacement allowed for AutoFlyToXY() etc. in metres.
	AutoXYLimitM = 200 // 200m
	// AutoXYToleranceM is the maximum accuracy we try to attain in XY navigation in metres
	AutoXYToleranceM = 0.3
)

// CancelAutoFlyToHeight stops any in-flight AutoFlyToHeight navigation.
// The drone should stop moving vertically.
func (tello *Tello) CancelAutoFlyToHeight() {
	tello.autoHeightMu.Lock()
	tello.autoHeight = false
	tello.autoHeightMu.Unlock()
}

// AutoFlyToHeight starts vertical movement to the specified height in decimetres
// (so a value of 10 means 1m).
// The func returns immediately and a Goroutine handles the navigation until either
// it is complete or cancelled via CancelAutoFlyToHeight().
// The caller may optionally listen on the 'done' channel for a signal that
// the navigation is complete (or has been cancelled).
func (tello *Tello) AutoFlyToHeight(dm int16) (done chan bool, err error) {
	//log.Printf("AutoFlyToHeight called with height: %d\n", dm)
	if dm > AutoHeightLimitDm || dm < -AutoHeightLimitDm {
		return nil, errors.New("Verical navigation limit exceeded")
	}
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
				//log.Println("Cancelled")
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
				tello.ctrlLy = autoPilotSpeedFast // full throttle if >40cm off target
			case delta > 0:
				tello.ctrlLy = autoPilotSpeedSlow // half throttle if <40cm off target
			case delta < -4:
				tello.ctrlLy = -autoPilotSpeedFast
			case delta < 0:
				tello.ctrlLy = -autoPilotSpeedSlow
			case delta == 0: // might need some 'tolerance' here?
				// we're there! Cancel...
				tello.autoHeightMu.Lock()
				tello.autoHeight = false
				tello.autoHeightMu.Unlock()
			}
			tello.ctrlMu.Unlock()
			//tello.sendStickUpdate()

			time.Sleep(autopilotPeriodMs * time.Millisecond)
		}
	}()

	return done, nil
}

// CancelAutoTurn stops any in-flight AutoTurnToYaw or AutoTurnByDeg navigation.
// The drone should stop rotating.
func (tello *Tello) CancelAutoTurn() {
	tello.autoYawMu.Lock()
	tello.autoYaw = false
	tello.autoYawMu.Unlock()
}

// AutoTurnToYaw starts rotational movement to the specified yaw in degrees.
// The yaw should be between -180 and +180 degrees.
// The func returns immediately and a Goroutine handles the navigation.
// The caller may optionally listen on the 'done' channel for a signal that
// the navigation is complete (may have been cancelled).
// You may explicitly cancel this operation via CancelAutoTurn().
func (tello *Tello) AutoTurnToYaw(targetYaw int16) (done chan bool, err error) {
	//log.Printf("AutoTurnToYaw called with target: %d\n", targetYaw)
	if targetYaw < -180 || targetYaw > 180 {
		return nil, errors.New("Target yaw must be between -180 and +180")
	}
	adjustedTarget := targetYaw
	if targetYaw < 0 {
		adjustedTarget = 360 + targetYaw
	}

	// are we already navigating?
	tello.autoYawMu.RLock()
	if tello.autoYaw {
		tello.autoYawMu.RUnlock()
		return nil, errors.New("Already navigating rotationally")
	}
	tello.autoYawMu.RUnlock()

	tello.autoYawMu.Lock()
	tello.autoYaw = true
	tello.autoYawMu.Unlock()

	done = make(chan bool, 1) // buffered so send doesn't block

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
				done <- true
				return
			}

			tello.fdMu.RLock()
			adjustedCurrent := tello.fd.IMU.Yaw
			tello.fdMu.RUnlock()
			if adjustedCurrent < 0 {
				adjustedCurrent = 360 + adjustedCurrent
			}

			delta := adjustedTarget - adjustedCurrent
			absDelta := int16Abs(delta)
			switch {
			case absDelta <= 180: //
			case delta > 0:
				delta = absDelta - 360
			case delta > 0:
				delta = 360 - absDelta
			}

			//log.Printf("Target: %d, Current: %d, Delta: %d\n", adjustedTarget, adjustedCurrent, delta)

			tello.ctrlMu.Lock()
			switch {
			case delta > 10:
				tello.ctrlLx = autoPilotSpeedFast
			case delta > 0:
				tello.ctrlLx = autoPilotSpeedSlow
			case delta < -10:
				tello.ctrlLx = -autoPilotSpeedFast
			case delta < 0:
				tello.ctrlLx = -autoPilotSpeedSlow
			case delta == 0: // might need some 'tolerance' here?
				// we're there! Cancel...
				tello.autoYawMu.Lock()
				tello.autoYaw = false
				tello.autoYawMu.Unlock()
			}
			tello.ctrlMu.Unlock()
			//tello.sendStickUpdate()

			time.Sleep(autopilotPeriodMs * time.Millisecond)
		}
	}()

	return done, nil
}

// AutoTurnByDeg starts rotational movement by the specified amount degrees.
// The amount should be between -180 and +180 degrees, where negative values cause an
// anticlockwise rotation and vice-versa.
// The func returns immediately and a Goroutine handles the navigation.
// The caller may optionally listen on the 'done' channel for a signal that
// the navigation is complete (may have been cancelled).
// You may explicitly cancel this operation via CancelAutoTurn().
func (tello *Tello) AutoTurnByDeg(delta int16) (done chan bool, err error) {

	if delta < -180 || delta > 180 {
		return nil, errors.New("Turn amount must be between -180 and +180")
	}

	// are we already navigating?
	tello.autoYawMu.RLock()
	if tello.autoYaw {
		tello.autoYawMu.RUnlock()
		return nil, errors.New("Already navigating rotationally")
	}
	tello.autoYawMu.RUnlock()

	tello.fdMu.RLock()
	adjustedTarget := tello.fd.IMU.Yaw
	tello.fdMu.RUnlock()

	adjustedTarget += delta
	switch {
	case adjustedTarget == 360:
		adjustedTarget = 0
	case adjustedTarget > 180:
		adjustedTarget = -(360 - adjustedTarget)
	case adjustedTarget < -180:
		adjustedTarget = 360 + adjustedTarget
	}

	return tello.AutoTurnToYaw(adjustedTarget)
}

// // autoWaitAndSetOrigin is run as a Goroutine after takeoff is initiated.
// func (tello *Tello) autoWaitAndSetOrigin() {

// }

// SetHome establishes the current MVO position and IMU yaw as the home
// point for autopilot operations.  It could be called after takeoff to establish a
// home coordinate, or during (non-autopilot) flight to set a waypoint.
func (tello *Tello) SetHome() (err error) {
	tello.autoXYMu.RLock()
	alreadyAuto := tello.autoXY
	tello.autoXYMu.RUnlock()
	if alreadyAuto {
		return errors.New("Cannot set origin during automatic flight")
	}
	tello.autoXYMu.Lock()
	tello.autoXY = false
	tello.fdMu.RLock()
	tello.homeX = tello.fd.MVO.PositionX
	tello.homeY = tello.fd.MVO.PositionY
	tello.homeYaw = tello.fd.IMU.Yaw
	tello.fdMu.RUnlock()
	if tello.homeYaw < 0 {
		tello.homeYaw += 360
	}
	tello.homeValid = true
	tello.autoXYMu.Unlock()
	return nil
}

// IsHomeSet tests whether the home point used for the travelling AutoFly... funcs is set.
func (tello *Tello) IsHomeSet() (set bool) {
	tello.autoXYMu.RLock()
	set = tello.homeValid
	tello.autoXYMu.RUnlock()
	return set
}

// CancelAutoFlyToXY stops any in-flight AutoFlyToXY navigation.
// The drone should stop.
func (tello *Tello) CancelAutoFlyToXY() {
	tello.autoXYMu.Lock()
	tello.autoXY = false
	tello.autoXYMu.Unlock()
}

// AutoFlyToXY starts horizontal movement to the specified (X, Y) location
// expressed in metres from the home point (which must have been previously set).
// The func returns immediately and a Goroutine handles the navigation until either
// it is complete or cancelled via CancelFlyToXY().
// The caller may optionally listen on the 'done' channel for a signal that
// the navigation is complete (or has been cancelled).
func (tello *Tello) AutoFlyToXY(targetX, targetY float32) (done chan bool, err error) {
	//log.Printf("FlyToXY called with XY: %d\n", dm)
	if targetX > AutoXYLimitM || targetY > AutoXYLimitM ||
		targetX < -AutoXYLimitM || targetY < -AutoXYLimitM {
		return nil, errors.New("Horizontal navigation limit exceeded")
	}
	// are we already navigating?
	tello.autoXYMu.RLock()
	already := tello.autoXY
	valid := tello.homeValid
	originX := tello.homeX
	originY := tello.homeY
	tello.autoXYMu.RUnlock()
	if already {
		return nil, errors.New("Already AutoFlying horizontally")
	}
	if !valid {
		return nil, errors.New("Cannot AutoFly as home point has not be set (or is invalid)")
	}

	tello.autoXYMu.Lock()
	tello.autoXY = true
	tello.autoXYMu.Unlock()

	// adjust target relative to origin -SHOULD WE ADJUST YAW TOO???
	targetX += originX
	targetY += originY

	done = make(chan bool, 1) // buffered so send doesn't block

	//log.Println("AutoXY set - starting goroutine")

	go func() {
		for {
			// has autoflight been cancelled?
			tello.autoXYMu.RLock()
			cancelled := tello.autoXY == false
			tello.autoXYMu.RUnlock()
			if cancelled {
				// stop vertical movement
				tello.ctrlMu.Lock()
				tello.ctrlRx = 0
				tello.ctrlRy = 0
				tello.ctrlMu.Unlock()
				tello.sendStickUpdate()
				done <- true
				return
			}

			// get current yaw & position
			tello.fdMu.RLock()
			currentYaw := tello.fd.IMU.Yaw
			currentX := tello.fd.MVO.PositionX
			currentY := tello.fd.MVO.PositionY
			tello.fdMu.RUnlock()

			deltaX, deltaY := calcXYdeltas(currentYaw, currentX, currentY, targetX, targetY)

			tello.ctrlMu.Lock()

			switch {
			case deltaX <= AutoXYToleranceM && deltaX >= -AutoXYToleranceM:
				tello.ctrlRx = 0
			case deltaX >= 1.5:
				tello.ctrlRx = autoPilotSpeedSlow // full throttle if =>1m off target
			case deltaX <= -1.5:
				tello.ctrlRx = -autoPilotSpeedSlow // full throttle if =>1m off target
			case deltaX > AutoXYToleranceM:
				tello.ctrlRx = autoPilotSpeedSlow // half throttle
			case deltaX < -AutoXYToleranceM:
				tello.ctrlRx = -autoPilotSpeedSlow // half throttle
			default:
				log.Fatalf("Invalid state in AutoFlyToXY() - deltaX=%f", deltaX)
			}
			switch {
			case deltaY <= AutoXYToleranceM && deltaY >= -AutoXYToleranceM:
				tello.ctrlRy = 0
			case deltaY >= 1.5:
				tello.ctrlRy = autoPilotSpeedSlow // full throttle if =>1m off target
			case deltaY <= -1.5:
				tello.ctrlRy = -autoPilotSpeedSlow // full throttle if =>1m off target
			case deltaY > AutoXYToleranceM:
				tello.ctrlRy = autoPilotSpeedSlow // half throttle
			case deltaY < -AutoXYToleranceM:
				tello.ctrlRy = -autoPilotSpeedSlow // half throttle
			default:
				log.Fatalf("Invalid state in AutoFlyToXY() - deltaY=%f", deltaY)
			}

			// log.Printf("Current %.2f,%.2f Yaw: %d - Target: %.2f,%.2f - Deltas X: %.2f, Y:%.2f - Throttles: %d,%d\n",
			// 	currentX, currentY, currentYaw, targetX, targetY, deltaX, deltaY, tello.ctrlRx, tello.ctrlRy)

			if tello.ctrlRx == 0.0 && tello.ctrlRy == 0.0 {
				// we're there! Cancel...
				tello.autoXYMu.Lock()
				tello.autoXY = false
				tello.autoXYMu.Unlock()
			}
			tello.ctrlMu.Unlock()
			//tello.sendStickUpdate()

			time.Sleep(autopilotPeriodMs * time.Millisecond)
		}
	}()

	return done, nil
}

func calcXYdeltas(yawDeg int16, currX, currY, targetX, targetY float32) (dx, dy float32) {
	adjustedYaw := float64(yawDeg)
	if adjustedYaw < 0 {
		adjustedYaw += 360.0
	}
	adjustedYaw *= math.Pi / 180

	dx = float32(math.Cos(adjustedYaw))*(targetX-currX) - float32(math.Sin(adjustedYaw))*(targetY-currY)
	dy = float32(math.Sin(adjustedYaw))*(targetX-currX) + float32(math.Cos(adjustedYaw))*(targetY-currY)

	return dx, dy
}

// Helper functions...
func int16Abs(x int16) int16 {
	if x < 0 {
		return -x
	}
	return x
}
