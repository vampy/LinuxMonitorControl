package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/vampy/LinuxMonitorControl/pkg/build"
	"github.com/vampy/LinuxMonitorControl/pkg/ddc"
	"github.com/vampy/LinuxMonitorControl/pkg/xorg"
	"golang.org/x/sync/errgroup"
)

// TODO systray https://www.reddit.com/r/golang/comments/bh0p2h/go_and_the_linux_system_tray/

func main() {
	app := cli.NewApp()
	app.Version = build.Version()

	app.Flags = []cli.Flag{
		&cli.IntFlag{
			Name:    "brightness",
			Value:   -1,
			Usage:   "brightness value 0-100",
			Aliases: []string{"b"},
		},
		&cli.IntFlag{
			Name:    "display",
			Value:   -1,
			Usage:   "display index to control. -1 means all displays",
			Aliases: []string{"d"},
		},
		&cli.IntFlag{
			Name:    "contrast",
			Value:   -1,
			Usage:   "contrast value 0-100",
			Aliases: []string{"c"},
		},
	}

	app.Before = func(c *cli.Context) error {
		numDisplays := xorg.NumActiveDisplays()
		fmt.Fprintf(c.App.Writer, "NumDisplays = %d\n\n", numDisplays)
		return nil
	}

	app.Action = func(c *cli.Context) error {
		ddc, err := ddc.New()
		if err != nil {
			return err
		}
		fmt.Printf("%v\n\n", ddc)

		brightness := c.Int("brightness")
		contrast := c.Int("contrast")
		display := c.Int("display")
		if brightness >= 0 && brightness <= 100 {
			fmt.Printf("Using brightness = %d\n\n", brightness)

			if display == -1 {
				// All
				numDisplays := xorg.NumActiveDisplays()

				errs, _ := errgroup.WithContext(context.Background())
				for i := 0; i < numDisplays; i++ {
					i := i
					errs.Go(func() error {
						return ddc.SetBrightness(i, brightness)
					})
				}

				err := errs.Wait()
				if err != nil {
					return err
				}
			} else {
				// One Display
				fmt.Printf("For Display = %d\n", display)
				err := ddc.SetBrightness(display, brightness)
				if err != nil {
					return err
				}
			}
		}

		if contrast >= 0 && contrast <= 100 {
			fmt.Printf("Using contrast = %d\n\n", contrast)
			if display != -1 {
				fmt.Printf("For Display = %d\n", display)
				err := ddc.SetContrast(display, contrast)
				if err != nil {
					return err
				}
			}
		}

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
