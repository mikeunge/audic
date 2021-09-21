package audio

import (
	"fmt"
	"bytes"
	"os/exec"
	"strconv"

	utils "github.com/mikeunge/audic/pkg/utils"
)

type Settings struct {
	Action  string
	Percent int
	Sink 	int
}

// getVolume - get the current volume
func getVolume(s *Settings) (string, error) {
	err := utils.CmdExists("pactl")
	if err != nil {
		return "", fmt.Errorf("command 'pactl' does not exist")
	}

	// the command gets all the available sinks, gets the volume from each one, then we get the sink WE want, after that we cleanup the string (sed) so we ONLY get the volume and not the text with the volume
	command := "pactl list sinks | grep '^[[:space:]]Volume:' | head -n $(( $SINK + "+ strconv.Itoa(s.Sink) +" )) | tail -n 1 | sed -e 's,.* \\([0-9][0-9]*\\)%.*,\\1,'"
	cmd := exec.Command("bash", "-c", command)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf(fmt.Sprint(err) + ": " + stderr.String())
	}

	vol := out.String()[:len(out.String())-1]
	return vol, nil
}

// ChangeVolume - execute the command so pluseaudio can increase/decrease or set the volume
func ChangeVolume(s *Settings) error {
	percent := strconv.Itoa(s.Percent) + "%"
	switch s.Action {
	case "increase":
		err := utils.CmdExists("pactl")
		if err != nil {
			return fmt.Errorf("command 'pactl' does not exist")
		}
		cmd := exec.Command("pactl", "set-sink-volume", strconv.Itoa(s.Sink), "+"+percent)
		err = cmd.Run()
		if err != nil {
			return err
		}
	case "decrease":
		err := utils.CmdExists("pactl")
		if err != nil {
			return fmt.Errorf("command 'pactl' does not exist")
		}
		cmd := exec.Command("pactl", "set-sink-volume", strconv.Itoa(s.Sink), "-"+percent)
		err = cmd.Run()
		if err != nil {
			return err
		}
	case "set":
		err := utils.CmdExists("pactl")
		if err != nil {
			return fmt.Errorf("command 'pactl' does not exist")
		}
		cmd := exec.Command("pactl", "set-sink-volume", strconv.Itoa(s.Sink), percent)
		err = cmd.Run()
		if err != nil {
			return err
		}
	case "toggle":
		err := utils.CmdExists("pactl")
		if err != nil {
			return fmt.Errorf("command 'pactl' does not exist")
		}
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
	}
	return nil
}

