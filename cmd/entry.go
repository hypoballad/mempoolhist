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
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hypoballad/mempoolhist/mempool"
	"github.com/hypoballad/mempoolhist/mempoolapi"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

// entryCmd represents the entry command
var entryCmd = &cobra.Command{
	Use:   "entry",
	Short: "Get the time of the first occurrence by specifying txid.",
	Long: `Get the time of the first occurrence by specifying txid. For example:

	./mempoolhist time 0c42e5341cd59d944d268ee34fd5e2c90896ce64c23989ca4f9ad005f7b11c5e -H`,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("entry called")
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
		var item *mempoolapi.JsonResp
		item, err = c.GetMementry(ctx, in)
		if err != nil {
			log.Fatalln(err)
		}
		var txentry mempool.Txentry
		if err = json.Unmarshal([]byte(item.Json), &txentry); err != nil {
			log.Fatalln(err)
		}
		var b []byte
		if viper.GetBool("entry.indent") {
			b, err = json.MarshalIndent(txentry, "", "	")
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Println(string(b))
		} else {
			b, err = json.Marshal(txentry)
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Println(string(b))
		}

	},
}

func init() {
	rootCmd.AddCommand(entryCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// entryCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// entryCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	entryCmd.Flags().BoolP("indent", "I", false, "marshal indent")
	viper.BindPFlag("entry.indent", entryCmd.Flags().Lookup("indent"))
}
