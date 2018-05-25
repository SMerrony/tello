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

const msgHdr = 0xcc // 204

type packet struct {
	header        byte
	size13        uint16
	crc8          byte
	toDrone       bool // the following 4 fields are encoded in a single byte in the raw packet
	fromDrone     bool
	packetType    byte // 3-bit
	packetSubtype byte // 3-bit
	message       uint16
	sequence      uint16
	payload       []byte
	crc16         uint16
}

// tello packet types, 3 and 7 currently unknown
const (
	ptExtended = 0
	ptGet      = 1
	ptData1    = 2
	ptData2    = 4
	ptSet      = 5
	ptFlip     = 6
)

// Tello messages
const (
	msgDoConnect         = 0x0001 // 1
	msgConnected         = 0x0002 // 2
	msgGetSSID           = 0x0011 // 17
	msgSetSSID           = 0x0012 // 18
	msgGetSSIDPass       = 0x0013 // 19
	msgSetSSIDPass       = 0x0014 // 20
	msgGetWifiRegion     = 0x0015 // 21
	msgSetWifiRegion     = 0x0016 // 22
	msgWifiStrength      = 0x001a // 26
	msgSetVideoBitrate   = 0x0020 // 32
	msgSetDynAdjRate     = 0x0021 // 33
	msgEisSetting        = 0x0024 // 36
	msgGetVideoSPSPPS    = 0x0025 // 37
	msgGetVideoBitrate   = 0x0028 // 40
	msgDoTakePic         = 0x0030 // 48
	msgSwitchPicVideo    = 0x0031 // 49
	msgDoStartRec        = 0x0032 // 50
	msgExposureVals      = 0x0034 // 52 (Get or set?)
	msgLightStrength     = 0x0035 // 53
	msgGetJPEGQuality    = 0x0037 // 55
	msgError1            = 0x0043 // 67
	msgError2            = 0x0044 // 68
	msgGetVersion        = 0x0045 // 69
	msgSetDateTime       = 0x0046 // 70
	msgGetActivationTime = 0x0047 // 71
	msgGetLoaderVersion  = 0x0049 // 73
	msgSetStick          = 0x0050 // 80
	msgDoTakeoff         = 0x0054 // 84
	msgDoLand            = 0x0055 // 85
	msgFlightStatus      = 0x0056 // 86
	msgSetHeightLimit    = 0x0058 // 88
	msgDoFlip            = 0x005c // 92
	msgDoThrowTakeoff    = 0x005d // 93
	msgDoPalmLand        = 0x005e // 94
	msgFileSize          = 0x0062 // 98
	msgFileData          = 0x0063 // 99
	msgFileDone          = 0x0064 // 100
	msgDoSmartVideo      = 0x0080 // 128
	msgGetSmartVideo     = 0x0081 // 129
	msgLogHeader         = 0x1050 // 4176
	msgLogData           = 0x1051 // 4177
	msgLogConfig         = 0x1052 // 4178
	msgDoBounce          = 0x1053 // 4179
	msgDoCalibration     = 0x1054 // 4180
	msgSetLowBattThresh  = 0x1055 // 4181
	msgGetHeightLimit    = 0x1056 // 4182
	msgGetLowBattThresh  = 0x1057 // 4183
	msgSetAttitude       = 0x1058 // 4184
	msgGetAttitude       = 0x1059 // 4185
)

// Flip types
const (
	flipForward = iota
	flipLeft
	flipBackward
	flipRight
	flipForwardLeft
	flipBackwardLeft
	flipBackwardRight
	flipForwardRight
)

// Smart Video messages
const (
	svStop   = 0
	sv360    = 1
	svCircle = 2
	svUpOut  = 3
)

// video bit rate (mbps)
const (
	vbrAuto = iota
	vbr1M
	vbr1M5
	vbr2M
	vbr3M
	vbr4M
)

// FlightData payload from the Tello
type FlightData struct {
	BatteryLow               bool
	BatteryLower             bool
	BatteryPercentage        int8
	BatteryState             bool
	CameraState              int8
	DownVisualState          bool
	DroneBatteryLeft         int16
	DroneFlyTimeLeft         int16
	DroneHover               bool
	EmOpen                   bool
	EmSky                    bool
	EmGround                 bool
	EastSpeed                int16
	ElectricalMachineryState int16
	FactoryMode              bool
	FlyMode                  int8
	FlySpeed                 int16
	FlyTime                  int16
	FrontIn                  bool
	FrontLSC                 bool
	FrontOut                 bool
	GravityState             bool
	GroundSpeed              int16
	Height                   int16
	ImuCalibrationState      int8
	ImuState                 bool
	LightStrength            int8
	NorthSpeed               int16
	OutageRecording          bool
	PowerState               bool
	PressureState            bool
	SmartVideoExitMode       int16
	TemperatureHeight        bool
	ThrowFlyTimer            int8
	WifiDisturb              int8
	WifiStrength             int8
	WindState                bool
}
