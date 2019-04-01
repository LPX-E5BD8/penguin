package model

import (
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

type Link struct {
	Title string
	URL   string
}

type ReleaseNoteItem struct {
	Class        string
	Content      string
	Tags         []string
	RelatedLinks []*Link
	RelatedBugs  []*Bug
}

// Analysis ReleaseNoteItem from *goquery.Selection
func (item *ReleaseNoteItem) Analysis(class string, selection *goquery.Selection, wg *sync.WaitGroup) *ReleaseNoteItem {
	wg.Add(1)
	go func() {
		defer wg.Done()
		item.Class = class
		item.Tags = analysisTags(selection)
		item.RelatedLinks = analysisLinks(selection)
		item.Content = strings.TrimSpace(selection.Text())
		item.RelatedBugs = analysisBugs(item.Content)
	}()
	return item
}

func analysisLinks(selection *goquery.Selection) []*Link {
	links := make([]*Link, 0)
	selection.Find("a.ulink").Each(func(i int, selection *goquery.Selection) {
		u, _ := selection.Attr("href")
		links = append(links, &Link{
			Title: compressStr(strings.TrimSpace(selection.Text())),
			URL:   baseURL + u,
		})
	})
	return links
}

func analysisTags(selection *goquery.Selection) []string {
	tagStr := selection.Find("span.bold").Find("strong").Text()
	return strings.Split(strings.Trim(tagStr, ":"), ";")
}

func analysisBugs(content string) []*Bug {
	regExp, _ := regexp.Compile("Bug #[\\d]+")
	bugStr := regExp.FindAllString(content, -1)

	wg := new(sync.WaitGroup)
	bugs := make([]*Bug, 0)
	for _, bug := range bugStr {
		info := strings.Split(bug, "#")
		if len(info) != 2 {
			Logger.Println("unknown bug: ", bugStr)
			continue
		}

		bugId, err := strconv.Atoi(info[1])
		if err != nil {
			Logger.Println("unknown bug id: ", bugId, info)
			continue
		}

		if bugId > 10000000 {
			continue
		}

		bugs = append(bugs, new(Bug).New(bugId).Analysis(wg))
	}

	wg.Wait()
	return bugs
}
