package cmd

import (
	"fmt"
	"math"
	"os"
	"path"
	"strconv"
	"sync"
	"time"

	"github.com/myzWILLmake/bgmgo/parser"
	"github.com/spf13/viper"

	"github.com/gen2brain/beeep"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update your subscriptions and download.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		update()
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

func notifyUpdate(updatedSubItems chan *SubItem) {
	title := fmt.Sprintf("Bgmgo Updated %d Subcription(s)", len(updatedSubItems))
	msg := ""
	for {
		if subItem, ok := <-updatedSubItems; ok {
			if subItem.Progress == math.Floor(subItem.Progress) {
				msg += fmt.Sprintf("%s updated to EP.%d.\n", subItem.Name, int(subItem.Progress))
			} else {
				msg += fmt.Sprintf("%s updated to EP.%.1f.\n", subItem.Name, subItem.Progress)
			}
		} else {
			break
		}
	}
	beeep.Notify(title, msg, "")
}

func doUpdateSubItem(subItem *SubItem, wg *sync.WaitGroup, updatedSubItems chan *SubItem) {
	defer wg.Done()

	source := subItem.Source
	if parser.ParserCtor[source] == nil {
		fmt.Println("Not valid data source: ", source)
		return
	}
	web := parser.ParserCtor[source]()

	maxEp := subItem.Progress
	err := web.Request([]string{subItem.Pattern})
	if err != nil {
		fmt.Println("Error: failed to request", source)
		fmt.Println(err)
		return
	}
	filterMap := map[string]int{
		"no":    0,
		"title": 1,
	}
	infos := web.ShowFindResult(filterMap, 2)

	newEpNumMap := make(map[float64]bool)
	selectNums := []int{}
	for _, info := range infos {
		title := info[len(info)-1]
		ep := parseEpisodeFromTitle(title)
		if ep > subItem.Progress && !newEpNumMap[ep] {
			newEpNumMap[ep] = true
			if ep > maxEp {
				maxEp = ep
			}
			idx, err := strconv.Atoi(info[0])
			if err == nil {
				selectNums = append(selectNums, idx)
			}
		}
	}

	if len(selectNums) > 0 {
		fmt.Println("Subcription updated:", subItem.Name, "updated to EP.", maxEp)

		magnets := web.GetMagnets(selectNums)
		needTrimMagnet := viper.GetBool("enable-trim-magnet")
		if needTrimMagnet {
			trimMagnets(magnets)
		}

		useNameAsFolder := viper.GetBool("use-name-as-subscription-folder")
		useSeasonFolder := viper.GetBool("use-season-folder")
		dir := viper.GetString("default-download-dir")

		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Println("Cannot create Download folder:", err)
				return
			}
		}

		if useNameAsFolder && subItem.Name != "" {
			dirTest := path.Join(dir, subItem.Name)
			if useSeasonFolder && subItem.Season != "" {
				dirTest = path.Join(dirTest, subItem.Season)
			}

			if _, err := os.Stat(dirTest); os.IsNotExist(err) {
				if err := os.MkdirAll(dirTest, 0755); err == nil {
					dir = dirTest
				} else {
					fmt.Println("Cannot create Bangumi folder:", err)
				}
			} else {
				dir = dirTest
			}
		}

		err := downloadMagnets(magnets, dir)
		if err != nil {
			fmt.Println("Cannot connect to aria2:", err)
			return
		}

		subItem.Progress = maxEp
		subItem.Time = time.Now().Unix()
		updatedSubItems <- subItem
	}
}

func update() {
	var wg sync.WaitGroup
	wg.Add(len(globalData.Sublist))
	updatedSubItems := make(chan *SubItem, len(globalData.Sublist))

	for _, subItem := range globalData.Sublist {
		go doUpdateSubItem(subItem, &wg, updatedSubItems)
	}

	fmt.Println("Checking update...")
	wg.Wait()
	close(updatedSubItems)

	fmt.Println("Update completed:", len(updatedSubItems), "subscription(s) updated.")
	if len(updatedSubItems) > 0 {
		notifyUpdate(updatedSubItems)
	}

	err := writeData()
	if err != nil {
		fmt.Println("Cannot update subscription progress", err)
		return
	}
}
