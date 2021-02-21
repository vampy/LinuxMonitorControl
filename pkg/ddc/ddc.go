package ddc

// https://en.wikipedia.org/wiki/Display_Data_Channel

import (
	"bufio"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/vampy/LinuxMonitorControl/pkg/xorg"
)

/*
#cgo CFLAGS: -I./ddcutil/src/src/
#cgo LDFLAGS: -L./ddcutil/lib -lddcutil -Wl,-rpath=./ddcutil/lib

#include "ddcutil_c_api.h"
*/

// https://www.ddcutil.com/command_getvcp/
type Value struct {
	FeatureCode  string
	CurrentValue int
	MaxValue     int
}

type Display struct {
	Index  int
	I2CBus string
	Name   string
}

type DDC struct {
	// Map from Index to Display
	displays map[int]*Display
}

// http://www.ddcutil.com/command_capabilities/
const featureCodeBrightness = 10
const featureCodeContrast = 12

// We depend on https://www.ddcutil.com/
func getBinaryPath() string {
	local := "./ddcutil/bin/ddcutil"
	_, err := exec.LookPath(local)
	if err == nil {
		return local
	}

	system := "ddcutil"
	_, err = exec.LookPath(system)
	if err == nil {
		return system
	}

	return system
}

func runCommand(name string, arg ...string) (string, error) {
	cmd := exec.Command(name, arg...)
	fmt.Printf("Running command: %s\n", cmd)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

func (d *DDC) String() string {
	var sb strings.Builder
	for key, value := range d.displays {
		sb.WriteString(fmt.Sprintf("%v => %v, ", key, value))
	}

	return fmt.Sprintf("Indices = %v, Displays = {%+v}", d.Indices, sb.String())
}

func New() (*DDC, error) {
	// cVersion := C.ddca_ddcutil_version_string()
	// version := C.GoString(cVersion)
	// fmt.Println(version)

	ddcutilPath := getBinaryPath()
	_, err := exec.LookPath(ddcutilPath)
	if err != nil {
		return nil, err
	}

	output, err := runCommand(ddcutilPath, "detect", "--nousb", "--async", "--brief")
	if err != nil {
		return nil, err
	}

	// Scan line by line
	ddc := &DDC{}
	ddc.displays = make(map[int]*Display)
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	var curDisplay *Display = nil
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// fmt.Println("line: " + line)

		// Add new Display
		if len(line) == 0 {
			if curDisplay != nil {
				ddc.displays[curDisplay.Index] = curDisplay
			}

			curDisplay = nil
			continue
		}

		// Process line
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		key := fields[0]
		value := fields[1]
		switch strings.ToLower(key) {
		case "display":
			index, err := strconv.Atoi(value)
			if err != nil {
				return nil, err
			}
			curDisplay = &Display{
				Index: index,
			}
		case "i2c bus:":
			if curDisplay == nil {
				continue
			}
			curDisplay.I2CBus = value

		case "monitor:":
			if curDisplay == nil {
				continue
			}

			// Value can have whitespaces
			value = strings.Join(fields[1:len(fields)-1], " ")
			curDisplay.Name = value
		}
	}

	return ddc, nil
}

func (d *DDC) Indices() []int {
	indices := make([]int, len(d.Displays()))
	for key, _ := range d.displays {
		indices = append(indices, key)
	}
	return indices
}

func (d *DDC) Displays() map[int]*Display {
	return d.displays
}

func (d *DDC) SetBrightness(displayIndex int, brightness int) error {
	_, ok := d.displays[displayIndex]
	if ok {
		ddcutilPath := getBinaryPath()
		_, err := exec.LookPath(ddcutilPath)
		if err != nil {
			return err
		}

		_, err = runCommand(
			ddcutilPath,
			"--display", strconv.Itoa(displayIndex),
			"setvcp", strconv.Itoa(featureCodeBrightness), strconv.Itoa(brightness))
		if err != nil {
			return err
		}
	} else {
		// Try with xorg xxbbacklight
		lighter, err := xorg.NewBacklight()
		if err != nil {
			return err
		}
		lighter.SetAll(brightness)
		// err = x.Set(float64(brightness))
		// if err != nil {
		// 	return err
		// }
	}

	return nil
}

// def get_laptop_display_brightness():
//     return int(float(run_command("xbacklight -get").strip()))

// def set_ddc_display_brightness(display_index, brightness):
//     # For the rest we use ddc https://www.ddcutil.com/
//     run_command(f"ddcutil --display {display_index} setvcp 10 {brightness}")
