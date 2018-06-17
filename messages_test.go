// tello project messages_test.go

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
	"bytes"
	"testing"
)

// use go test -count=1 to bypass test caching

func TestPacketToBuffer(t *testing.T) {
	// create a minimal packet
	var p packet

	p.header = msgHdr
	p.toDrone = true
	p.packetType = ptSet
	p.messageID = msgDoTakeoff
	p.sequence = 0

	b := packetToBuffer(p)

	correct := []byte{0xcc, 0x58, 0, 0x7c, 0x68, 0x54, 0, 0, 0, 0xb2, 0x89}

	if !bytes.Equal(correct, b) {
		t.Error("Buffer encoding incorrect")
	}
}

func TestByteToFloat32(t *testing.T) {
	var b = []byte{
		0, 0, 0, 0,
		128, 63, 0, 0, 112, 65,
	}
	var r float32
	if r = bytesToFloat32(b[0:5]); r != 0 {
		t.Errorf("Expected 0 got, %f\n", r)
	}
	if r = bytesToFloat32(b[2:7]); r != 1 {
		t.Errorf("Expected 1 got, %f\n", r)
	}
	if r = bytesToFloat32(b[6:]); r != 15 {
		t.Errorf("Expected 15 got, %f\n", r)
	}
}
