package cmd

import (
	"fmt"

	"github.com/myzWILLmake/bgmgo/parser"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var findCmd = &cobra.Command{
	Use:   "find",
	Short: "Find torrent descriptions from a torrent site by a provided pattern.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		find(args)
	},
}

var findSource string

func init() {
	rootCmd.AddCommand(findCmd)

	findCmd.Flags().StringVarP(&findSource, "source", "s", "", "set data source where to search magnet links")
}

func find(args []string) {
	if findSource == "" {
		findSource = viper.GetString("data-source")
	}

	if parser.ParserCtor[findSource] == nil {
		fmt.Println("Error: Not a valid data source", findSource)
		return
	}

	web := parser.ParserCtor[findSource]()
	err := web.Request(args)
	if err != nil {
		fmt.Println("Error:", err)
	}

	filterMap := map[string]int{
		"no":    0,
		"type":  1,
		"team":  2,
		"size":  3,
		"title": 4,
	}

	infos := web.ShowFindResult(filterMap, 5)

	n := len(infos)
	if n == 0 {
		fmt.Println("Not found any result, please check your pattern")
		return
	}
	fmt.Printf("Found %d record(s):\n", n)
	for _, info := range infos {
		for _, s := range info {
			fmt.Print(s, "  ")
		}
		fmt.Println()
	}
}
