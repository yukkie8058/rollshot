package main

import (
	"flag"
	"os"
	"runtime/pprof"

	"fyne.io/fyne/v2/app"
	"github.com/yukkie8058/rollshot/data"
	"github.com/yukkie8058/rollshot/internal"
)

func main() {
	var (
		profile = flag.Bool("profile", false, "Enable CPU profiling")
	)
	flag.Parse()

	if *profile {
		f, err := os.Create("cpu.pprof")
		if err != nil {
			panic(err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			panic(err)
		}
		defer pprof.StopCPUProfile()
	}

	a := app.New()
	internal.ShowMainWindow(a, data.NewImageList())
	a.Run()
}
