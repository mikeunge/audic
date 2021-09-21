package audio

import (
	"fmt"
	"bytes"
	"os/exec"
	"strconv"

	utils "github.com/mikeunge/audic/pkg/utils"
)

type Audio struct {
	Action  string
	Percent int
	Sink 	int
}

// GetVolume - get the current volume
func GetVolume() (int, error) {
	err := utils.CmdExists("pactl")
	if err != nil {
		return -1, fmt.Errorf("command 'pactl' does not exist")
	}
	sink := "1"	// todo: pass the sink to use
	// the command gets all the available sinks, gets the volume from each one, then we get the sink WE want, after that we cleanup the string (sed) so we ONLY get the volume and not the text with the volume
	command := "pactl list sinks | grep '^[[:space:]]Volume:' | head -n $(( $SINK + "+sink+" )) | tail -n 1 | sed -e 's,.* \\([0-9][0-9]*\\)%.*,\\1,'"
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

// ChangeVolume - execute the command so pluseaudio can increase/decrease or set the volume
func ChangeVolume(ad Audio) error {
	percent := strconv.Itoa(ad.Percent) + "%"
	switch ad.Action {
	case "increase":
		err := utils.CmdExists("pactl")
		if err != nil {
			return fmt.Errorf("command 'pactl' does not exist")
		}
		cmd := exec.Command("pactl", "--", "set-sink-volume", "0", "+"+percent)
		err = cmd.Run()
		if err != nil {
			return err
		}
	case "decrease":
		err := utils.CmdExists("pactl")
		if err != nil {
			return fmt.Errorf("command 'pactl' does not exist")
		}
		cmd := exec.Command("pactl", "--", "set-sink-volume", "0", "-"+percent)
		err = cmd.Run()
		if err != nil {
			return err
		}
	case "set":
		err := utils.CmdExists("pactl")
		if err != nil {
			return fmt.Errorf("command 'pactl' does not exist")
		}
		cmd := exec.Command("pactl", "--", "set-sink-volume", "0", percent)
		err = cmd.Run()
		if err != nil {
			return err
		}
	case "mute":
		err := utils.CmdExists("pactl")
		if err != nil {
			return fmt.Errorf("command 'pactl' does not exist")
		}
		cmd := exec.Command("pactl", "--", "set-sink-mute", "0", "toggle")
		err = cmd.Run()
		if err != nil {
			return err
		}
	case "volume":
		vol, err := GetVolume()
		if err != nil {
			return err
		}
		fmt.Printf("Volume: %d\n", vol)
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

