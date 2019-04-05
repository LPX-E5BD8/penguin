package ui

import (
	"github.com/jroimartin/gocui"
)

func contentKeyBindings(g *gocui.Gui) error {
	if err := g.SetKeybinding(contentView, gocui.KeyArrowUp, gocui.ModNone, pageUp); err != nil {
		return err
	}
	if err := g.SetKeybinding(contentView, gocui.KeyArrowDown, gocui.ModNone, pageDown); err != nil {
		return err
	}
	if err := g.SetKeybinding(contentView, gocui.KeyPgup, gocui.ModNone, pageUp); err != nil {
		return err
	}
	if err := g.SetKeybinding(contentView, gocui.KeyPgdn, gocui.ModNone, pageDown); err != nil {
		return err
	}
	return nil
}
