// flightCommands.go

// This file contains the high-level Tello flight command API

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

// TakeOff sends a normal takeoff request to the Tello
func (tello *Tello) TakeOff() {
	tello.ctrlMu.Lock()
	defer tello.ctrlMu.Unlock()
	// create the command packet
	var pkt packet

	// populate the command packet fields we need
	pkt.header = msgHdr
	pkt.toDrone = true
	pkt.packetType = ptSet
	pkt.messageID = msgDoTakeoff
	tello.ctrlSeq++
	pkt.sequence = tello.ctrlSeq

	// pack the packet into raw format and calculate CRCs etc.
	buff := packetToBuffer(pkt)

	// send the command packet
	tello.ctrlConn.Write(buff)
}

// ThrowTakeOff initiates a 'throw and go' launch
func (tello *Tello) ThrowTakeOff() {
	tello.ctrlMu.Lock()
	defer tello.ctrlMu.Unlock()
	// create the command packet
	var pkt packet

	// populate the command packet fields we need
	pkt.header = msgHdr
	pkt.toDrone = true
	pkt.packetType = ptGet
	pkt.messageID = msgDoThrowTakeoff
	tello.ctrlSeq++
	pkt.sequence = tello.ctrlSeq

	// pack the packet into raw format and calculate CRCs etc.
	buff := packetToBuffer(pkt)

	// send the command packet
	tello.ctrlConn.Write(buff)
}

func (tello *Tello) landit(palm bool) {
	tello.ctrlMu.Lock()
	defer tello.ctrlMu.Unlock()
	// create the command packet
	var pkt packet

	// populate the command packet fields we need
	pkt.header = msgHdr
	pkt.toDrone = true
	pkt.packetType = ptSet
	pkt.messageID = msgDoLand
	tello.ctrlSeq++
	pkt.sequence = tello.ctrlSeq
	pkt.payload = make([]byte, 1)
	if palm {
		pkt.payload[0] = 1
	} else {
		pkt.payload[0] = 0
	}

	// pack the packet into raw format and calculate CRCs etc.
	buff := packetToBuffer(pkt)

	// send the command packet
	tello.ctrlConn.Write(buff)
}

// Land sends a normal Land request to the Tello
func (tello *Tello) Land() {
	tello.landit(false)
}

// PalmLand initiates a Palm Landing
func (tello *Tello) PalmLand() {
	tello.landit(true)
}

// Bounce toggles the bouncing mode of the Tello
func (tello *Tello) Bounce() {
	tello.ctrlMu.Lock()
	defer tello.ctrlMu.Unlock()
	// create the command packet
	var pkt packet

	// populate the command packet fields we need
	pkt.header = msgHdr
	pkt.toDrone = true
	pkt.packetType = ptSet
	pkt.messageID = msgDoBounce
	tello.ctrlSeq++
	pkt.sequence = tello.ctrlSeq
	pkt.payload = make([]byte, 1)
	if tello.ctrlBouncing {
		pkt.payload[0] = 0x31
		tello.ctrlBouncing = false
	} else {
		pkt.payload[0] = 0x30
		tello.ctrlBouncing = true
	}

	// pack the packet into raw format and calculate CRCs etc.
	buff := packetToBuffer(pkt)

	// send the command packet
	tello.ctrlConn.Write(buff)
}

// Hover simply sets the sticks to zero - useful as a panic action!
func (tello *Tello) Hover() {
	tello.ctrlMu.Lock()
	tello.ctrlLx = 0
	tello.ctrlLy = 0
	tello.ctrlRx = 0
	tello.ctrlRy = 0
	tello.ctrlThrottle = 0
	tello.ctrlMu.Unlock()
}

// Forward tells the drone to start moving forward at a given speed between 0 and 100
func (tello *Tello) Forward(pct int) {
	var speed int16
	if pct > 0 {
		speed = int16(pct) * 327 // /100 * 32767
	}
	tello.UpdateSticks(StickMessage{Rx: 0, Ry: speed, Lx: 0, Ly: 0, Throttle: 0})
}

// Backward tells the drone to start moving Backward at a given speed between 0 and 100
func (tello *Tello) Backward(pct int) {
	var speed int16
	if pct > 0 {
		speed = int16(pct) * 327 // /100 * 32767
	}
	tello.UpdateSticks(StickMessage{Rx: 0, Ry: -speed, Lx: 0, Ly: 0, Throttle: 0})
}

// Left tells the drone to start moving Left at a given speed between 0 and 100
func (tello *Tello) Left(pct int) {
	var speed int16
	if pct > 0 {
		speed = int16(pct) * 327 // /100 * 32767
	}
	tello.UpdateSticks(StickMessage{Rx: -speed, Ry: 0, Lx: 0, Ly: 0, Throttle: 0})
}

// Right tells the drone to start moving Right at a given speed between 0 and 100
func (tello *Tello) Right(pct int) {
	var speed int16
	if pct > 0 {
		speed = int16(pct) * 327 // /100 * 32767
	}
	tello.UpdateSticks(StickMessage{Rx: speed, Ry: 0, Lx: 0, Ly: 0, Throttle: 0})
}

// Up tells the drone to start moving Up at a given speed between 0 and 100
func (tello *Tello) Up(pct int) {
	var speed int16
	if pct > 0 {
		speed = int16(pct) * 327 // /100 * 32767
	}
	tello.UpdateSticks(StickMessage{Rx: 0, Ry: 0, Lx: 0, Ly: speed, Throttle: 0})
}

// Down tells the drone to start moving Down at a given speed between 0 and 100
func (tello *Tello) Down(pct int) {
	var speed int16
	if pct > 0 {
		speed = int16(pct) * 327 // /100 * 32767
	}
	tello.UpdateSticks(StickMessage{Rx: 0, Ry: 0, Lx: 0, Ly: -speed, Throttle: 0})
}

// Clockwise tells the drone to start rotating Clockwise at a given speed between 0 and 100
func (tello *Tello) Clockwise(pct int) {
	var speed int16
	if pct > 0 {
		speed = int16(pct) * 327 // /100 * 32767
	}
	tello.UpdateSticks(StickMessage{Rx: 0, Ry: 0, Lx: speed, Ly: 0, Throttle: 0})
}

// TurnRight is an alias for Clockwise()
func (tello *Tello) TurnRight(pct int) {
	tello.Clockwise(pct)
}

// Anticlockwise tells the drone to start rotating Anticlockwise at a given speed between 0 and 100
func (tello *Tello) Anticlockwise(pct int) {
	var speed int16
	if pct > 0 {
		speed = int16(pct) * 327 // /100 * 32767
	}
	tello.UpdateSticks(StickMessage{Rx: 0, Ry: 0, Lx: -speed, Ly: 0, Throttle: 0})
}

// TurnLeft is an alias for Anticlockwise()
func (tello *Tello) TurnLeft(pct int) {
	tello.Anticlockwise(pct)
}

// CounterClockwise is an alias for Anticlockwise()
func (tello *Tello) CounterClockwise(pct int) {
	tello.Anticlockwise(pct)
}
