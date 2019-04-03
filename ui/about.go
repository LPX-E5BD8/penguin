package ui

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

// about content
var aboutContent = []string{
	"Penguin, MySQL Doc Tool.",
	"Powered by gocui.",
	"Press F12 for more detials.",
}

// about details content
var aboutDetailsContent = []string{
	"  Author: Liipx",
	"    Mail: lipengxiang_@outlook.com",
	" -----------------------------------------",
	"  Github: https://github.com/liipx/penguin",
	" Version: 0.1.0",
	" ",
	"                       Press Enter to quit",
}

// key bindings for about view
func aboutKeyBinding(g *gocui.Gui) error {
	// del about details
	if err := g.SetKeybinding(aboutDetails, gocui.KeyEnter, gocui.ModNone, delAbout); err != nil {
		return err
	}

	return nil
}

func showAboutDetails(g *gocui.Gui, v *gocui.View) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView(aboutDetails, maxX/2-22, maxY/2-3, maxX/2+22, maxY/2+5); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = "About"
		for _, detail := range aboutDetailsContent {
			if _, err := fmt.Fprintln(v, detail); err != nil {
				continue
			}
		}

		if _, err := setCurrentViewOnTop(g, aboutDetails); err != nil {
			return err
		}
	}
	return nil
}

func delAbout(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView(aboutDetails); err != nil {
		return err
	}

	if _, err := setCurrentViewOnTop(g, lastViewName); err != nil {
		return err
	}

	return nil
}
