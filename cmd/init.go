/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"log"

	"github.com/hypoballad/mempoolhist/mempool"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/syndtr/goleveldb/leveldb"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Run getrawmemory to download the initial data",
	Long: `Run getrawmemory to download the initial data. For example:

./mempoolhist init`,
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println("init called")
		db, err = leveldb.OpenFile(viper.GetString("root.db"), nil)
		if err != nil {
			log.Fatalln(err)
		}
		defer db.Close()
		cli := mempool.BitcoinCli{
			Dir: viper.GetString("bitcoin.dir"),
			Bin: viper.GetString("bitcoin.bin"),
		}
		job := mempool.Job{
			Cli:   cli,
			DB:    db,
			Debug: viper.GetBool("root.debug"),
		}
		if err = job.DownloadRawMempool(true); err != nil {
			log.Fatalln(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
