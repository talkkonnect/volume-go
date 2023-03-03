package volume

import (
	"fmt"
	"strconv"
	"strings"
)

func cmdEnv() []string {
	return nil
}

func getVolumeCmd(outputdevice string) []string {
	return []string{"osascript", "-e", "output volume of (get volume settings)"}
}

func parseVolume(out string) (int, error) {
	out = strings.TrimSuffix(out, "\n")
	if out == "missing value" {
		return 0, fmt.Errorf("failed to get volume settings: %s", out)
	}
	return strconv.Atoi(out)
}

func setVolumeCmd(volume int, outputdevice string) []string {
	return []string{"osascript", "-e", "set volume output volume " + strconv.Itoa(volume)}
}

func increaseVolumeCmd(diff int, outputdevice string) []string {
	return []string{"osascript", "-e", "set volume output volume ((output volume of (get volume settings)) + " + strconv.Itoa(diff) + ")"}
}

func getMutedCmd(outputdevice string) []string {
	return []string{"osascript", "-e", "output muted of (get volume settings)"}
}

func parseMuted(out string) (bool, error) {
	switch strings.TrimSpace(out) {
	case "true":
		return true, nil
	case "false":
		return false, nil
	}
	return false, fmt.Errorf("unknown muted status: %s", out)
}

func muteCmd(outputdevice string) []string {
	return []string{"osascript", "-e", "set volume output muted true"}
}

func unmuteCmd(outputdevice string) []string {
	return []string{"osascript", "-e", "set volume output muted false"}
}
