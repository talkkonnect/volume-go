package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

var volumePattern = regexp.MustCompile(`\d+%`)
var useAmixer bool = true

//Example For CM108 (USB Sound Card)
var OutputVolDevice string ="Speaker"
var OutputMuteDevice string = "Speaker"
var OutputVolRegex string = "Playback"

//Example For WM8960 (RESPEAKER)
//var OutputVolDevice string ="Speaker"
//var OutputMuteDevice string = "Right Output Mixer PCM"
//var OutputVolRegex string = "Front Left:"


func main() {

	vol, err := GetVolume(OutputVolDevice)
	if err != nil {
		log.Fatalf("get volume failed: %+v", err)
	}
	fmt.Printf("current volume: %d\n", vol)

	var NewVol int = 100
	if vol < 100 {
		NewVol = vol + 1
	}

	err = SetVolume(NewVol, OutputVolDevice)
	if err != nil {
		log.Fatalf("set volume failed: %+v", err)
	}
	fmt.Printf("set volume success\n")

	err = Mute(OutputVolDevice)
	if err != nil {
		log.Fatalf("mute failed: %+v", err)
	} else {
		log.Println("mute success")
	}

	err = Unmute(OutputVolDevice)
	if err != nil {
		log.Fatalf("unmute failed: %+v", err)
	} else {
		log.Println("unmute success")

	}
}

func GetVolume(outputdevice string) (int, error) {
	out, err := execCmd(getVolumeCmd(outputdevice))
	if err != nil {
		return 0, err
	}
	return parseVolume(string(out))
}

func getVolumeCmd(outputdevice string) []string {
	return []string{"amixer", "get", outputdevice}
}

func execCmd(cmdArgs []string) ([]byte, error) {
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	cmd.Env = append(os.Environ(), cmdEnv()...)
	out, err := cmd.Output()
	if err != nil {
		err = fmt.Errorf(`failed to execute "%v" (%+v)`, strings.Join(cmdArgs, " "), err)
	}
	return out, err
}

func parseVolume(out string) (int, error) {
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		s := strings.TrimLeft(line, " \t")
		if useAmixer && strings.Contains(s, OutputVolRegex) && strings.Contains(s, "%") ||
			!useAmixer && strings.HasPrefix(s, "Volume:") {
			volumeStr := volumePattern.FindString(s)
			return strconv.Atoi(volumeStr[:len(volumeStr)-1])
		}
	}
	return 0, errors.New("no volume found")
}

func cmdEnv() []string {
	return []string{"LANG=C", "LC_ALL=C"}
}

func SetVolume(volume int, outputdevice string) error {
	if volume < 0 || 100 < volume {
		return errors.New("out of valid volume range")
	}
	_, err := execCmd(setVolumeCmd(volume, outputdevice))
	return err
}

func setVolumeCmd(volume int, outputdevice string) []string {
	if useAmixer {
		return []string{"amixer", "set", outputdevice, strconv.Itoa(volume) + "%"}
	} else if _, err := strconv.Atoi(outputdevice); err == nil {
		return []string{"pactl", "set-sink-volume", outputdevice, strconv.Itoa(volume) + "%"}
	}
	return []string{"pactl", "set-sink-volume", "0", strconv.Itoa(volume) + "%"}
}

// Mute mutes the audio.
func Mute(outputdevice string) error {
	_, err := execCmd(muteCmd(outputdevice))
	return err
}

func muteCmd(outputdevice string) []string {
	if useAmixer {
		return []string{"amixer", "set", outputdevice, "mute"}
	} else if _, err := strconv.Atoi(outputdevice); err == nil {
		return []string{"pactl", "set-sink-mute", outputdevice, "1"}
	}
	return []string{"pactl", "set-sink-mute", "0", "1"}
}

func Unmute(outputdevice string) error {
	_, err := execCmd(unmuteCmd(outputdevice))
	return err
}

func unmuteCmd(outputdevice string) []string {
	if useAmixer {
		return []string{"amixer", "set", outputdevice, "unmute"}
	} else if _, err := strconv.Atoi(outputdevice); err == nil {
		return []string{"pactl", "set-sink-mute", outputdevice, "0"}
	}
	return []string{"pactl", "set-sink-mute", "0", "0"}
}

