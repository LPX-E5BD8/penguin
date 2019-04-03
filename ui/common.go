package ui

import (
	"fmt"

	"github.com/jroimartin/gocui"
	"github.com/liipx/penguin/model"
)

const (
	titleView    = "HeaderView"
	sidebarView  = "SidebarView"
	optionView   = "OptionView"
	contentView  = "LoggerView"
	aboutView    = "AboutView"
	aboutDetails = "AboutDetails"
	loggerView   = "ContentView"
)

var (
	// views can visited
	viewArr        = []string{sidebarView, optionView, contentView}
	currentViewIdx = 0
	currentVersion = model.Version55
	lastViewName   = sidebarView
)

var (
	// usage items
	usageContent = []string{
		"← ↑: Move cursor",
		"Enter: Confirm selection",
		"Tab: Change view focus",
		"^C: Exit",
	}
)

// check if view has net line
func hasNextLine(v *gocui.View) bool {
	return getLine(v, 1) != ""
}

// get view current line content
func getLine(v *gocui.View, offset int) string {
	var l string
	var err error

	_, cy := v.Cursor()
	if l, err = v.Line(cy + offset); err != nil {
		l = ""
	}

	return l
}

// move cursor down
func cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		if !hasNextLine(v) {
			return nil
		}

		if err := v.SetCursor(cx, cy+1); err != nil {
			ox, oy := v.Origin()
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}

	}
	return nil
}

// move cursor up
func cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	return nil
}

// change view in focus
func nextView(g *gocui.Gui, v *gocui.View) error {
	nextIndex := (currentViewIdx + 1) % len(viewArr)
	name := viewArr[nextIndex]

	if _, err := setCurrentViewOnTop(g, name); err != nil {
		return err
	}

	currentViewIdx = nextIndex
	lastViewName = viewArr[currentViewIdx]

	return nil
}

// move view on top level
func setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	return g.SetViewOnTop(name)
}

// exit main loop
func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

// log print
func viewLogPrintln(g *gocui.Gui, a ...interface{}) {
	g.Update(func(gui *gocui.Gui) error {
		v, err := g.View(loggerView)
		if err != nil {
			model.Logger.Println(err)
			return err
		}
		_, _ = fmt.Fprintln(v, a...)
		return nil
	})
}
