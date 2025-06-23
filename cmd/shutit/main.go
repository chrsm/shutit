package main

import (
	"flag"
	"log"
	"runtime"

	"bits.chrsm.org/shutit"
	"bits.chrsm.org/shutit/internal/api"
)

var (
	installF bool
)

func main() {
	runtime.LockOSThread()

	flag.BoolVar(&installF, "install", false, "whether to install shutit as a service")
	flag.Parse()

	if installF {
		log.Print("installing service")
		err := api.Install()
		if err != nil {
			log.Fatalf("error installing shutit: %s", err)
		}

		return
	}

	log.Print("shutit running", shutit.BuildDate, shutit.BuildCommit)
	shutit, err := api.NewShutit()
	if err != nil {
		log.Fatalf("error starting shutit: %s", err)
	}

	shutit.Exec()
}
