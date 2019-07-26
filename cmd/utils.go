package cmd

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/matzoe/argo/rpc"
	"github.com/spf13/viper"
)

func parseEpisodeFromTitle(title string) float64 {
	blacklistPatterns := []string{
		`x?(1080|720|480)p?`,
		`\d+\s*月新番`,
		`x26[45]`,
		`10bit`,
		`ma10p`,
		`\bv\d`,
		`big5`,
		`mp4`,
		`assx2`,
	}

	epNumPattern := `(?:\D*|^)([\d.]+)(-[\d.]+)?(?:\D*|$)`

	rawTokens := regexp.MustCompile(`[[\]【】_\s]`).Split(title, -1)
	tokens := []string{}
	for _, token := range rawTokens {
		token = strings.ToLower(token)
		if matched, _ := regexp.MatchString(`\d`, token); matched {
			matched := false
			for _, pattern := range blacklistPatterns {
				if m, _ := regexp.MatchString(pattern, token); m {
					matched = true
					break
				}
			}
			if !matched {
				token = strings.TrimSpace(token)
				tokens = append(tokens, token)
			}
		}
	}

	parseEpisode := func(token string) float64 {
		re := regexp.MustCompile(epNumPattern)
		s := re.ReplaceAllString(token, `$1$2`)
		s = strings.Split(s, "-")[0]
		ans, err := strconv.ParseFloat(string(s), 64)
		if err != nil {
			fmt.Println("Cannot parseEpisode:", err, s)
			return -1
		}
		return ans
	}

	for idx := len(tokens) - 1; idx >= 0; idx-- {
		token := tokens[idx]
		token = regexp.MustCompile(`\s*(end|完)$`).ReplaceAllString(token, ``)
		token = regexp.MustCompile(`\s*v\d+$`).ReplaceAllString(token, ``)

		if matched, _ := regexp.MatchString(epNumPattern, token); matched {
			return parseEpisode(token)
		}
	}

	fmt.Println("Unable to parse episode:", title)
	return -1
}

func trimMagnets(magnets []string) {
	for idx, s := range magnets {
		if strings.Index(s, "&") != -1 {
			magnets[idx] = s[:strings.Index(s, "&")]
		}
	}
}

func downloadMagnets(magnets []string, dir string) error {
	address := viper.GetString("aria2-rpc-address")
	token := viper.GetString("aria2-rpc-token")
	client, err := rpc.New(context.Background(), address, token, 10*time.Second, rpc.DummyNotifier{})
	if err != nil {
		return err
	}

	downDir := map[string]interface{}{"dir": dir}
	defer client.Close()
	for _, magnet := range magnets {
		_, err := client.AddURI(magnet, downDir)
		if err != nil {
			fmt.Println("Failed to download magnet link:", err)
			fmt.Printf("\t%s\n", magnet)
		}
	}
	return nil
}
