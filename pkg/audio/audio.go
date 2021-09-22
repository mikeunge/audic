package audio

import (
	"fmt"
	"os/exec"
	"strconv"

	"github.com/mikeunge/audic/pkg/utils"
)

type Settings struct {
	Action   string
	Percent  int
	Sink     int
	SetSink  bool
	SinkPath string
}

// Controller :: decide what to do and take action accordingly
func Controller(s *Settings) error {
	// check if needed (pactl) command exists
	err := utils.CmdExists("pactl")
	if err != nil {
		return fmt.Errorf("command 'pactl' does not exist")
	}

	// make sure the requested sink exists
	if !sinkExists(s.Sink) {
		return fmt.Errorf("provided sink does not exist")
	}

	switch s.Action {
	case "toggle":
		cmd := exec.Command("pactl", "set-sink-mute", strconv.Itoa(s.Sink), "toggle")
		err = cmd.Run()
		if err != nil {
			return err
		}
	case "show":
		vol, err := getVolume(s)
		if err != nil {
			return err
		}
		fmt.Printf("Volume: %s%%\n", vol)
	case "gui":
		err := utils.CmdExists("pavucontrol")
		if err != nil {
			return fmt.Errorf("command 'pavucontrol' does not exist")
		}
		fmt.Println("Press Ctrl+C to exit pavucontrol")
		cmd := exec.Command("pavucontrol")
		err = cmd.Run()
		if err != nil {
			return err
		}
	default:
		err := changeVolume(s)
		if err != nil {
			return err
		}
	}
	return nil
}
