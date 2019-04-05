package ui

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

func Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if maxY < 48 {
		return fmt.Errorf("y is too small:%d", maxY)
	}
	// set content
	if v, err := g.SetView(titleView, 0, 0, 30, 10); err != nil {
		if err != gocui.ErrUnknownView || v == nil {
			return err
		}

		v.Wrap = true
		v.Title = "Usage"
		for _, content := range usageContent {
			if _, err := fmt.Fprintln(v, content); err != nil {
				return err
			}
		}
	}

	// set sidebar
	if v, err := g.SetView(sidebarView, 0, 11, 30, maxY-35-1); err != nil {
		if err != gocui.ErrUnknownView || v == nil {
			return err
		}

		v.Wrap = true
		v.Highlight = true
		v.Title = "MySQL Change Note"

		// set content
		for _, content := range sidebarContent {
			if _, err := fmt.Fprintln(v, content); err != nil {
				continue
			}
		}

		if _, err := setCurrentViewOnTop(g, sidebarView); err != nil {
			return err
		}
	}

	// set options
	if v, err := g.SetView(optionView, 0, maxY-35, 30, maxY-5-1); err != nil {
		if err != gocui.ErrUnknownView || v == nil {
			return err
		}

		v.Title = "Options"
		v.Highlight = true
	}

	// set about
	if v, err := g.SetView(aboutView, 0, maxY-5, 30, maxY-1); err != nil {
		if err != gocui.ErrUnknownView || v == nil {
			return err
		}

		v.Wrap = true
		v.Title = "About Penguin"
		// set content
		for _, content := range aboutContent {
			if _, err := fmt.Fprintln(v, content); err != nil {
				continue
			}
		}
	}

	// set main
	if v, err := g.SetView(contentView, 31, 0, maxX-1, maxY-5-1); err != nil {
		if err != gocui.ErrUnknownView || v == nil {
			return err
		}
	}

	// set log view
	if v, err := g.SetView(loggerView, 31, maxY-5, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView || v == nil {
			return err
		}

		v.Autoscroll = true
		v.Title = "Log"
	}

	return nil
}

func keyBindings(g *gocui.Gui) error {
	// global key bind
	if err := globalKeyBindings(g); err != nil {
		return err
	}

	// sidebar key bind
	if err := sidebarKeyBindings(g); err != nil {
		return err
	}

	// content key bind
	if err := contentKeyBindings(g); err != nil {
		return err
	}

	return nil
}

func globalKeyBindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}

	for _, viewName := range viewArr {
		if err := g.SetKeybinding(viewName, gocui.KeyTab, gocui.ModNone, nextView); err != nil {
			return err
		}
		// show about details
		if err := g.SetKeybinding(viewName, gocui.KeyF12, gocui.ModNone, showAboutDetails); err != nil {
			return err
		}
	}

	// about view key binding
	if err := aboutKeyBinding(g); err != nil {
		return err
	}

	return nil
}
