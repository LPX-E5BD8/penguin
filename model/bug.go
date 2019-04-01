package model

import (
	"bytes"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/kr/pretty"
)

type Bug struct {
	ID           int
	URL          string
	Title        string
	SubmitTime   time.Time
	ModifiedTime time.Time
	Reporter     string
	Status       string
	Category     string
	Version      []string
	Tags         []string
	Triage       string
	Severity     string
	OS           string
	CPUArch      []string
}

const bugApiTemplate = "https://bugs.mysql.com/bug.php?id=%d"
const timeFMT = "2 Jan 2006 15:04"

// New prepare a new bug struct
func (bug *Bug) New(id int) *Bug {
	bug.ID = id
	bug.URL = fmt.Sprintf(bugApiTemplate, id)
	return bug
}

// Analysis prepare a new bug struct
func (bug *Bug) Analysis(wg *sync.WaitGroup) *Bug {
	wg.Add(1)
	go func() {
		defer wg.Done()
		res, err := HTTPGetWithCache(bug.URL, CacheDir)
		if err != nil {
			Logger.Println("bug analysis get url: ", bug.URL, err)
			return
		}

		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(res))
		doc.Find("#bugheader").Find("td").Each(func(i int, selection *goquery.Selection) {
			value := compressStr(strings.TrimSpace(selection.Text()))
			switch i {
			case 0:
				bug.Title = value
			case 1:
				subTime, _ := time.Parse(timeFMT, value)
				bug.SubmitTime = subTime
			case 2:
				mTime, _ := time.Parse(timeFMT, value)
				bug.ModifiedTime = mTime
			case 3:
				bug.Reporter = value
			case 5:
				bug.Status = value
			case 7:
				bug.Category = value
			case 8:
				bug.Severity = value
			case 9:
				bug.Version = strings.Split(value, ",")
			case 10:
				bug.OS = value
			case 12:
				bug.CPUArch = strings.Split(value, ",")
			case 13:
				if strings.Index(value, "Triage") >= 0 {
					bug.Triage = value
					break
				}
				bug.Tags = strings.Split(value, ",")
			case 14:
				bug.Triage = value
			}
		})
		pretty.Println(bug)
	}()
	return bug
}
