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
	"log"
	"net"

	"github.com/hypoballad/mempoolhist/mempool"
	"github.com/hypoballad/mempoolhist/mempoolapi"
	"github.com/robfig/cron"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/syndtr/goleveldb/leveldb"
	"google.golang.org/grpc"
)

type server struct {
	mempoolapi.UnimplementedMempoolServiceServer
}

func (s server) GetMementry(ctx context.Context, in *mempoolapi.TxidParam) (*mempoolapi.JsonResp, error) {
	var item mempoolapi.JsonResp
	resp, err := mempool.GetMementry(db, in.Txid)
	if err != nil {
		return &item, err
	}
	txentry, err := json.Marshal(resp)
	if err != nil {
		return &item, err
	}
	item.Json = string(txentry)
	return &item, nil
}

func (s server) GetMementryTime(ctx context.Context, in *mempoolapi.TxidParam) (*mempoolapi.TimeResp, error) {
	var item mempoolapi.TimeResp
	entryTime, err := mempool.GetMementryTime(db, in.Txid)
	if err != nil {
		return &item, err
	}
	item.Uts = entryTime
	return &item, nil
}

func (s server) FindMempoolhist(ctx context.Context, in *mempoolapi.TimerangeParam) (*mempoolapi.MemHistArray, error) {
	var item mempoolapi.MemHistArray
	hist, err := mempool.FindMempoolhist(db, in.Start, in.Stop, in.Asc)
	if err != nil {
		// log.Println(err)
		return &item, err
	}
	// fmt.Printf("len(hist): %d\n", len(hist))
	memhist := []*mempoolapi.MemHist{}
	for _, mh := range hist {
		h := mempoolapi.MemHist{
			Uts:  mh.Uts,
			Txid: mh.Txid,
		}
		memhist = append(memhist, &h)
	}
	item.Memhist = memhist
	return &item, nil
}

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "A server that downloads mempool data according to a schedule and sends it using gRPC",
	Long: `A server that downloads mempool data according to a schedule and sends it using gRPC. For example:

./mempoolhist server`,
	Run: func(cmd *cobra.Command, args []string) {
		db, err = leveldb.OpenFile(viper.GetString("root.db"), nil)
		if err != nil {
			log.Fatalln(err)
		}
		defer db.Close()
		c := cron.New()
		cli := mempool.BitcoinCli{
			Dir: viper.GetString("bitcoin.dir"),
			Bin: viper.GetString("bitcoin.bin"),
		}
		job := mempool.Job{
			Cli:   cli,
			DB:    db,
			Debug: viper.GetBool("root.debug"),
		}
		if err := job.DownloadMementry(true); err != nil {
			log.Printf("download mementry err: %v", err)
		}
		log.Printf("add schedule: %s\n", viper.GetString("job.download_mementry"))
		c.AddFunc(viper.GetString("job.download_mementry"), func() {
			// fmt.Println("DownloadRowmempool")
			if err := job.DownloadMementry(false); err != nil {
				log.Printf("download mementry err: %v", err)
			}
		})
		c.Start()
		addr := viper.GetString("root.addr")
		var lis net.Listener
		lis, err = net.Listen("tcp", addr)
		if err != nil {
			log.Fatalln("failed to listen: %v", err)
		}
		s := grpc.NewServer()
		log.Printf("listen to %s\n", addr)
		mempoolapi.RegisterMempoolServiceServer(s, &server{})
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
