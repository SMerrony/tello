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

import (
	"math"
)

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
		recLen := int(uint8(data[pos+1])) + int(uint8(data[pos+2]))<<8
		logRecType := uint16(data[pos+4]) + uint16(data[pos+5])<<8
		//log.Printf("Flight Log - Rec type: %x, len:%d\n", logRecType, recLen)
		xorBuf := make([]byte, 256)
		xorVal := data[pos+6]
		switch logRecType {
		case logRecNewMVO:
			//log.Println("NewMOV rec found")
			for i := 0; i < recLen && pos+i < len(data); i++ {
				xorBuf[i] = data[pos+i] ^ xorVal
			}
			offset := 10
			flags := data[offset+76]
			tello.fdMu.Lock()
			if flags&logValidVelX != 0 {
				tello.fd.MVO.VelocityX = (int16(xorBuf[offset+2]) + int16(xorBuf[offset+3])<<8)
			}
			if flags&logValidVelY != 0 {
				tello.fd.MVO.VelocityY = (int16(xorBuf[offset+4]) + int16(xorBuf[offset+5])<<8)
			}
			if flags&logValidVelZ != 0 {
				tello.fd.MVO.VelocityZ = -(int16(xorBuf[offset+6]) + int16(xorBuf[offset+7])<<8)
			}
			if flags&logValidPosY != 0 {
				tello.fd.MVO.PositionY = bytesToFloat32(xorBuf[offset+8 : offset+13])
			}
			if flags&logValidPosX != 0 {
				tello.fd.MVO.PositionX = bytesToFloat32(xorBuf[offset+12 : offset+17])
			}
			if flags&logValidPosZ != 0 {
				tello.fd.MVO.PositionZ = bytesToFloat32(xorBuf[offset+16 : offset+21])
			}
			tello.fdMu.Unlock()
		case logRecIMU:
			//log.Println("IMU rec found")
			for i := 0; i < recLen && pos+i < len(data); i++ {
				xorBuf[i] = data[pos+i] ^ xorVal
			}
			offset := 10
			tello.fdMu.Lock()
			tello.fd.IMU.QuaternionW = bytesToFloat32(xorBuf[offset+48 : offset+53])
			tello.fd.IMU.QuaternionX = bytesToFloat32(xorBuf[offset+52 : offset+57])
			tello.fd.IMU.QuaternionY = bytesToFloat32(xorBuf[offset+56 : offset+61])
			tello.fd.IMU.QuaternionZ = bytesToFloat32(xorBuf[offset+60 : offset+65])
			tello.fd.IMU.Temperature = (int16(xorBuf[offset+106]) + int16(xorBuf[offset+107])<<8) / 100
			tello.fd.IMU.Yaw = quatToYawDeg(tello.fd.IMU.QuaternionX,
				tello.fd.IMU.QuaternionY,
				tello.fd.IMU.QuaternionZ,
				tello.fd.IMU.QuaternionW)
			tello.fdMu.Unlock()
		}
		pos += recLen
	}
}

// QuatToEulerDeg converts a quaternion set into pitch, roll & yaw expressed in degrees
func QuatToEulerDeg(qX, qY, qZ, qW float32) (pitch, roll, yaw int) {
	const degree = math.Pi / 180.0
	qqX := float64(qX)
	qqY := float64(qY)
	qqZ := float64(qZ)
	qqW := float64(qW)
	sqX := qqX * qqX
	sqY := qqY * qqY
	sqZ := qqZ * qqZ

	sinR := 2.0 * (qqW*qqX + qqY*qqZ)
	cosR := 1 - 2*(sqX+sqY)
	roll = int(math.Round(math.Atan2(sinR, cosR) / degree))

	sinP := 2.0 * (qqW*qqY - qqZ*qqX)
	if sinP > 1.0 {
		sinP = 1.0
	}
	if sinP < -1.0 {
		sinP = -1
	}
	pitch = int(math.Round(math.Asin(sinP) / degree))

	sinY := 2.0 * (qqW*qqZ + qqX*qqY)
	cosY := 1.0 - 2*(sqY+sqZ)
	yaw = int(math.Round(math.Atan2(sinY, cosY) / degree))

	return pitch, roll, yaw
}

// faster func just for getting yaw internally
func quatToYawDeg(qX, qY, qZ, qW float32) (yaw int16) {
	const degree = math.Pi / 180.0
	qqX := float64(qX)
	qqY := float64(qY)
	qqZ := float64(qZ)
	qqW := float64(qW)
	sqY := qqY * qqY
	sqZ := qqZ * qqZ

	sinY := 2.0 * (qqW*qqZ + qqX*qqY)
	cosY := 1.0 - 2*(sqY+sqZ)
	yaw = int16(math.Round(math.Atan2(sinY, cosY) / degree))

	return yaw
}
