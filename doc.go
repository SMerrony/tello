/*Package tello provides an unofficial, easy-to-use, standalone API for the Ryze Tello® drone.

Disclaimer

Tello is a registered trademark of Ryze Tech.  The author(s) of this package is/are in no way affiliated with Ryze, DJI, or Intel.
The package has been developed by gathering together information from a variety of sources on the Internet
(especially the generous contributors at  https://tellopilots.com), and by examining data packets sent to/from the Tello.
The package will probably be extended as more knowledge of the drone's protocol is obtained.

Use this package at your own risk.  The author(s) is/are in no way responsible for any damage caused either to or by the
drone when using this software.

Features

The following features have been implemented...
  * Stick-based flight control, ie. for joystick, game-, or flight-controller
  * Drone built-in flight commands, eg. Takeoff(), PalmLand()
  * Macro-level flight control, eg. Forward(), Up()
  * Autopilot commands, eg. AutoFlyToHeight(), AutoTurnToYaw()
  * Enriched flight data (some log data is added) for real-time telemetry
  * Video stream support
  * Picture taking/saving
  * Multiple drone support - Untested
An example application using this package is available at http://github.com/SMerrony/telloterm

This documentation should be consulted alongside https://github.com/SMerrony/tello/blob/master/ImplementationChart.md

Concepts

Connection Types

The drone provides two types of connection: a 'control' connection which handles all commands
to and from the drone including flight, status and (still) pictures, and a 'video' connection which
provides an H.264 video stream from the forward-facing camera.  You must establish a control connection to use the drone,
but the video connection is optional and cannot be started unless a control connection is running.

Funcs vs. Channels

Certain functionality is made available in two forms: single-shot function calls and streaming (channel) data flows.
Eg. GetFlightData() vs. StreamFlightData(), and UpdateSticks() vs. StartStickListener().

Use whichever paradigm you prefer, but be aware that the channel-based calls should return immediately (the channels are buffered)
whereas the function-based options could conceivably cause your application to pause very briefly if the Tello is very busy.
(In practice, the author has not found this to be an issue.)

*/
package tello
