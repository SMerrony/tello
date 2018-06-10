# Tello® Package Implementation Chart
A list of the currently-known Tello functions and whether the package handles them.

| Tello Function | Dir | Package Implementation | Comments |
| -------------- | --- | ---------------------- | -------- |
| Connect  | → | ControlConnect(), ControlConnectDefault() | These funcs wait up to 3s for the Tello to respond |
| Connected | ← | ControlConnected() | (See comments for Connect) |
| Get SSID | ↔ | GetSSID() | SSID is stored in FlightData when it is received |
| Set SSID | → |  |  |
| Get SSID Password | → |  |  |
| Set SSID Password | → |  |  |
| Get Wifi Region | → |  |  |
| Set Wifi Region | → |  |  | 
| Wifi Strength | ← | Y | Handled by package - stored in FlightData |
| Set Video Bit-Rate | → | SetVideoBitrate() |  |
| Set Video Dyn. Adj. Rate | → |  |  |
| Set EIS | → |  |  |
| Request Video Start | → | StartVideo() | Use VideoConnect() first, also see VideoDisconnect() |
| Get Video Bit-Rate | → |  |  |
| Take Picture | ↔ | TakePicture() | Can also be a response, see also NumPics() and SaveAllPics() |
| Set Video Aspect | → |  |  |
| Start Recording | → |  |  |
| Exposure Values | | | |
| Light Strength | ← | Y | Handled by package - stored in FlightData |
| Get JPEG Quality | → |  |  |
| Error 1 | ← |  |  |
| Error 2 | ← |  |  |
| Get Version | → |  |  |
| Set Date & Time | ↔ | Y | Handled internally by package |
| Get Activation Time | → |  |  |
| Get Loader Version | → |  |  |
| Set Sticks | → | UpdateSticks(), StartStickListener() | also, keepAlive sends these |
| Take Off | → | TakeOff() | Ignored on receipt |
| Land | ↔ | Land(), StopLanding() | Ignored on receipt |
| Flight Status | ← | GetFlightData(), StreamFlightData() |  |
| Set Height Limit | → |  |  |
| Flip | → | Flip()  | Also see macro commands below eg. BackFlip() |
| Throw Take Off | → | ThrowTakeOff() |  |
| Palm Land | → | PalmLand() |  |
| File Size | ← | Y | Handled by package internally |
| File Data | ← | Y |  Handled by package internally |
| EOF | ← | Y | Handled by package internally |
| Start Smart Video | → |  |  |
| Get Smart Video Status | → |  |  |
| Log Header | ← |  | Currently ignored |
| Log Data | ← |  |  |
| Log Config. | ← |  |  |
| Bounce | → | Bounce() | Toggles the Bounce mode |
| Calibration | → |  |  |
| Set Low Battery Threshold | → |  |  |
| Get Height Limit | ↔ | GetMaxHeight() | MaxHeight stored in FlightData when it is received |
| Get Low Battery Threshold | → |  |  |
| Get Attitude (Limit?) | → |  |  |
| Set Attitude (Limit?) | → |  |  |

## Macro Commands

| Command | Comments |
| ------- | -------- |
| Hover() | Stop motion |
| Forward(), Backward(), Left(), Right() | Start moving at given percentage of max speed |
| Up(), Down() | Start moving at given percentage of max speed |
| Clockwise(), Anticlockwise() | aliases: TurnLeft(), TurnRight(), CounterClockwise() - Start turning at given percentage of max rate |
| SetSportsMode(), SetFastMode, SetSlowMode | |
| BackFlip(), BackLeftFlip(), BackRightFlip(), ForwardFlip(), etc. |  |
