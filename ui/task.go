package ui

import (
	"github.com/jroimartin/gocui"
	"github.com/liipx/penguin/model"
)

var options = make(map[string]map[string][]*model.ReleaseNoteItem)

// releaseCache cache struct
var releaseCache = map[string]*model.ReleaseInfo{}

func cacheBuild(g *gocui.Gui) {
	for _, ver := range []string{model.Version55, model.Version56, model.Version57, model.Version80} {
		viewLogPrintln(g, "[Info]", "Loading mysql", ver, "docs ...")
		info, err := model.NewReleaseInfo(ver)
		if err != nil {
			viewLogPrintln(g, "[Error]", err)
			continue
		}

		cacheOptions(info)
		releaseCache[ver] = info
		viewLogPrintln(g, "[Done] Mysql", ver, "doc loaded.")
	}
	viewLogPrintln(g, "[Done] All version cached.")
}

// cacheOptions
func cacheOptions(info *model.ReleaseInfo) {
	options[info.Version] = make(map[string][]*model.ReleaseNoteItem)
	tmpMap := map[string]struct{}{}
	for _, note := range info.Info {
		for _, item := range note.Items {
			for _, tag := range item.Tags {
				if tag == "" {
					tag = "No Tags"
				}

				if _, ok := options[info.Version][tag]; !ok {
					options[info.Version][tag] = make([]*model.ReleaseNoteItem, 0)
				}

				if _, ok := tmpMap[tag+item.Content]; !ok {
					options[info.Version][tag] = append(options[info.Version][tag], item)
					tmpMap[tag+item.Content] = struct{}{}
				}
			}
		}
	}
}
