package websites

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type dmhy struct {
	WebBase
}

func (web *dmhy) Request(args []string) error {
	for idx, s := range args {
		args[idx] = regexp.MustCompile(`\s`).ReplaceAllString(s, `+`)
	}
	url := web.urlBase + strings.Join(args, "+")
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return err
	}

	infos := make([]map[string]string, 0)
	count := 0
	doc.Find("table.tablesorter").Find("tbody").Find("tr").Each(func(idx int, tr *goquery.Selection) {
		info := make(map[string]string)
		info["no"] = strconv.Itoa(count)
		count++
		tr.Find("td").Each(func(idx int, tablecell *goquery.Selection) {
			switch idx {
			case 0:
				info["date"] = strings.TrimSpace(tablecell.Find("span").Text())
			case 1:
				info["type"] = strings.TrimSpace(tablecell.Text())
			case 2:
				tablecell.Children().Each(func(idx int, item *goquery.Selection) {
					if item.Is("span") && item.HasClass("tag") {
						info["organization"] = strings.TrimSpace(item.Find("a").Text())
					} else if item.Is("a") {
						info["title"] = strings.TrimSpace(item.Text())
					}
				})
			case 3:
				href, ok := tablecell.Find("a.arrow-magnet").Attr("href")
				if ok {
					info["magnet"] = strings.TrimSpace(href)
				}
			case 4:
				info["size"] = strings.TrimSpace(tablecell.Text())
			case 5:
				info["torrentSourceNum"] = strings.TrimSpace(tablecell.Text())
			case 6:
				info["downloadedNum"] = strings.TrimSpace(tablecell.Text())
			case 7:
				info["finishedNum"] = strings.TrimSpace(tablecell.Text())
			case 8:
				info["publisher"] = strings.TrimSpace(tablecell.Find("a").Text())
			}
		})
		infos = append(infos, info)
	})
	web.resCache = infos
	return nil
}

func (web *dmhy) ShowFindResult(filterMap map[string]int, len int) [][]string {
	infos := [][]string{}
	for _, row := range web.resCache {
		info := make([]string, len)
		for key, idx := range filterMap {
			info[idx] = row[key]
		}
		infos = append(infos, info)
	}

	return infos
}

func (web *dmhy) GetMagnets(selectNums []int) []string {
	res := make([]string, len(selectNums))
	for idx, num := range selectNums {
		res[idx] = web.resCache[num]["magnet"]
	}
	return res
}

func DmhyCtor() WebParser {
	dmhy := &dmhy{}
	dmhy.urlBase = "https://share.dmhy.org/topics/list?keyword="
	return dmhy
}
