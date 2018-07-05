# Tello® Package Implementation Chart
A list of the currently-known Tello functions and whether or how this package handles them.

This should be consulted alongside the GoDoc package documentation.  See https://godoc.org/github.com/SMerrony/tello or your local equivalent.

| ID (Hex) | Tello Function | Dir | Package Implementation | Comments |
| -------- | -------------- | --- | ---------------------- | -------- |
| 0x0001 | Connect  | → | ControlConnect(), ControlConnectDefault() | These funcs wait up to 3s for the Tello to respond |
| 0x0002 | Connected | ← | ControlConnected() | (See comments for Connect) |
| 0x0011 | Query SSID | ↔ | GetSSID() | SSID is stored in FlightData when it is received |
| 0x0012 | Set SSID | → |  |  |
| 0x0013 | Query SSID Password | → |  |  |
| 0x0014 | Set SSID Password | → |  |  |
| 0x0015 | Query Wifi Region | → |  |  |
| 0x0016 | Set Wifi Region | → |  |  | 
| 0x001a | Wifi Strength | ← | Y | Handled internally by package - stored in FlightData |
| 0x0020 | Set Video Bit-Rate | → | SetVideoBitrate() |  |
| 0x0021 | Set Video Dyn. Adj. Rate | → |  |  |
| 0x0024 | Set EIS | → |  |  |
| 0x0025 | Request Video Start | → | StartVideo() | Use VideoConnect() first, also see VideoDisconnect() |
| 0x0028 | Query Video Bit-Rate | ↔ | GetVideoBitrate() |  |
| 0x0030 | Take Picture | ↔ | TakePicture() | Can also be a response, see also NumPics() and SaveAllPics() |
| 0x0031 | Set Video Aspect | ↔ | SetVideoNormal() & SetVideoWide() |  |
| 0x0032 | Start Recording | → |  |  |
| 0x0034 | Exposure Values | | | |
| 0x0035 | Light Strength | ← | Y | Handled internally by package - stored in FlightData |
| 0x0037 | Query JPEG Quality | → |  |  |
| 0x0043 | Error 1 | ← |  |  |
| 0x0044 | Error 2 | ← |  |  |
| 0x0045 | Query Version | ↔ | GetVersion() |  |
| 0x0046 | Set Date & Time | ↔ | Y | Handled internally by package |
| 0x0047 | Query Activation Time | → |  |  |
| 0x0049 | Query Loader Version | → |  |  |
| 0x0050 | Set Sticks | → | UpdateSticks(), StartStickListener() | also, keepAlive sends these |
| 0x0054 | Take Off | → | TakeOff() | Ignored on receipt |
| 0x0055 | Land | ↔ | Land(), StopLanding() | Ignored on receipt |
| 0x0056 | Flight Status | ← | GetFlightData(), StreamFlightData() |  |
| 0x0058 | Set Height Limit | → |  |  |
| 0x005c | Flip | → | Flip()  | Also see macro commands below eg. BackFlip() |
| 0x005d | Throw Take Off | → | ThrowTakeOff() |  |
| 0x005e | Palm Land | → | PalmLand() |  |
| 0x0062 | File Size | ← | Y | Handled internally by package |
| 0x0063 | File Data | ← | Y |  Handled internally by package |
| 0x0064 | EOF | ← | Y | Handled internally by package |
| 0x0080 | Start Smart Video | → | StartSmartVideo(), StopSmartVideo() |  |
| 0x0081 | Smart Video Status | ← |  |  |
| 0x1050 | Log Header | ↔ |  | Handled internally by package |
| 0x1051 | Log Data | ← |  | Some MOV and IMU data are captured and added to FlightData |
| 0x1052 | Log Config. | ← |  |  |
| 0x1053 | Bounce | → | Bounce() | Toggles the Bounce mode |
| 0x1054 | Calibration | → |  |  |
| 0x1055 | Set Low Battery Threshold | ↔ | SetLowBatteryThreshold() | (See godoc) |
| 0x1056 | Query Height Limit | ↔ | GetMaxHeight() | MaxHeight stored in FlightData when it is received |
| 0x1057 | Query Low Battery Threshold | ↔ | GetLowBatteryThreshold() |  |
| 0x1058 | Query Attitude (Limit?) | → |  |  |
| 0x1059 | Set Attitude (Limit?) | → |  |  |

## Macro and Flight Commands

| Low-Level Command | Macro Command | Comments |
| ----------------- | ------------- | -------- |
| UpdateSticks() | Set joystick position (macro commands below) |
| | Hover() | Stop motion |
| | Forward(), Backward(), Left(), Right(), Up(), Down()| Start moving at given percentage of max speed |
| |Clockwise(), Anticlockwise() | aliases: TurnLeft(), TurnRight(), CounterClockwise() - Start turning at given percentage of max rate |
| | AutoFlyToHeight(), AutoTurnToYaw(), AutoTurnByDeg(), AutoFlyToXY() | Fly automatically to specified height/yaw/pos (can use concurrently) |
| SetSportsMode() | Also SetFastMode(), SetSlowMode() |
| Flip() | Also BackFlip(), BackLeftFlip(), BackRightFlip(), ForwardFlip(), etc. |
| StartSmartVideo(), StopSmartVideo() | eg. 360 rotation, circle, up-and-out |
