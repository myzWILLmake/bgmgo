package cmd

import (
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type SubItem struct {
	No       int
	Name     string
	Time     int64
	Pattern  string
	Progress float64
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all your subscriptions.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		list()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func list() {
	if len(globalData.Sublist) == 0 {
		fmt.Println("No subscription in your list, use \"bgmgo sub\" to add some subscriptions.")
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	showPattern := viper.GetBool("show-pattern-in-subscription-list")
	if showPattern {
		table.SetHeader([]string{"Sub-Number", "Name", "Last-Update-Time", "Progress", "Pattern"})
	} else {
		table.SetHeader([]string{"Sub-Number", "Name", "Last-Update-Time", "Progress"})
	}

	keys := []int{}
	for key := range globalData.Sublist {
		keys = append(keys, key)
	}
	sort.Ints(keys)

	for _, key := range keys {
		subItem := globalData.Sublist[key]

		var tableItem []string
		if showPattern {
			tableItem = make([]string, 5)
			tableItem[4] = subItem.Pattern
		} else {
			tableItem = make([]string, 4)
		}
		tableItem[0] = strconv.Itoa(subItem.No)
		tableItem[1] = subItem.Name
		if subItem.Time == 0 {
			tableItem[2] = "-"
		} else {
			tableItem[2] = time.Unix(subItem.Time, 0).Format("2006-01-02 15:04")
		}
		if subItem.Progress == math.Floor(subItem.Progress) {
			tableItem[3] = fmt.Sprintf("%d", int(subItem.Progress))
		} else {
			tableItem[3] = fmt.Sprintf("%.1f", subItem.Progress)
		}
		table.Append(tableItem)
	}

	table.Render()
}
