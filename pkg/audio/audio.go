package audio

import (
	"fmt"
	"os/exec"
	"strconv"

	utils "github.com/mikeunge/audic/pkg/utils"
)

type Settings struct {
	Action  string
	Percent int
	Sink 	int
}


// sinkExists :: check if the provided sink exists or not
func sinkExists(sink int) bool {
	cmd := "pactl list sinks | grep 'Sink #"+ strconv.Itoa(sink) +"'"
	out, err := exec.Command("bash", "-c", cmd).Output()
	// Check if we didn't get any errors and if the grep output isn't empty.
	if err != nil || string(out) == "" {
		return false
	}
	return true
}


// getVolume :: get the current volume
func getVolume(s *Settings) (string, error) {
	// the command gets all the available sinks, gets the volume from each one, then we get the sink WE want, after that we cleanup the string (sed) so we ONLY get the volume and not the text with the volume
	command := "pactl list sinks | grep '^[[:space:]]Volume:' | head -n $(( $SINK + "+ strconv.Itoa(s.Sink) +" )) | tail -n 1 | sed -e 's,.* \\([0-9][0-9]*\\)%.*,\\1,'"
	out, err := exec.Command("bash", "-c", command).Output()
	if err != nil {
		return "", err
	}
	vol := string(out)[:len(string(out))-1]
	return vol, nil
}


// changeVolume :: as the name suggests, this function changes the volume according to the settings.Action
func changeVolume(s *Settings) error {
	var cmd exec.Cmd
	percent := strconv.Itoa(s.Percent) + "%"

	switch s.Action {
	case "increase":
		cmd = *exec.Command("pactl", "set-sink-volume", strconv.Itoa(s.Sink), "+"+percent)
	case "decrease":
		cmd = *exec.Command("pactl", "set-sink-volume", strconv.Itoa(s.Sink), "-"+percent)
	case "set":
		cmd = *exec.Command("pactl", "set-sink-volume", strconv.Itoa(s.Sink), percent)
	default:
		return fmt.Errorf("action was not found, aborting")
	}

	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
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
		cmd := exec.Command("pavucontrol")
		fmt.Println("Press Ctrl+C to exit pavucontrol")
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

