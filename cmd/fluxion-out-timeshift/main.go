package main

import (
	"github.com/yosisa/fluxion-timeshift"
	"github.com/yosisa/fluxion/plugin"
)

func main() {
	plugin.New("out-timeshift", timeshift.Factory).Run()
}
