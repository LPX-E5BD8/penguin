package model

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// ReleaseInfo all release note for a mysql version
type ReleaseInfo struct {
	Version string
	Info    []*ReleaseNote
}

func (ri ReleaseInfo) String() string {
	t := template.Must(template.New("ReleaseInfo").Parse(ReleaseInfoTemplate))
	buf := bytes.NewBuffer(make([]byte, 0))
	_ = t.Execute(buf, ri)
	return buf.String()
}

const (
	Version55 = "5.5"
	Version56 = "5.6"
	Version57 = "5.7"
	Version80 = "8.0"
)

const baseURL = "https://dev.mysql.com"
const releaseNoteAPITemplate = baseURL + "/doc/relnotes/mysql/%s/en/"

func NewReleaseInfo(version string) (*ReleaseInfo, error) {
	uri := fmt.Sprintf(releaseNoteAPITemplate, version)
	res, err := HTTPGetWithCache(uri, CacheDir)
	if err != nil {
		return nil, err
	}

	ri := &ReleaseInfo{Version: version}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(res))
	if err != nil {
		return nil, err
	}

	doc.Find("span.section").Find("a").Each(func(i int, selection *goquery.Selection) {
		href, _ := selection.Attr("href")
		if rn := new(ReleaseNote).setMeta(selection.Text(), uri+href); rn != nil {
			ri.Info = append(ri.Info, rn)
		}
	})

	wg := &sync.WaitGroup{}
	for _, rn := range ri.Info {
		wg.Add(1)
		err = Pool.Submit(func() {
			func(note *ReleaseNote) {
				defer wg.Done()
				if err := note.Analysis(); err != nil {
					Logger.Println("note.Analysis() err:", err, "ref:", note.URL)
				}
			}(rn)
		})
	}
	wg.Wait()
	return ri, err
}

// ReleaseNote mysql release note
type ReleaseNote struct {
	Version string             // MySQL version to sort
	URL     string             // change note url
	RelType string             // released type: GA/GC/DM
	RelTime time.Time          // released time
	IsRel   bool               // is released
	Items   []*ReleaseNoteItem // change details
}

func (rn *ReleaseNote) setMeta(info, url string) *ReleaseNote {
	var err error
	meta := strings.Split(info, " ")
	if len(meta) < 4 || !strings.EqualFold(meta[0], "Changes") {
		return nil
	}

	rn.URL = url
	rn.Version = meta[3]

	// parse release time
	if len(meta[4]) > 3 {
		timeStr := meta[4][1 : len(meta[4])-1]
		rn.RelTime, err = time.Parse("2006-01-02", timeStr)
		if err == nil {
			rn.IsRel = true
		}
	}

	// parse release type
	relTypeInfo := strings.Split(info, ",")
	if len(relTypeInfo) > 1 {
		rn.RelType = strings.TrimSpace(relTypeInfo[1][:len(relTypeInfo[1])-1])
	}

	return rn
}

// Analysis ReleaseNote from document
func (rn *ReleaseNote) Analysis() error {
	if rn.URL == "" {
		return fmt.Errorf("release note has no URL")
	}

	res, err := HTTPGetWithCache(rn.URL, CacheDir)
	if err != nil {
		return err
	}

	wg := new(sync.WaitGroup)
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(res))
	if err != nil {
		return err
	}

	doc.Find("div.simplesect").Each(func(i int, selection *goquery.Selection) {
		class := strings.TrimSpace(selection.Find("div.titlepage").Text())
		selection.Find("li.listitem").Each(func(i int, selection *goquery.Selection) {
			// Analysis ReleaseNoteItem
			rn.Items = append(rn.Items, new(ReleaseNoteItem).Analysis(rn, class, selection, wg))
		})
	})

	wg.Wait()
	return nil
}
