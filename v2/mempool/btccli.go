package mempool

import (
	"fmt"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
)

type MempoolResp struct {
	Fee     float64  `json:"fee"`
	Time    int64    `json:"time"`
	Height  int64    `json:"height"`
	Depends []string `json:"depends"`
}

type BitcoinConf struct {
	Rpcconnect string
	Rpcuser    string
	Rpcpasswd  string
	Rpcport    int
}

func (b BitcoinConf) connect() (client *rpcclient.Client, err error) {
	connCfg := &rpcclient.ConnConfig{
		Host:         fmt.Sprintf("%s:%d", b.Rpcconnect, b.Rpcport),
		User:         b.Rpcuser,
		Pass:         b.Rpcpasswd,
		HTTPPostMode: true,
		DisableTLS:   true,
	}
	client, err = rpcclient.New(connCfg, nil)
	if err != nil {
		return
	}
	return
}

func (b BitcoinConf) GetrawmempoolVerbose() (resp map[string]MempoolResp, err error) {
	var client *rpcclient.Client
	client, err = b.connect()
	if err != nil {
		return
	}
	defer client.Shutdown()
	var result map[string]btcjson.GetRawMempoolVerboseResult
	result, err = client.GetRawMempoolVerbose()
	if err != nil {
		return
	}
	// var result map[string]MepoolResp
	resp = make(map[string]MempoolResp)
	for k, v := range result {
		resp[k] = MempoolResp{
			Fee:     v.Fee,
			Time:    v.Time,
			Height:  v.Height,
			Depends: v.Depends,
		}
	}
	return
}

func (b BitcoinConf) Getrawmempool() (resp []string, err error) {
	var client *rpcclient.Client
	client, err = b.connect()
	if err != nil {
		return
	}
	defer client.Shutdown()
	var result []*chainhash.Hash
	result, err = client.GetRawMempool()
	if err != nil {
		return
	}
	resp = []string{}
	for _, hash := range result {
		resp = append(resp, hash.String())
	}
	return
}

func (b BitcoinConf) GetMempoolEntry(txid string) (resp MempoolResp, err error) {
	var client *rpcclient.Client
	client, err = b.connect()
	if err != nil {
		return
	}
	defer client.Shutdown()
	var result *btcjson.GetMempoolEntryResult
	result, err = client.GetMempoolEntry(txid)
	if err != nil {
		//fmt.Printf("err %+v\n", err)
		return
	}
	resp = MempoolResp{
		Fee:     result.Fee,
		Time:    result.Time,
		Height:  result.Height,
		Depends: result.Depends,
	}

	return
}
