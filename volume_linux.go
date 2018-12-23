// +build !windows,!darwin

package volume


import (
	"errors"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

var (
useAmixer bool
outputdevice string = "Speaker"
)

const (
//	outputdevice string = "Master"
//	outputdevice string = "Speaker" // Modified from Master to Speaker for talkkonnect for raspberry pi
)

func init() {
	if _, err := exec.LookPath("pactl"); err != nil {
		useAmixer = true
	}
}

func cmdEnv() []string {
	return []string{"LANG=C", "LC_ALL=C"}
}

func getVolumeCmd(outputdevice string) []string {
	if useAmixer {
		return []string{"amixer", "get", outputdevice}
	}
	return []string{"pactl", "list", "sinks"}
}

var volumePattern = regexp.MustCompile(`\d+%`)

func parseVolume(out string) (int, error) {
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		s := strings.TrimLeft(line, " \t")
		if useAmixer && strings.Contains(s, "Playback") && strings.Contains(s, "%") ||
			!useAmixer && strings.HasPrefix(s, "Volume:") {
			volumeStr := volumePattern.FindString(s)
			return strconv.Atoi(volumeStr[:len(volumeStr)-1])
		}
	}
	return 0, errors.New("no volume found")
}

func setVolumeCmd(volume int, outputdevice string) []string {
	if useAmixer {
		return []string{"amixer", "set", outputdevice, strconv.Itoa(volume) + "%"}
	}
	return []string{"pactl", "set-sink-volume", "0", strconv.Itoa(volume) + "%"}
}

func increaseVolumeCmd(diff int, outputdevice string) []string {
	var sign string
	if diff >= 0 {
		sign = "+"
	} else if useAmixer {
		diff = -diff
		sign = "-"
	}
	if useAmixer {
		return []string{"amixer", "set", outputdevice, strconv.Itoa(diff) + "%" + sign}
	}
	return []string{"pactl", "--", "set-sink-volume", "0", sign + strconv.Itoa(diff) + "%"}
}

func getMutedCmd(outputdevice string) []string {
	if useAmixer {
		return []string{"amixer", "get", outputdevice}
	}
	return []string{"pactl", "list", "sinks"}
}

func parseMuted(out string) (bool, error) {
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		s := strings.TrimLeft(line, " \t")
		if useAmixer && strings.Contains(s, "Playback") && strings.Contains(s, "%") ||
			!useAmixer && strings.HasPrefix(s, "Mute: ") {
			if strings.Contains(s, "[off]") || strings.Contains(s, "yes") {
				return true, nil
			} else if strings.Contains(s, "[on]") || strings.Contains(s, "no") {
				return false, nil
			}
		}
	}
	return false, errors.New("no muted information found")
}

func muteCmd(outputdevice string) []string {
	if useAmixer {
		return []string{"amixer", "set", outputdevice, "mute"}
	}
	return []string{"pactl", "set-sink-mute", "0", "1"}
}

func unmuteCmd(outputdevice string) []string {
	if useAmixer {
		return []string{"amixer", "set", outputdevice, "unmute"}
	}
	return []string{"pactl", "set-sink-mute", "0", "0"}
}
