package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	_ "github.com/mikeunge/audio-controller/pkg/utils"
)

// Specify the default increse/decrease volume.
const defaultVolume = 10
const appName = "audio-controller"
const appNameShort = "audic"
const version = "0.1.2"

type audio struct {
	Action  string
	Percent int
}

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

// cmdExists - make sure a command exists before executing it
func cmdExists(cmdlet string) error {
	cmd := exec.Command("which", cmdlet)	// we use `which` for the probing
	var out bytes.Buffer
	var stderr bytes.Buffer
	// redirect the stdout and stderr to the created bytebuffer.
	// this ensures the user doesn't see anything from the command probing.
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	// run the command and check for errors.
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf(fmt.Sprint(err) + ": " + stderr.String())
	}
	return nil
}

// inRange - make sure that the next element is in range before assigning it.
func inRange(min int, max int) bool {
	if min < max {
		return true
	}
	return false
}

// volume - get the current volume
func volume() (int, error) {
	err := cmdExists("pactl")
	if err != nil {
		return -1, fmt.Errorf("command 'pactl' does not exist")
	}
	// this is a weird way of executing a command but it's the only way this pipe works
	command := "pactl list sinks | grep '^[[:space:]]Volume:' | head -n $(( $SINK + 1 )) | tail -n 1 | sed -e 's,.* \\([0-9][0-9]*\\)%.*,\\1,'"
	cmd := exec.Command("bash", "-c", command)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		return -1, fmt.Errorf(fmt.Sprint(err) + ": " + stderr.String())
	}
	tmp_vol := out.String()[:len(out.String())-1]
	vol, err := strconv.Atoi(tmp_vol)
	if err != nil {
		return -1, err
	}
	return vol, nil
}

// changeVolume - execute the command so pluseaudio can increase/decrease or set the volume
func changeVolume(ad audio) error {
	percent := strconv.Itoa(ad.Percent) + "%"
	switch ad.Action {
	case "increase":
		err := cmdExists("pactl")
		if err != nil {
			return fmt.Errorf("command 'pactl' does not exist")
		}
		cmd := exec.Command("pactl", "--", "set-sink-volume", "0", "+"+percent)
		err = cmd.Run()
		if err != nil {
			return err
		}
	case "decrease":
		err := cmdExists("pactl")
		if err != nil {
			return fmt.Errorf("command 'pactl' does not exist")
		}
		cmd := exec.Command("pactl", "--", "set-sink-volume", "0", "-"+percent)
		err = cmd.Run()
		if err != nil {
			return err
		}
	case "set":
		err := cmdExists("pactl")
		if err != nil {
			return fmt.Errorf("command 'pactl' does not exist")
		}
		cmd := exec.Command("pactl", "--", "set-sink-volume", "0", percent)
		err = cmd.Run()
		if err != nil {
			return err
		}
	case "mute":
		err := cmdExists("pactl")
		if err != nil {
			return fmt.Errorf("command 'pactl' does not exist")
		}
		cmd := exec.Command("pactl", "--", "set-sink-mute", "0", "toggle")
		err = cmd.Run()
		if err != nil {
			return err
		}
	case "volume":
		vol, err := volume()
		if err != nil {
			return err
		}
		fmt.Printf("Volume: %d\n", vol)
	case "gui":
		err := cmdExists("pavucontrol")
		if err != nil {
			return fmt.Errorf("command 'pavucontrol' does not exist")
		}
		cmd := exec.Command("pavucontrol")
		fmt.Println("Press Ctrl+C to exit pavucontrol")
		err = cmd.Run()
		if err != nil {
			return err
		}
	}
	return nil
}

// argparser - parse the provided arguments
func argparser() (audio, error) {
	var ad audio
	var err error

	argv := os.Args[1:]
	if len(argv) <= 0 {
		return ad, fmt.Errorf("too few arguments provided, try audic --help for more information")
	}

	switch strings.ToLower(argv[0]) {
	case "up", "--increase", "-i":
		ad.Action = "increase"
		if inRange(1, len(argv)) {
			ad.Percent, err = strconv.Atoi(argv[1])
			if err != nil {
				return ad, fmt.Errorf("'%s' is not a valid number\n", argv[1])
			}
		} else {
			ad.Percent = defaultVolume
		}
	case "down", "--decrease", "-d":
		ad.Action = "decrease"
		if inRange(1, len(argv)) {
			ad.Percent, err = strconv.Atoi(argv[1])
			if err != nil {
				return ad, fmt.Errorf("'%s' is not a valid number\n", argv[1])
			}
		} else {
			ad.Percent = defaultVolume
		}
	case "set", "--set", "-s":
		ad.Action = "set"
		if inRange(1, len(argv)) {
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
	var ad audio

	ad, err := argparser()
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
	err = changeVolume(ad)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}
