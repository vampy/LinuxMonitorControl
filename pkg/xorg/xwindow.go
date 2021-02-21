// MIT License

// Copyright (c) 2016 kbinani

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
// https://github.com/kbinani/screenshot/blob/master/internal/xwindow/xwindow.go

package xorg

import (
	"image"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xinerama"
)

func NumActiveDisplays() (num int) {
	defer func() {
		e := recover()
		if e != nil {
			num = 0
		}
	}()

	c, err := xgb.NewConn()
	if err != nil {
		return 0
	}
	defer c.Close()

	err = xinerama.Init(c)
	if err != nil {
		return 0
	}

	reply, err := xinerama.QueryScreens(c).Reply()
	if err != nil {
		return 0
	}

	num = int(reply.Number)
	return num
}

func GetDisplayBounds(displayIndex int) (rect image.Rectangle) {
	defer func() {
		e := recover()
		if e != nil {
			rect = image.ZR
		}
	}()

	c, err := xgb.NewConn()
	if err != nil {
		return image.ZR
	}
	defer c.Close()

	err = xinerama.Init(c)
	if err != nil {
		return image.ZR
	}

	reply, err := xinerama.QueryScreens(c).Reply()
	if err != nil {
		return image.ZR
	}

	if displayIndex >= int(reply.Number) {
		return image.ZR
	}

	primary := reply.ScreenInfo[0]
	x0 := int(primary.XOrg)
	y0 := int(primary.YOrg)

	screen := reply.ScreenInfo[displayIndex]
	x := int(screen.XOrg) - x0
	y := int(screen.YOrg) - y0
	w := int(screen.Width)
	h := int(screen.Height)
	rect = image.Rect(x, y, x+w, y+h)
	return rect
}
