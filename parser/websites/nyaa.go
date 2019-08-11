package websites

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type nyaa struct {
	WebBase
}

func (web *nyaa) Request(args []string) error {
	for idx, s := range args {
		args[idx] = regexp.MustCompile(`\s`).ReplaceAllString(s, `+`)
	}
	url := web.urlBase + strings.Join(args, `+`)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return err
	}

	infos := make([]map[string]string, 0)
	count := 0
	doc.Find("table.torrent-list").Find("tbody").Find("tr").Each(func(idx int, tr *goquery.Selection) {
		info := make(map[string]string)
		info["no"] = strconv.Itoa(count)
		count++
		tr.Find("td").Each(func(idx int, tablecell *goquery.Selection) {
			switch idx {
			case 0:
				t, ok := tablecell.Find("a").Attr("title")
				if ok {
					info["type"] = strings.TrimSpace(t)
				}
			case 1:
				tablecell.Find("a").Each(func(idx int, a *goquery.Selection) {
					if !a.HasClass("comments") {
						t, ok := a.Attr("title")
						if ok {
							info["title"] = strings.TrimSpace(t)
						}
					}
				})
			case 2:
				tablecell.Find("a").Each(func(idx int, a *goquery.Selection) {
					if idx == 1 {
						href, ok := a.Attr("href")
						if ok {
							info["magnet"] = strings.TrimSpace(href)
						}
					}
				})
			case 3:
				info["size"] = strings.TrimSpace(tablecell.Text())
			case 4:
				info["date"] = strings.TrimSpace(tablecell.Text())
			case 5:
				info["seeders"] = strings.TrimSpace(tablecell.Text())
			case 6:
				info["leechers"] = strings.TrimSpace(tablecell.Text())
			case 7:
				info["finishedNum"] = strings.TrimSpace(tablecell.Text())
			}
		})
		infos = append(infos, info)
	})

	web.resCache = infos

	return nil
}

func NyaaCtor() WebParser {
	nyaa := &nyaa{}
	nyaa.urlBase = "https://nyaa.si/?f=0&c=0_0&q="
	return nyaa
}
