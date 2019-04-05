package ui

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

func Run() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Highlight = true
	g.SelFgColor = gocui.ColorGreen

	g.SetManagerFunc(Layout)

	if err := keyBindings(g); err != nil {
		fmt.Println(err)
		return
	}

	go cacheBuild(g)

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		fmt.Println(err)
		return
	}
}
