package main

import (
	"fmt"
	"os"

	"github.com/akamensky/argparse"

	audio "github.com/mikeunge/audic/pkg/audio"
)

// Specify the default increse/decrease volume.
const defaultVolume = 10
const appName = "audic"
const version = "0.1.3.1"

func parser(s *audio.Settings) error {
	parser := argparse.NewParser(appName, "The easiest way to control your audio.")

	// Parser options
	volume := parser.NewCommand("volume", "Set/Increase/Decrease the volume")
	sinkV := volume.Int("S", "sink", &argparse.Options{Required: false, Help: "Set the sink you want to control", Default: 0})
	increase := volume.Int("i", "increase", &argparse.Options{Required: false, Help: "Increase the volume by N"})
	decrease := volume.Int("d", "decrease", &argparse.Options{Required: false, Help: "Decrease the volume by N"})
	set := volume.Int("s", "set", &argparse.Options{Required: false, Help: "Set the volume to N"})
	show := volume.Flag("m", "show", &argparse.Options{Required: false, Help: "Show the volume eg. Volume: 80%"})
	toggle := parser.NewCommand("toggle", "Toggle the audio (mute/unmute)")
	sinkT := toggle.Int("S", "sink", &argparse.Options{Required: false, Help: "Set the sink you want to control", Default: 0})
	gui := parser.NewCommand("gui", "Open a GUI (requires pavucontrol)")

	// Parse the options
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(""))
		return fmt.Errorf("parser returned with an error, no arguments provided")
	}

	if toggle.Happened() {
		s.Action = "toggle"
		s.Sink = *sinkT
		return nil
	}

	if gui.Happened() {
		s.Action = "gui"
		return nil
	}

	if volume.Happened() {
		s.Sink = *sinkV
		if (*increase != 0) {
			s.Action = "increase"
			s.Percent = *increase
		} else if (*decrease != 0) {
			s.Action = "decrease"
			s.Percent = *decrease
		} else if (*set != 0) {
			s.Action = "set"
			s.Percent = *set
		} else if *show {
			s.Action = "show"
		} else {
			return fmt.Errorf("no subcommand provided, use -h/--help for more information")
		}
	}
	return nil
}


func main() {
	var s audio.Settings

	err := parser(&s)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	err = audio.ChangeVolume(&s)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}
