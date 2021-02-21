package xorg

import (
	"encoding/binary"
	"fmt"
	"math"

	"github.com/BurntSushi/xgb/randr"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xprop"
)

type Backlight struct {
	x *xgbutil.XUtil

	resources *randr.GetScreenResourcesCurrentReply
	atom      xproto.Atom
}

func NewBacklight() (*Backlight, error) {
	x, err := xgbutil.NewConn()
	if err != nil {
		return nil, err
	}
	if err := randr.Init(x.Conn()); err != nil {
		return nil, err
	}

	resources, err := randr.GetScreenResourcesCurrent(x.Conn(), x.RootWin()).Reply()
	if err != nil {
		return nil, err
	}

	atom, err := xprop.Atom(x, "Backlight", false)
	if err != nil {
		return nil, err
	}

	b := &Backlight{
		x:         x,
		resources: resources,
		atom:      atom,
	}

	return b, nil
}

// If the return is -1 then the output is not valid
func (b *Backlight) getRawBrightessForOutput(output randr.Output) (int, error) {
	atomNew, err := xprop.Atom(b.x, "Backlight", false)
	if err != nil {
		return -1, err
	}
	b.atom = atomNew

	var prop *randr.GetOutputPropertyReply
	prop, err = randr.GetOutputProperty(b.x.Conn(), output, atomNew, xproto.AtomNone, 0, 4, false, false).Reply()
	if err != nil {
		// Try with legacy API
		atomLegacy, err := xprop.Atom(b.x, "BACKLIGHT", false)
		if err != nil {
			return -1, err
		}

		prop, err = randr.GetOutputProperty(b.x.Conn(), output, atomLegacy, xproto.AtomNone, 0, 4, false, false).Reply()
		if err != nil {
			return -1, err
		}

		b.atom = atomLegacy
	}

	if prop.Type != xproto.AtomInteger ||
		prop.NumItems != 1 ||
		prop.Format != 32 {
		return -1, fmt.Errorf("Invalid return type for getRawBrightessForOutput")
	}

	return int(binary.LittleEndian.Uint32(prop.Data)), nil
}

func (b *Backlight) getRawRangeForOutput(output randr.Output) (int, int, error) {
	if b.atom == xproto.AtomNone {
		atom, err := xprop.Atom(b.x, "Backlight", false)
		if err != nil {
			return 0, 0, err
		}
		b.atom = atom
	}

	query, err := randr.QueryOutputProperty(b.x.Conn(), output, b.atom).Reply()
	if err != nil {
		return 0, 0, err
	}

	if query.Range && query.Length == 2 {
		return int(query.ValidValues[0]), int(query.ValidValues[1]), nil
	}

	return 0, 0, fmt.Errorf("Invalid return type for getRawRangeForOutput")
}

// Gets first valid value
// func (b *Backlighter) Get() (float64, error) {
// 	rawBacklight, err := backlightGet(b.x, b.output)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return float64(rawBacklight-b.min) / float64(b.max-b.min), nil
// }

func (b *Backlight) SetAll(brightness int) {

	for i := 0; i < len(b.resources.Outputs); i++ {
		output := b.resources.Outputs[i]

		current, err := b.getRawBrightessForOutput(output)

		// Invalid output, continue
		if current == -1 || err != nil {
			continue
		}

		min, max := 0, 0
		min, max, err = b.getRawRangeForOutput(output)
		if err != nil {
			continue
		}
		currentNormalized := float64(current-min) * 100 / float64(max-min)

		// backlight_set
		data := make([]byte, 4)
		newNormalized := uint32(min) + uint32(math.Ceil(float64(brightness)*float64(max-min)/100))
		binary.LittleEndian.PutUint32(data, newNormalized)
		randr.ChangeOutputProperty(b.x.Conn(), output, b.atom, xproto.AtomInteger, 32, xproto.PropModeReplace, 1, data)
		fmt.Printf("Xorg: Output = %d, Current = %d, New = %d\n", i, int(currentNormalized), int(newNormalized))
	}
}
