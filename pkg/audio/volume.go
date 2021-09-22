package audio

import (
	"fmt"
	"os/exec"
	"strconv"
)

// getVolume :: returns the current volume for the given sink.
func getVolume(s *Settings) (string, error) {
	// the command gets all the available sinks, gets the volume from each one, then we get the sink WE want, after that we cleanup the string (sed) so we ONLY get the volume and not the text with the volume
	command := "pactl list sinks | grep '^[[:space:]]Volume:' | head -n $(( $SINK + " + strconv.Itoa(s.Sink) + " )) | tail -n 1 | sed -e 's,.* \\([0-9][0-9]*\\)%.*,\\1,'"
	out, err := exec.Command("bash", "-c", command).Output()
	if err != nil {
		return "", err
	}
	vol := string(out)[:len(string(out))-1]
	return vol, nil
}

// changeVolume :: as the name suggests, this function changes the volume according to the settings.Action.
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
