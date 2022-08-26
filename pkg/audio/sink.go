package audio

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
)

// SinkVal :: check the value of provided sink.
//
// Check if we deal with the default value of sink (-1) or if the user assigned one.
// If we deal with the default, we check if we can read the value from "cache", else, we assign 0
func SinkVal(val *int, settings *Settings) {
	var sink int
	if *val == -1 {
		tmp, err := GetSink(&settings.SinkPath)
		if err == nil {
			sink = tmp
		} else {
			settings.SetSink = true
			sink = 0
		}
	} else {
		settings.SetSink = true
		sink = *val
	}
	settings.Sink = sink
}

// listSink :: list all available sinks 
func listSink() (string, error) {
	cmd := "pactl list sinks"
	out, err := exec.Command("bash", "-c", cmd).Output()
	// Check if we didn't get any errors.
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// sinkExists :: check if the provided sink exists or not
func sinkExists(sink int) bool {
	cmd := "pactl list sinks | grep 'Sink #" + strconv.Itoa(sink) + "'"
	out, err := exec.Command("bash", "-c", cmd).Output()
	// Check if we didn't get any errors and if the grep output isn't empty.
	if err != nil || string(out) == "" {
		return false
	}
	return true
}

// WriteSink :: write the sink to "cache".
func WriteSink(settings *Settings) error {
	// create a new byte array that can hold up to two byte (65535)
	if settings.Sink >= 65536 {
		return fmt.Errorf("cannot write sink to cache, sink size is too large")
	}
	data := make([]byte, 2)
	binary.LittleEndian.PutUint16(data, uint16(settings.Sink))

	err := ioutil.WriteFile(settings.SinkPath, data, 0664)
	if err != nil {
		return err
	}

	return nil
}

// GetSink :: check if a sink is "cached", if so, laod and return it.
func GetSink(path *string) (int, error) {
	if _, err := os.Stat(*path); os.IsNotExist(err) {
		return 0, fmt.Errorf(fmt.Sprintf("could not find file %s\n", *path))
	}

	// read the file ([]bytes) and convert it back into a uint16, then return as an int.
	slice, err := ioutil.ReadFile(*path)
	if err != nil {
		return 0, err
	}
	data := binary.LittleEndian.Uint16(slice)
	return int(data), nil
}
