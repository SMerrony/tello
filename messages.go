// messages.go

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
	"encoding/binary"
	"math"
)

const msgHdr = 0xcc // 204

// packet is our internal representation of the messages passed to/from the Tello
type packet struct {
	header        byte
	size13        uint16
	crc8          byte
	fromDrone     bool // the following 4 fields are encoded in a single byte in the raw packet
	toDrone       bool
	packetType    uint8 // 3-bit
	packetSubtype uint8 // 3-bit
	messageID     uint16
	sequence      uint16
	payload       []byte
	crc16         uint16
}

const minPktSize = 11 // smallest possible raw packet

// tello packet types, 3 and 7 currently unknown
const (
	ptExtended = 0
	ptGet      = 1
	ptData1    = 2
	ptData2    = 4
	ptSet      = 5
	ptFlip     = 6
)

// Tello message IDs
const (
	msgDoConnect           = 0x0001 // 1
	msgConnected           = 0x0002 // 2
	msgQuerySSID           = 0x0011 // 17
	msgSetSSID             = 0x0012 // 18
	msgQuerySSIDPass       = 0x0013 // 19
	msgSetSSIDPass         = 0x0014 // 20
	msgQueryWifiRegion     = 0x0015 // 21
	msgSetWifiRegion       = 0x0016 // 22
	msgWifiStrength        = 0x001a // 26
	msgSetVideoBitrate     = 0x0020 // 32
	msgSetDynAdjRate       = 0x0021 // 33
	msgEisSetting          = 0x0024 // 36
	msgQueryVideoSPSPPS    = 0x0025 // 37
	msgQueryVideoBitrate   = 0x0028 // 40
	msgDoTakePic           = 0x0030 // 48
	msgSwitchPicVideo      = 0x0031 // 49
	msgDoStartRec          = 0x0032 // 50
	msgExposureVals        = 0x0034 // 52 (Get or set?)
	msgLightStrength       = 0x0035 // 53
	msgQueryJPEGQuality    = 0x0037 // 55
	msgError1              = 0x0043 // 67
	msgError2              = 0x0044 // 68
	msgQueryVersion        = 0x0045 // 69
	msgSetDateTime         = 0x0046 // 70
	msgQueryActivationTime = 0x0047 // 71
	msgQueryLoaderVersion  = 0x0049 // 73
	msgSetStick            = 0x0050 // 80
	msgDoTakeoff           = 0x0054 // 84
	msgDoLand              = 0x0055 // 85
	msgFlightStatus        = 0x0056 // 86
	msgSetHeightLimit      = 0x0058 // 88
	msgDoFlip              = 0x005c // 92
	msgDoThrowTakeoff      = 0x005d // 93
	msgDoPalmLand          = 0x005e // 94
	msgFileSize            = 0x0062 // 98
	msgFileData            = 0x0063 // 99
	msgFileDone            = 0x0064 // 100
	msgDoSmartVideo        = 0x0080 // 128
	msgSmartVideoStatus    = 0x0081 // 129
	msgLogHeader           = 0x1050 // 4176
	msgLogData             = 0x1051 // 4177
	msgLogConfig           = 0x1052 // 4178
	msgDoBounce            = 0x1053 // 4179
	msgDoCalibration       = 0x1054 // 4180
	msgSetLowBattThresh    = 0x1055 // 4181
	msgQueryHeightLimit    = 0x1056 // 4182
	msgQueryLowBattThresh  = 0x1057 // 4183
	msgSetAttitude         = 0x1058 // 4184
	msgQueryAttitude       = 0x1059 // 4185
)

// FlipType represents a flip direction.
type FlipType int

// Flip types...
const (
	FlipForward FlipType = iota
	FlipLeft
	FlipBackward
	FlipRight
	FlipForwardLeft
	FlipBackwardLeft
	FlipBackwardRight
	FlipForwardRight
)

// SvCmd is Smart Video flight command.
type SvCmd byte

// Smart Video flight commands...
const (
	Sv360    SvCmd = 1 << 2 // Slowly rotate around 360 degrees.
	SvCircle       = 2 << 2 // Circle around a point in front of the drone.
	SvUpOut        = 3 << 2 // Perform the 'Up and Out' manouvre.
)

// VBR is a Video Bit Rate, the int value is meaningless.
type VBR byte

// VBR settings...
const (
	VbrAuto VBR = iota // let the Tello choose the best for the current connection
	Vbr1M              // Set the VBR to 1Mbps
	Vbr1M5             // Set the VBR to 1.5Mbps
	Vbr2M              // Set the VBR to 2Mbps
	Vbr3M              // Set the VBR to 3Mbps
	Vbr4M              // Set the VBR to 4mbps
)

const (
	vmNormal = 0
	vmWide   = 1
)

// fileType is the type of file being sent to/from the drone
type fileType byte

// Known File Types...
const (
	ftJPEG fileType = 1
)

type fileData struct {
	fileType  fileType // 1 = JPEG
	fileSize  int
	fileBytes []byte
}

type fileInternal struct {
	fID          uint16
	filetype     fileType
	expectedSize int
	accumSize    int
	pieces       []filePiece
}

type filePiece struct {
	fID       uint16
	numChunks int
	chunks    []fileChunk
}

type fileChunk struct {
	fID       uint16
	pieceNum  uint32
	chunkNum  uint32
	chunkLen  uint16
	chunkData []byte
}

// FlightData holds our current knowledge of the drone's state.
// This data is not all sent at once from the drone, different fields may be updated
// at varying rates.
type FlightData struct {
	BatteryLow               bool
	BatteryCritical          bool
	BatteryMilliVolts        int16
	BatteryPercentage        int8
	BatteryState             bool
	CameraState              uint8
	DownVisualState          bool
	DroneFlyTimeLeft         int16
	DroneHover               bool
	EmOpen                   bool
	EastSpeed                int16
	ElectricalMachineryState uint8
	FactoryMode              bool
	Flying                   bool
	FlyMode                  uint8
	FlyTime                  int16
	FrontIn                  bool
	FrontLSC                 bool
	FrontOut                 bool
	GravityState             bool
	GroundSpeed              int16
	Height                   int16 // seems to be in decimetres
	IMU                      IMUData
	ImuCalibrationState      int8
	ImuState                 bool
	LightStrength            uint8
	LowBatteryThreshold      uint8
	MaxHeight                uint8
	MVO                      MVOData
	NorthSpeed               int16
	OnGround                 bool
	OutageRecording          bool
	OverTemp                 bool
	PowerState               bool
	PressureState            bool
	SmartVideoExitMode       int16
	SSID                     string
	ThrowFlyTimer            int8
	VerticalSpeed            int16
	Version                  string
	VideoBitrate             VBR
	WifiInterference         uint8
	WifiStrength             uint8
	WindState                bool
}

// MVOData comes from the flight log messages
type MVOData struct {
	PositionX, PositionY, PositionZ float32
	VelocityX, VelocityY, VelocityZ int16
}

// IMUData comes from the flight log messages
type IMUData struct {
	QuaternionW,
	QuaternionX, QuaternionY, QuaternionZ float32
	Temperature int16
	Yaw         int16 // derived from Quat fields, -180 > degrees > +180
}

// StickMessage holds the signed 16-bit values of a joystick update.
// Each value can range from -32768 to 32767
type StickMessage struct {
	Rx, Ry, Lx, Ly int16
}

const logRecordSeparator = 'U'

// flight log message IDs
const (
	logRecNewMVO = 0x001d
	logRecIMU    = 0x0800
	// TODO - there are many more
)

const (
	logValidVelX = 0x01
	logValidVelY = 0x02
	logValidVelZ = 0x04
	logValidPosX = 0x10
	logValidPosY = 0x20
	logValidPosZ = 0x40
)

// utility funcs for message handling

// bufferToPacket takes a raw buffer of bytes and populates our packet struct
func bufferToPacket(buff []byte) (pkt packet) {
	pkt.header = buff[0]
	pkt.size13 = (uint16(buff[1]) + uint16(buff[2])<<8) >> 3
	pkt.crc8 = buff[3]
	pkt.fromDrone = (buff[4] & 0x80) == 1
	pkt.toDrone = (buff[4] & 0x40) == 1
	pkt.packetType = uint8((buff[4] >> 3) & 0x07)
	pkt.packetSubtype = uint8(buff[4] & 0x07)
	pkt.messageID = (uint16(buff[6]) << 8) | uint16(buff[5])
	pkt.sequence = (uint16(buff[8]) << 8) | uint16(buff[7])
	payloadSize := pkt.size13 - 11
	if payloadSize > 0 {
		pkt.payload = make([]byte, payloadSize)
		copy(pkt.payload, buff[9:9+payloadSize])
	}
	pkt.crc16 = uint16(buff[pkt.size13-1])<<8 + uint16(buff[pkt.size13-2])
	return pkt
}

// newPacket returns a packet with some fields populated
func newPacket(pt uint8, cmd uint16, seq uint16, payloadSize int) (pkt packet) {
	pkt.header = msgHdr
	pkt.toDrone = true
	pkt.packetType = pt
	pkt.messageID = cmd
	pkt.sequence = seq
	if payloadSize > 0 {
		pkt.payload = make([]byte, payloadSize)
	}
	return pkt
}

// pack the packet into raw buffer format and calculate CRCs etc.
func packetToBuffer(pkt packet) (buff []byte) {
	// create a buffer of the right size
	payloadSize := len(pkt.payload)
	packetSize := minPktSize + payloadSize
	buff = make([]byte, packetSize)

	// copy each field, manipulating if necessary
	buff[0] = pkt.header
	buff[1] = byte(packetSize << 3)
	buff[2] = byte(packetSize >> 5)
	buff[3] = calculateCRC8(buff[0:3])
	buff[4] = pkt.packetSubtype + (pkt.packetType << 3)
	if pkt.toDrone {
		buff[4] |= 0x40
	}
	if pkt.fromDrone {
		buff[4] |= 0x80
	}
	buff[5] = byte(pkt.messageID)
	buff[6] = byte(pkt.messageID >> 8)
	buff[7] = byte(pkt.sequence)
	buff[8] = byte(pkt.sequence >> 8)

	for p := 0; p < payloadSize; p++ {
		buff[9+p] = pkt.payload[p]
	}
	crc16 := calculateCRC16(buff[0 : 9+payloadSize])
	buff[9+payloadSize] = byte(crc16)
	buff[10+payloadSize] = byte(crc16 >> 8)

	return buff
}

func payloadToFlightData(pl []byte) (fd FlightData) {
	fd.Height = int16(pl[0]) + int16(pl[1])<<8
	fd.NorthSpeed = int16(uint16(pl[2]) | uint16(pl[3])<<8)
	fd.EastSpeed = int16(pl[4]) | int16(pl[5])<<8
	fd.VerticalSpeed = int16(pl[6]) | int16(pl[7])<<8
	fd.FlyTime = int16(pl[8]) | int16(pl[9])<<8

	fd.ImuState = (pl[10] & 1) == 1
	fd.PressureState = (pl[10] >> 1 & 1) == 1
	fd.DownVisualState = (pl[10] >> 2 & 1) == 1
	fd.PowerState = (pl[10] >> 3 & 1) == 1
	fd.BatteryState = (pl[10] >> 4 & 1) == 1
	fd.GravityState = (pl[10] >> 5 & 1) == 1
	// what is bit 6?
	fd.WindState = (pl[10] >> 7 & 1) == 1

	fd.ImuCalibrationState = int8(pl[11])
	fd.BatteryPercentage = int8(pl[12])
	fd.DroneFlyTimeLeft = int16(pl[13]) + int16(pl[14])<<8
	fd.BatteryMilliVolts = int16(pl[15]) + int16(pl[16])<<8

	fd.Flying = (pl[17] & 1) == 1
	fd.OnGround = (pl[17] >> 1 & 1) == 1
	fd.EmOpen = (pl[17] >> 2 & 1) == 1
	fd.DroneHover = (pl[17] >> 3 & 1) == 1
	fd.OutageRecording = (pl[17] >> 4 & 1) == 1
	fd.BatteryLow = (pl[17] >> 5 & 1) == 1
	fd.BatteryCritical = (pl[17] >> 6 & 1) == 1
	fd.FactoryMode = (pl[17] >> 7 & 1) == 1

	fd.FlyMode = uint8(pl[18])
	fd.ThrowFlyTimer = int8(pl[19])
	fd.CameraState = uint8(pl[20])
	fd.ElectricalMachineryState = uint8(pl[21])

	fd.FrontIn = (pl[22] & 1) == 1
	fd.FrontOut = (pl[22] >> 1 & 1) == 1
	fd.FrontLSC = (pl[22] >> 2 & 1) == 1
	fd.OverTemp = (pl[23] & 1) == 1

	return fd
}

func payloadToFileInfo(pl []byte) (fType fileType, fSize uint32, fID uint16) {
	fType = fileType(pl[0])
	fSize = uint32(pl[1]) + uint32(pl[2])<<8 + uint32(pl[3])<<16 + uint32(pl[4])<<24
	fID = uint16(pl[5]) + uint16(pl[6])<<8
	return fType, fSize, fID
}

func payloadToFileChunk(pl []byte) (fc fileChunk) {
	fc.fID = uint16(pl[0]) + uint16(pl[1])<<8
	fc.pieceNum = uint32(pl[2]) + uint32(pl[3])<<8 + uint32(pl[4])<<16 + uint32(pl[5])<<24
	fc.chunkNum = uint32(pl[6]) + uint32(pl[7])<<8 + uint32(pl[8])<<16 + uint32(pl[9])<<24
	fc.chunkLen = uint16(pl[10]) + uint16(pl[11])<<8
	fc.chunkData = pl[12:]
	return fc
}

func bytesToFloat32(b []byte) (fl float32) {
	return math.Float32frombits(binary.LittleEndian.Uint32(b))
}
