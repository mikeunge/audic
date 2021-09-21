package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	utils "github.com/mikeunge/audic/pkg/utils"
	audio "github.com/mikeunge/audic/pkg/audio"
)

// Specify the default increse/decrease volume.
const defaultVolume = 10
const appName = "audio-controller"
const appNameShort = "audic"
const version = "0.1.2"


// help - print useful information
func help() {
	fmt.Printf("%s (%s) -- v%s\n\n", appName, appNameShort, version)
	fmt.Printf("Usage: %s set N   --   set the volume to N%%\n\n", appNameShort)
	fmt.Println("Available commands/arguments:")
	fmt.Println("  -     up N\t Increase the volume by N percent")
	fmt.Println("  -   down N\t Decrease the volume by N percent")
	fmt.Println("  -    set N\t Set the volume to N percent")
	fmt.Println("  - volume\t Show the current volume")
	fmt.Println("  -   mute\t Mute/Unmute the audio")
	fmt.Println("  -   gui\t Spawn a GUI via pavucontrol")
	fmt.Println("  -   help\t Show this message and exit")
}


// argparser - parse the provided arguments
func argparser() (audio.Audio, error) {
	var ad audio.Audio
	var err error

	argv := os.Args[1:]
	if len(argv) <= 0 {
		return ad, fmt.Errorf("too few arguments provided, try audic --help for more information")
	}

	switch strings.ToLower(argv[0]) {
	case "up", "--increase", "-i":
		ad.Action = "increase"
		if utils.InRange(1, len(argv)) {
			ad.Percent, err = strconv.Atoi(argv[1])
			if err != nil {
				return ad, fmt.Errorf("'%s' is not a valid number\n", argv[1])
			}
		} else {
			ad.Percent = defaultVolume
		}
	case "down", "--decrease", "-d":
		ad.Action = "decrease"
		if utils.InRange(1, len(argv)) {
			ad.Percent, err = strconv.Atoi(argv[1])
			if err != nil {
				return ad, fmt.Errorf("'%s' is not a valid number\n", argv[1])
			}
		} else {
			ad.Percent = defaultVolume
		}
	case "set", "--set", "-s":
		ad.Action = "set"
		if utils.InRange(1, len(argv)) {
			ad.Percent, err = strconv.Atoi(argv[1])
			if err != nil {
				return ad, fmt.Errorf("'%s' is not a valid number\n", argv[1])
			}
		} else {
			return ad, fmt.Errorf("you need to specify a volume to set")
		}
	case "mute", "--mute", "-m":
		ad.Action = "mute"
	case "volume", "--volume", "-v":
		ad.Action = "volume"
	case "gui", "--gui", "-g":
		ad.Action = "gui"
	case "help", "--help", "-h":
		help()
		os.Exit(0)
	default:
		help()
		return ad, fmt.Errorf("'%s' is not a valid argument\n", argv[0])
	}
	return ad, nil
}

func main() {
	var ad audio.Audio

	ad, err := argparser()
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
	err = audio.ChangeVolume(ad)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}
