package cmd

import (
	"fmt"
	"math"
	"os"
	"path"
	"strconv"
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

func notifyUpdate(updatedSubItems []*SubItem) {
	title := fmt.Sprintf("Bgmgo Updated %d Subcription(s)", len(updatedSubItems))
	msg := ""
	for _, subItem := range updatedSubItems {
		if subItem.Progress == math.Floor(subItem.Progress) {
			msg += fmt.Sprintf("%s updated to EP.%d.\n", subItem.Name, int(subItem.Progress))
		} else {
			msg += fmt.Sprintf("%s updated to EP.%.1f.\n", subItem.Name, subItem.Progress)
		}
	}
	beeep.Notify(title, msg, "")
}

func update() {
	dataSource := viper.GetString("data-source")
	if parser.ParserCtor[dataSource] == nil {
		fmt.Println("Not valid data source: ", dataSource)
		return
	}

	updatedSubItems := []*SubItem{}
	web := parser.ParserCtor[dataSource]()

	for _, subItem := range globalData.Sublist {
		maxEp := subItem.Progress
		web.Request([]string{subItem.Pattern})
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
			dir := viper.GetString("default-download-dir")

			if _, err := os.Stat(dir); os.IsNotExist(err) {
				if err := os.MkdirAll(dir, 0755); err != nil {
					fmt.Println("Cannot create Download folder:", err)
					return
				}
			}

			if useNameAsFolder && subItem.Name != "" {
				dirTest := path.Join(dir, subItem.Name)

				if _, err := os.Stat(dirTest); os.IsNotExist(err) {
					if err := os.Mkdir(dirTest, 0755); err == nil {
						dir = dirTest
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
			updatedSubItems = append(updatedSubItems, subItem)
		}
	}

	fmt.Println("Update completed:", len(updatedSubItems), "subscription(s) updated.")
	if len(updatedSubItems) > 0 {
		notifyUpdate(updatedSubItems)
	}
	writeData()
}
