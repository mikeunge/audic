# audic

Small audio-controller for pulseaudio written in go.

## About

'audic' is a small cli wrapper that uses 'pactl' under the hood.
This tool makes it really easy to interact with pulseaudio via the pactl.

If you find any _bugs_ or you got a _feature_ you'd need, create a pr and I'll try to implement/fix it.

## Help

### Prerequisite

- PulseAudio + pactl
- Golang

### Build / Install

The available **make** flags are:

`make install` - will build and install the binary to your **/usr/local/bin** folder

`make run`- compiles and runs the binary with the **help** flag

`make build` - build the binary and do nothing after

### How to use **audic**

_usage_: audic \<Command> [-h|--help]

Commands:

- volume  Set/Increase/Decrease the volume
- toggle  Toggle the audio (mute/unmute)
- gui     Open a GUI (requires pavucontrol)
- about   Display information about the project

### Subcommands

_usage_: audic volume [-h|--help] [-S|--sink \<integer>] [-i|--increase \<integer>]
             [-d|--decrease \<integer>] [-s|--set \<integer>] [-m|--show]

Arguments:

- -h --help      Print help information
- -S --sink      Set the sink you want to control. Default: -1
- -i --increase  Increase the volume by N
- -d --decrease  Decrease the volume by N
- -s --set       Set the volume to N
- -m --show      Show the volume

### Features

`audic` works on _all_ Linux distributions where _pulseaudio_ is installed.

It's equiped with a "_sink_" cache, so you don't have to type it everytime.
