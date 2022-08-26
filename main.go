package main

import (
	"fmt"
	"os"

	"github.com/akamensky/argparse"

	"github.com/mikeunge/audic/pkg/audio"
)

const (
	appName    = "audic"
	appVersion = "0.1.6"
	appRelease = "@release"
	appRepo    = "github.com/mikeunge/audic"
	appDev     = "mikeunge"
)

func parser(settings *audio.Settings) error {
	parser := argparse.NewParser(appName, "The easiest way to control your audio.")

	// Parser options
	volume := parser.NewCommand("volume", "Set/Increase/Decrease the volume")
	sinkV := volume.Int("S", "sink", &argparse.Options{Required: false, Help: "Set the sink you want to control", Default: -1})
	increase := volume.Int("i", "increase", &argparse.Options{Required: false, Help: "Increase the volume by N"})
	decrease := volume.Int("d", "decrease", &argparse.Options{Required: false, Help: "Decrease the volume by N"})
	set := volume.Int("s", "set", &argparse.Options{Required: false, Help: "Set the volume to N"})
	show := volume.Flag("m", "show", &argparse.Options{Required: false, Help: "Show the volume"})
	toggle := parser.NewCommand("toggle", "Toggle the audio (mute/unmute)")
	sinkT := toggle.Int("S", "sink", &argparse.Options{Required: false, Help: "Set the sink you want to control", Default: -1})
	gui := parser.NewCommand("gui", "Open a GUI (requires pavucontrol)")
	about := parser.NewCommand("about", "Display information about the project")
	listSinks := parser.NewCommand("list", "List available sinks")

	// Parse the input
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage("")) // print help
		return fmt.Errorf("parser returned with an error, no arguments provided")
	}

	if toggle.Happened() {
		settings.Action = "toggle"
		audio.SinkVal(sinkT, settings)
		return nil
	}

	if gui.Happened() {
		settings.Action = "gui"
		return nil
	}

	if listSinks.Happened() {
		settings.Action = "listSink"
		return nil
	}

	if about.Happened() {
		fmt.Printf("%s - v%s %s\n\n", appName, appVersion, appRelease)
		fmt.Printf("Repository: %s\n", appRepo)
		fmt.Printf("Developed:  %s\n\n", appDev)
		fmt.Print(parser.Usage(""))
		os.Exit(0)
	}

	// parse if volume command was set
	// check for all the subcommands and their values
	if volume.Happened() {
		audio.SinkVal(sinkV, settings)
		if *increase != 0 {
			settings.Action = "increase"
			settings.Percent = *increase
		} else if *decrease != 0 {
			settings.Action = "decrease"
			settings.Percent = *decrease
		} else if *set != 0 {
			settings.Action = "set"
			settings.Percent = *set
		} else if *show {
			settings.Action = "show"
		} else {
			return fmt.Errorf("no subcommand provided, use -h/--help for more information")
		}
	}
	return nil
}

func main() {
	var settings audio.Settings
	settings.SinkPath = "/tmp/audic.sink" // where the sink cache is located
	settings.SetSink = false              // trigger, if we should write to cache (will get changed later)

	// parse the provided arguments
	err := parser(&settings)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	// pass the instructions to the controller
	err = audio.Controller(&settings)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	// write the sink to cache
	if settings.SetSink {
		// we really don't care if it failed or not
		// this is only interesting when we debug the program
		_ = audio.WriteSink(&settings)
	}
	os.Exit(0)
}
