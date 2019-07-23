package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

type Data struct {
	Sublist  map[int]SubItem
	SubMaxNo int
}

var workDir string
var globalData *Data

var rootCmd = &cobra.Command{
	Use:   "bgmgo",
	Short: "An anime subscription tool.",
	Long:  ``,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	cobra.OnInitialize(initData)

	rootCmd.PersistentFlags().StringVar(&workDir, "workDir", "", "working directory (default is $HOME/.bgmgo)")
}

func initConfig() {
	if workDir == "" {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		workDir = path.Join(home, ".bgmgo")
		if _, err := os.Stat(workDir); os.IsNotExist(err) {
			os.Mkdir(workDir, 0755)
		}
	}

	viper.AutomaticEnv()

	initDefaultConfig()

	configFilePath := path.Join(workDir, "config.json")
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		if err := viper.WriteConfigAs(configFilePath); err != nil {
			fmt.Println("Cannot create config file:", err)
			return
		}
	}

	viper.AddConfigPath(workDir)
	viper.SetConfigFile(configFilePath)

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Cannot read config file:", err)
		return
	}
}

func initData() {
	dataFilePath := path.Join(workDir, "data.json")
	globalData = &Data{map[int]SubItem{}, 0}
	if _, err := os.Stat(dataFilePath); os.IsNotExist(err) {
		err := writeData()
		if err != nil {
			return
		}
	}

	dataJson, err := ioutil.ReadFile(dataFilePath)
	if err != nil {
		fmt.Println("Cannot read data file:", err)
		return
	}

	err = json.Unmarshal(dataJson, globalData)
	if err != nil {
		fmt.Println("Data format invalid:", err)
		return
	}
}

func writeData() error {
	dataFilePath := path.Join(workDir, "data.json")

	dataJson, err := json.Marshal(globalData)
	if err != nil {
		fmt.Println("Cannot write data:", err)
		return err
	}

	err = ioutil.WriteFile(dataFilePath, dataJson, 0644)
	if err != nil {
		fmt.Println("Cannot write data:", err)
		return err
	}
	return nil
}

func initDefaultConfig() {
	viper.SetDefault("data-source", "dmhy")
	viper.SetDefault("aria2-rpc-address", "http://localhost:6800/jsonrpc")
	viper.SetDefault("aria2-rpc-token", "")
	viper.SetDefault("enable-trim-magnet", true)

	home, _ := homedir.Dir()

	viper.SetDefault("default-download-dir", home)
	viper.SetDefault("use-name-as-subscription-folder", false)
	viper.SetDefault("use-pattern-as-subscription-folder", false)
}
