package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/vampy/LinuxMonitorControl/pkg/build"
	"github.com/vampy/LinuxMonitorControl/pkg/ddc"
	"github.com/vampy/LinuxMonitorControl/pkg/xorg"
	"golang.org/x/sync/errgroup"
)

// TODO systray https://www.reddit.com/r/golang/comments/bh0p2h/go_and_the_linux_system_tray/

var brightness = flag.Int("brightness", 0, "brightness of all the monitors")

func main() {

	flag.Parse()
	flagset := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) { flagset[f.Name] = true })

	// if flag.NArg() == 0 {
	// 	flag.Usage()
	// 	os.Exit(1)
	// }

	build.Print()

	ddc, err := ddc.New()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n\n", ddc)

	numDisplays := xorg.NumActiveDisplays()
	fmt.Printf("NumDisplays = %d\n\n", numDisplays)
	if flagset["brightness"] {
		fmt.Printf("Using brightness = %d\n\n", *brightness)

		errs, _ := errgroup.WithContext(context.Background())
		for i := 0; i < numDisplays; i++ {
			i := i
			errs.Go(func() error {
				return ddc.SetBrightness(i, *brightness)
			})
		}

		err := errs.Wait()
		if err != nil {
			panic(err)
		}
	}
}
