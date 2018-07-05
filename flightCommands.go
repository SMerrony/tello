// flightCommands.go

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

// TakeOff sends a normal takeoff request to the Tello.
// Any previously set origin is invalidated.
func (tello *Tello) TakeOff() {
	tello.ctrlMu.Lock()

	tello.autoXYMu.Lock()
	tello.homeValid = false // origin is invalidated until flying and reset
	tello.autoXYMu.Unlock()

	tello.ctrlSeq++
	pkt := newPacket(ptSet, msgDoTakeoff, tello.ctrlSeq, 0)
	tello.ctrlConn.Write(packetToBuffer(pkt))

	tello.ctrlMu.Unlock()
}

// ThrowTakeOff initiates a 'throw and go' launch.
// Any previously set origin is invalidated.
func (tello *Tello) ThrowTakeOff() {
	tello.ctrlMu.Lock()

	tello.autoXYMu.Lock()
	tello.homeValid = false // origin is invalidated until flying and reset
	tello.autoXYMu.Unlock()

	tello.ctrlSeq++
	pkt := newPacket(ptGet, msgDoThrowTakeoff, tello.ctrlSeq, 0)
	tello.ctrlConn.Write(packetToBuffer(pkt))

	tello.ctrlMu.Unlock()
}

// Land sends a normal Land request to the Tello.
func (tello *Tello) Land() {
	tello.ctrlMu.Lock()
	defer tello.ctrlMu.Unlock()

	tello.ctrlSeq++
	pkt := newPacket(ptSet, msgDoLand, tello.ctrlSeq, 1)
	pkt.payload[0] = 0 // see StopLanding() for use of this field
	tello.ctrlConn.Write(packetToBuffer(pkt))
}

// StopLanding cancels a land command.
func (tello *Tello) StopLanding() {
	tello.ctrlMu.Lock()
	defer tello.ctrlMu.Unlock()

	tello.ctrlSeq++
	pkt := newPacket(ptSet, msgDoLand, tello.ctrlSeq, 1)
	pkt.payload[0] = 1
	tello.ctrlConn.Write(packetToBuffer(pkt))
}

// PalmLand initiates a Palm Landing.
func (tello *Tello) PalmLand() {
	tello.ctrlMu.Lock()
	defer tello.ctrlMu.Unlock()

	tello.ctrlSeq++
	pkt := newPacket(ptSet, msgDoPalmLand, tello.ctrlSeq, 1)
	pkt.payload[0] = 0
	tello.ctrlConn.Write(packetToBuffer(pkt))
}

// Bounce toggles the bouncing mode of the Tello.
func (tello *Tello) Bounce() {
	tello.ctrlMu.Lock()
	defer tello.ctrlMu.Unlock()

	tello.ctrlSeq++
	pkt := newPacket(ptSet, msgDoBounce, tello.ctrlSeq, 1)
	if tello.ctrlBouncing {
		pkt.payload[0] = 0x31
		tello.ctrlBouncing = false
	} else {
		pkt.payload[0] = 0x30
		tello.ctrlBouncing = true
	}
	tello.ctrlConn.Write(packetToBuffer(pkt))
}

// Flip sends a flip flight command to the Tello.
func (tello *Tello) Flip(dir FlipType) {
	tello.ctrlMu.Lock()
	defer tello.ctrlMu.Unlock()

	tello.ctrlSeq++
	pkt := newPacket(ptFlip, msgDoFlip, tello.ctrlSeq, 1)
	pkt.payload[0] = byte(dir)
	tello.ctrlConn.Write(packetToBuffer(pkt))
}

// StartSmartVideo begins a preprogrammed 'smart video' flight action.
func (tello *Tello) StartSmartVideo(cmd SvCmd) {
	tello.ctrlMu.Lock()
	defer tello.ctrlMu.Unlock()

	tello.ctrlSeq++
	pkt := newPacket(ptSet, msgDoSmartVideo, tello.ctrlSeq, 1)
	pkt.payload[0] = byte(cmd) | 0x01
	tello.ctrlConn.Write(packetToBuffer(pkt))
}

// StopSmartVideo begins a preprogrammed 'smart video' flight action.
func (tello *Tello) StopSmartVideo(cmd SvCmd) {
	tello.ctrlMu.Lock()
	defer tello.ctrlMu.Unlock()

	tello.ctrlSeq++
	pkt := newPacket(ptSet, msgDoSmartVideo, tello.ctrlSeq, 1)
	pkt.payload[0] = byte(cmd)
	tello.ctrlConn.Write(packetToBuffer(pkt))
}

// *** The following are 'macro' commands which are here purely
// *** to make the Tello easier to use in some circumstances.

// Hover simply sets the sticks to zero which should halt all motion - useful as a panic action!
func (tello *Tello) Hover() {
	tello.ctrlMu.Lock()
	tello.ctrlLx = 0
	tello.ctrlLy = 0
	tello.ctrlRx = 0
	tello.ctrlRy = 0
	tello.ctrlMu.Unlock()
}

// Forward tells the drone to start moving forward at a given speed between 0 and 100.
func (tello *Tello) Forward(pct int) {
	var speed int16
	if pct > 0 {
		speed = int16(pct) * 327 // /100 * 32767
	}
	tello.UpdateSticks(StickMessage{Rx: 0, Ry: speed, Lx: 0, Ly: 0})
}

// Backward tells the drone to start moving Backward at a given speed between 0 and 100.
func (tello *Tello) Backward(pct int) {
	var speed int16
	if pct > 0 {
		speed = int16(pct) * 327 // /100 * 32767
	}
	tello.UpdateSticks(StickMessage{Rx: 0, Ry: -speed, Lx: 0, Ly: 0})
}

// Left tells the drone to start moving Left at a given speed between 0 and 100.
func (tello *Tello) Left(pct int) {
	var speed int16
	if pct > 0 {
		speed = int16(pct) * 327 // /100 * 32767
	}
	tello.UpdateSticks(StickMessage{Rx: -speed, Ry: 0, Lx: 0, Ly: 0})
}

// Right tells the drone to start moving Right at a given speed between 0 and 100.
func (tello *Tello) Right(pct int) {
	var speed int16
	if pct > 0 {
		speed = int16(pct) * 327 // /100 * 32767
	}
	tello.UpdateSticks(StickMessage{Rx: speed, Ry: 0, Lx: 0, Ly: 0})
}

// Up tells the drone to start moving Up at a given speed between 0 and 100.
func (tello *Tello) Up(pct int) {
	var speed int16
	if pct > 0 {
		speed = int16(pct) * 327 // /100 * 32767
	}
	tello.UpdateSticks(StickMessage{Rx: 0, Ry: 0, Lx: 0, Ly: speed})
}

// Down tells the drone to start moving Down at a given speed between 0 and 100.
func (tello *Tello) Down(pct int) {
	var speed int16
	if pct > 0 {
		speed = int16(pct) * 327 // /100 * 32767
	}
	tello.UpdateSticks(StickMessage{Rx: 0, Ry: 0, Lx: 0, Ly: -speed})
}

// Clockwise tells the drone to start rotating Clockwise at a given speed between 0 and 100.
func (tello *Tello) Clockwise(pct int) {
	var speed int16
	if pct > 0 {
		speed = int16(pct) * 327 // /100 * 32767
	}
	tello.UpdateSticks(StickMessage{Rx: 0, Ry: 0, Lx: speed, Ly: 0})
}

// TurnRight is an alias for Clockwise().
func (tello *Tello) TurnRight(pct int) {
	tello.Clockwise(pct)
}

// Anticlockwise tells the drone to start rotating Anticlockwise at a given speed between 0 and 100.
func (tello *Tello) Anticlockwise(pct int) {
	var speed int16
	if pct > 0 {
		speed = int16(pct) * 327 // /100 * 32767
	}
	tello.UpdateSticks(StickMessage{Rx: 0, Ry: 0, Lx: -speed, Ly: 0})
}

// TurnLeft is an alias for Anticlockwise().
func (tello *Tello) TurnLeft(pct int) {
	tello.Anticlockwise(pct)
}

// CounterClockwise is an alias for Anticlockwise().
func (tello *Tello) CounterClockwise(pct int) {
	tello.Anticlockwise(pct)
}

// SetSportsMode sets the sports mode of flight to the given value.
func (tello *Tello) SetSportsMode(sports bool) {
	tello.ctrlMu.Lock()
	tello.ctrlSportsMode = sports
	tello.ctrlMu.Unlock()
}

// SetFastMode sets the 'fast' or 'sports' mode of flight.
func (tello *Tello) SetFastMode() {
	tello.SetSportsMode(true)
}

// SetSlowMode sets the 'slow' or 'normal' mode of flight.
func (tello *Tello) SetSlowMode() {
	tello.SetSportsMode(false)
}

// Flips...

// BackFlip - flip backwards.
func (tello *Tello) BackFlip() { tello.Flip(FlipBackward) }

// BackLeftFlip - flip backwards and to the left.
func (tello *Tello) BackLeftFlip() { tello.Flip(FlipBackwardLeft) }

// BackRightFlip - flip backwards and to the right.
func (tello *Tello) BackRightFlip() { tello.Flip(FlipBackwardRight) }

// ForwardFlip - flip forwards.
func (tello *Tello) ForwardFlip() { tello.Flip(FlipForward) }

// ForwardRightFlip - flip forwardsand to the right.
func (tello *Tello) ForwardRightFlip() { tello.Flip(FlipForwardRight) }

// ForwardLeftFlip - flip forward and to the left.
func (tello *Tello) ForwardLeftFlip() { tello.Flip(FlipForwardLeft) }

// LeftFlip - flip to the left.
func (tello *Tello) LeftFlip() { tello.Flip(FlipLeft) }

// RightFlip - flip to the right.
func (tello *Tello) RightFlip() { tello.Flip(FlipRight) }

// *** End of 'macro' commands ***
