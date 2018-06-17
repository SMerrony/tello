// tello package flog.go - handle the flight logs from the drone

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

func (tello *Tello) ackLogHeader(id []byte) {
	tello.ctrlMu.Lock()
	defer tello.ctrlMu.Unlock()
	tello.ctrlSeq++
	pkt := newPacket(ptData1, msgLogHeader, tello.ctrlSeq, 3)
	pkt.payload[1] = id[0]
	pkt.payload[2] = id[1]
	tello.ctrlConn.Write(packetToBuffer(pkt))
}

func (tello *Tello) parseLogPacket(data []byte) {
	pos := 1
	if len(data) < 2 {
		return
	}
	for pos < len(data)-6 {
		if data[pos] != logRecordSeparator {
			//log.Println("Error parsing log record (bad separator)")
			break
		}
		recLen := int(uint8(data[pos+1]))
		if data[pos+2] != 0 {
			//log.Println("Error parsing log record (too long)")
			break
		}
		logRecType := uint16(data[pos+3]) + uint16(data[pos+4])<<8
		xorBuf := make([]byte, 256)
		xorVal := data[pos+6]
		switch logRecType {
		case logRecNewMVO:
			//log.Println("NewMOV rec found")
			for i := 0; i < 18; i++ {
				xorBuf[i] = data[pos+i] ^ xorVal
			}
			offset := 12
			tello.fdMu.Lock()
			tello.fd.VelocityX = int16(xorBuf[offset]) + int16(xorBuf[offset+1])<<8
			offset += 2
			tello.fd.VelocityY = int16(xorBuf[offset]) + int16(xorBuf[offset+1])<<8
			offset += 2
			tello.fd.VelocityZ = int16(xorBuf[offset]) + int16(xorBuf[offset+1])<<8
			offset += 2
			tello.fd.PositionX = bytesToFloat32(xorBuf[offset : offset+5])
			offset += 4
			tello.fd.PositionY = bytesToFloat32(xorBuf[offset : offset+5])
			offset += 4
			tello.fd.PositionZ = bytesToFloat32(xorBuf[offset : offset+5])
			tello.fdMu.Unlock()
			// log.Printf("Decoded log velocities %d, %d, %d\n",
			// 	tello.fd.VelocityX, tello.fd.VelocityY, tello.fd.VelocityZ)
			// log.Printf("Decoded log positions %f, %f, %f\n",
			// 	tello.fd.PositionX, tello.fd.PositionY, tello.fd.PositionZ)
		case logRecIMU:
			//log.Println("IMU rec found")
		}
		pos += recLen
	}
}
