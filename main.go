package main

import (
	"flag"
	"fmt"
	"github.com/docker/machine/libmachine/drivers/plugin"
	"github.com/yandex-cloud/docker-machine-driver-yandex/driver"
	"os"
)

// Version will be added once we start the build process
var Version string

func main() {
	version := flag.Bool("v", false, "prints current docker-machine-driver-yandex version")
	flag.Parse()
	if *version {
		fmt.Printf("Version: %s\n", Version)
		os.Exit(0)
	}
	plugin.RegisterDriver(driver.NewDriver())
}
