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

	"github.com/hypoballad/mempoolhist/v2/mempoolapi"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

// timeCmd represents the time command
var timeCmd = &cobra.Command{
	Use:   "time",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("time called")
		if len(args) != 1 {
			log.Fatalln("arg txid is required.")
		}
		txid := args[0]
		addr := viper.GetString("root.addr")
		conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalf("did not connect: %v\n", err)
		}
		defer conn.Close()
		c := mempoolapi.NewMempoolServiceClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		in := &mempoolapi.TxidParam{Txid: txid}
		var item *mempoolapi.TimeResp
		item, err = c.GetMementryTime(ctx, in)
		if err != nil {
			log.Fatalln(err)
		}
		if viper.GetBool("time.human") {
			tm := time.Unix(item.Uts, 0)
			fmt.Printf("%s\n", tm.Format("2006-01-02 15:04:05"))
		} else {
			fmt.Printf("%d\n", item.Uts)
		}
	},
}

func init() {
	rootCmd.AddCommand(timeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// timeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// timeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	timeCmd.Flags().BoolP("human", "H", false, "print sizes timestamp human readable format")
	viper.BindPFlag("time.human", timeCmd.Flags().Lookup("human"))
}
