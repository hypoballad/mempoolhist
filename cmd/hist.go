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
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hypoballad/mempoolhist/mempoolapi"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

// histCmd represents the hist command
var histCmd = &cobra.Command{
	Use:   "hist",
	Short: "Get the time and txid",
	Long: `Get the current time to its duration by specifying time.Duration in the go language. For example:

./mempoolhist hist 5m -H`,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("hist called")
		if len(args) != 1 {
			log.Fatalln("arg duration is required.")
		}
		duration := args[0]
		addr := viper.GetString("root.addr")
		conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalf("did not connect: %v\n", err)
		}
		defer conn.Close()
		c := mempoolapi.NewMempoolServiceClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		end := time.Now().UTC()
		d, err := time.ParseDuration(duration)
		if err != nil {
			log.Fatalln(err)
		}
		start := end.Add(-1 * d)
		asc := viper.GetBool("hist.asc")
		in := &mempoolapi.TimerangeParam{
			Start: start.Unix(),
			Stop:  end.Unix(),
			Asc:   asc,
		}
		var item *mempoolapi.MemHistArray
		item, err = c.FindMempoolhist(ctx, in)
		if err != nil {
			log.Fatalln(err)
		}
		for _, memhist := range item.Memhist {
			ts := fmt.Sprintf("%d", memhist.Uts)
			if viper.GetBool("hist.human") {
				tm := time.Unix(memhist.Uts, 0)
				ts = tm.Format("2006-01-02 15:04:05")
			}
			fmt.Printf("%s %s\n", ts, memhist.Txid)
		}
	},
}

func init() {
	rootCmd.AddCommand(histCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// histCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// histCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	histCmd.Flags().Bool("asc", false, "sort mode")
	histCmd.Flags().BoolP("human", "H", false, "print sizes timestamp human readable format")
	viper.BindPFlag("hist.asc", histCmd.Flags().Lookup("asc"))
	viper.BindPFlag("hist.human", histCmd.Flags().Lookup("human"))
}
