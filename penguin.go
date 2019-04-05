package main

import (
	_ "net/http/pprof"

	"github.com/liipx/penguin/ui"
)

func main() {
	ui.Run()
}
