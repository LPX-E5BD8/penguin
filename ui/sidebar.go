package ui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/jroimartin/gocui"
	"github.com/liipx/penguin/model"
)

// sidebar items
const sign = " * "

var sidebarContent = []string{
	sign + model.Version55,
	sign + model.Version56,
	sign + model.Version57,
	sign + model.Version80,
}

// key bindings for sidebar
func sidebarKeyBinding(g *gocui.Gui) error {
	if err := g.SetKeybinding(sidebarView, gocui.KeyEnter, gocui.ModNone, selectVersion); err != nil {
		return err
	}
	if err := g.SetKeybinding(optionView, gocui.KeyEnter, gocui.ModNone, selectOption); err != nil {
		return err
	}

	return nil
}

// select from version view
func selectVersion(g *gocui.Gui, v *gocui.View) error {
	currentVersion = strings.Replace(getLine(v, 0), sign, "", -1)
	go g.Update(func(gui *gocui.Gui) error {
		v, err := g.View(optionView)
		if err != nil {
			return err
		}
		v.Clear()
		v.Title = "Tags for MySQL " + currentVersion
		tmp := make([]string, 0)
		for tag := range options[currentVersion] {
			tmp = append(tmp, tag)
		}
		sort.Strings(tmp)
		for _, tag := range tmp {
			_, _ = fmt.Fprintln(v, sign+tag)
		}
		return nil
	})

	go g.Update(func(gui *gocui.Gui) error {
		v, err := g.View(contentView)
		if err != nil {
			return err
		}
		v.Clear()
		_, _ = fmt.Fprintln(v, releaseCache[currentVersion])
		return nil
	})

	return nil
}

// select from option view
func selectOption(g *gocui.Gui, v *gocui.View) error {
	tag := strings.Replace(getLine(v, 0), sign, "", -1)
	go g.Update(func(gui *gocui.Gui) error {
		v, err := g.View(contentView)
		if err != nil {
			return err
		}
		v.Clear()
		v.Wrap = true
		x, _ := v.Size()
		for _, item := range options[currentVersion][tag] {
			_, _ = fmt.Fprintln(v, "Version :", item.ReleaseNote.Version)
			_, _ = fmt.Fprintln(v, item.Content)

			if len(item.RelatedBugs) > 0 {
				_, _ = fmt.Fprintln(v, " ")
			}

			for _, bug := range item.RelatedBugs {
				_, _ = fmt.Fprintln(v,
					fmt.Sprintf("BUG #%d [%s]: %s", bug.ID, bug.Title, bug.URL),
				)
			}

			_, _ = fmt.Fprintln(v, strings.Repeat("-", x))
		}

		return nil
	})

	return nil
}
