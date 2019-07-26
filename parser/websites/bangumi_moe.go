package websites

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type bangumiMoe struct {
	WebBase
}

func (web *bangumiMoe) Request(args []string) error {
	argsStr := strings.Join(args, " ")
	jsonStr := fmt.Sprintf(`{"query": "%s"}`, argsStr)
	req, err := http.NewRequest("POST", web.urlBase, bytes.NewBufferString(jsonStr))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	statuscode := res.StatusCode
	if statuscode != 200 {
		return fmt.Errorf("status code is %d", statuscode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var jsonRes map[string]interface{}
	json.Unmarshal(body, &jsonRes)

	infos := make([]map[string]string, 0)

	torrents := jsonRes["torrents"].([]interface{})
	for idx, torrent := range torrents {
		info := make(map[string]string)
		info["no"] = strconv.Itoa(idx)

		torrentMap := torrent.(map[string]interface{})

		t, _ := time.Parse(time.RFC3339, torrentMap["publish_time"].(string))
		info["date"] = t.Format("2006-01-02 15:04")

		tmp := torrentMap["category_tag"].(map[string]interface{})["synonyms"].([]interface{})
		if len(tmp) > 0 {
			info["type"] = tmp[0].(string)
		} else {
			info["type"] = "-"
		}
		info["title"] = torrentMap["title"].(string)
		info["team"] = torrentMap["team"].(map[string]interface{})["name"].(string)
		info["magnet"] = torrentMap["magnet"].(string)
		info["size"] = torrentMap["size"].(string)
		info["publisher"] = torrentMap["uploader"].(map[string]interface{})["username"].(string)
		info["torrentSourceNum"] = strconv.Itoa(int(torrentMap["seeders"].(float64)))
		info["downloadedNum"] = strconv.Itoa(int(torrentMap["downloads"].(float64)))
		info["finishedNum"] = strconv.Itoa(int(torrentMap["finished"].(float64)))

		infos = append(infos, info)
	}

	web.resCache = infos
	return nil
}

func BangumiMoeCtor() WebParser {
	b := &bangumiMoe{}
	b.urlBase = "https://bangumi.moe/api/v2/torrent/search"
	return b
}
