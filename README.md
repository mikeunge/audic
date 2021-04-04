# audic

Small audio-controller for pulseaudio written in go.

### About

'audic' is a small cli wrapper that uses 'pactl' under the hood.
This tool makes it really easy to interact with pulseaudio via the pactl.

If you find any _bugs_ or you got a _feature_ you'd need, create a pr and I'll try to implement/fix it.

### Help

#### Prerequisite

-   PulseAudio + pactl
-   Golang

#### Usage

The available **make** flags are:

´make install´ - will build and install the binary to your **/usr/local/bin** folder
´make run´ - compiles and runs the binary with the **help** flag
´make build´ - build the binary and do nothing after

How to use **audic**:

-   audic set 100 - sets the audio volume to 100%
-   audic up - increases the audio volume by 10% (_this is the fixed default increase_)
-   audic up 25 - increases the audio volume by 25% (_this overrides the default increase_)

#### Available arguments

-   up (--increase / -i)
    -   Increase the volume by N %
-   down (--decrease / -d)
    -   Decrease the volume by N %
-   set (--set / -s)
    -   Set the volume N %
-   mute (--mute / -m)
    -   Mute the audio
-   unmute (--unmute / -u)
    -   Unmute the audio
    -   Unmute the audio
