package main

import (
	"fmt"
	"github.com/docker/machine/libmachine/mcnflag"
	"github.com/yandex-cloud/docker-machine-driver-yandex/driver"
)

func main() {
	d := &driver.Driver{}
	for _, f := range d.GetCreateFlags() {
		cmdFlag := fmt.Sprintf("--%s", f.String())
		defValue := f.Default()

		envVar, desc := "", ""
		switch v := f.(type) {
		case mcnflag.StringFlag:
			envVar = v.EnvVar
			defValue = v.Default()
			desc = v.Usage
		case mcnflag.BoolFlag:
			envVar = v.EnvVar
			defValue = false
			desc = v.Usage
		case mcnflag.IntFlag:
			envVar = v.EnvVar
			defValue = v.Default()
			desc = v.Usage
		case mcnflag.StringSliceFlag:
			envVar = v.EnvVar
			defValue = v.Default()
			desc = v.Usage
		}

		combinedOne := fmt.Sprintf("``%s`` or ``$%s``", cmdFlag, envVar)
		fmt.Printf("| %-55s | %-49s | %-24v | | \n", combinedOne, desc, defValue)
	}
}
