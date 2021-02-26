package mempool

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/cheggaaa/pb/v3"
	"github.com/rs/xid"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type Job struct {
	Cli   BitcoinCli
	DB    *leveldb.DB
	Debug bool
}

type Mempoolhist struct {
	Uts  int64
	Txid string
}

func FindMempoolhist(db *leveldb.DB, start, stop int64, asc bool) (hist []Mempoolhist, err error) {
	hist = []Mempoolhist{}
	from := fmt.Sprintf("mempoolhist::%d", start)
	to := fmt.Sprintf("mempoolhist::%d", stop)
	iter := db.NewIterator(&util.Range{Start: []byte(from), Limit: []byte(to)}, nil)
	//iter := db.NewIterator(util.BytesPrefix([]byte("mempoolhist::")), nil)
	for iter.Next() {
		key := string(iter.Key())
		fields := strings.Split(key, "::")
		if len(fields) != 3 {
			err = fmt.Errorf("key error: %s", fields)
			return
		}
		ts := fields[1]
		var uts int64
		uts, err = strconv.ParseInt(ts, 10, 64)
		if err != nil {
			return
		}
		hist = append(hist, Mempoolhist{
			Uts:  uts,
			Txid: string(iter.Value()),
		})
		//fmt.Printf("%+v\n", hist)
	}
	iter.Release()
	if err = iter.Error(); err != nil {
		log.Println(err)
		return
	}

	if asc {
		return
	}

	sort.Slice(hist, func(i, j int) bool {
		return !(hist[i].Uts < hist[j].Uts)
	})

	return
}

func FindAllMempoolhist(db *leveldb.DB) (hist []Mempoolhist, err error) {
	hist = []Mempoolhist{}
	iter := db.NewIterator(util.BytesPrefix([]byte("mempoolhist::")), nil)
	for iter.Next() {
		key := string(iter.Key())
		fields := strings.Split(key, "::")
		if len(fields) != 2 {
			err = fmt.Errorf("key error: %s", fields)
			return
		}
		ts := fields[1]
		var uts int64
		uts, err = strconv.ParseInt(ts, 10, 64)
		if err != nil {
			return
		}
		hist = append(hist, Mempoolhist{
			Uts:  uts,
			Txid: string(iter.Value()),
		})
	}
	iter.Release()
	err = iter.Error()
	return
}

func SaveMempoolhist(db *leveldb.DB, txid string, uts int64) (err error) {
	guid := xid.New()
	key := fmt.Sprintf("mempoolhist::%d::%s", uts, guid.String())
	err = db.Put([]byte(key), []byte(txid), nil)
	return
}

func SaveMementry(db *leveldb.DB, txid string, txentry Txentry) (err error) {
	key := fmt.Sprintf("mempool::%s", txid)
	b, err := json.Marshal(txentry)
	if err != nil {
		return
	}
	err = db.Put([]byte(key), b, nil)
	return
}

func GetMementry(db *leveldb.DB, txid string) (txentry Txentry, err error) {
	key := fmt.Sprintf("mempool::%s", txid)
	b, err := db.Get([]byte(key), nil)
	if err != nil {
		return
	}
	if err = json.Unmarshal(b, &txentry); err != nil {
		return
	}
	return
}

func GetMementryTime(db *leveldb.DB, txid string) (entryTime int64, err error) {
	key := fmt.Sprintf("mempool::%s", txid)
	b, err := db.Get([]byte(key), nil)
	if err != nil {
		return
	}
	var txentry Txentry
	if err = json.Unmarshal(b, &txentry); err != nil {
		return
	}
	entryTime = txentry.Time
	return
}

func IsMementry(db *leveldb.DB, txid string) bool {
	key := fmt.Sprintf("mempool::%s", txid)
	_, err := db.Get([]byte(key), nil)
	if err != nil {
		return false
	}
	return true
}

func SaveMissingTx(db *leveldb.DB, txid string) (err error) {
	key := fmt.Sprintf("missingtx::%s", txid)
	err = db.Put([]byte(key), []byte(txid), nil)
	return
}

func IsMissingTx(db *leveldb.DB, txid string) bool {
	key := fmt.Sprintf("missingtx::%s", txid)
	_, err := db.Get([]byte(key), nil)
	if err != nil {
		return false
	}
	return true
}

func (j Job) DownloadRawMempool(progress bool) (err error) {
	var resp map[string]Txentry
	resp, err = j.Cli.GetrawmempoolVerbose(j.Debug)
	if err != nil {
		return
	}
	var bar *pb.ProgressBar
	if progress {
		bar = pb.StartNew(len(resp))
	}
	for txid, txentry := range resp {
		if progress {
			bar.Increment()
		}
		if IsMissingTx(j.DB, txid) {
			continue
		}
		if IsMementry(j.DB, txid) {
			continue
		}
		if err = SaveMempoolhist(j.DB, txid, txentry.Time); err != nil {
			return
		}
		if err = SaveMementry(j.DB, txid, txentry); err != nil {
			return
		}

	}
	if progress {
		bar.Finish()
	}
	return
}

func (j Job) DownloadMementry(progress bool) (err error) {
	var resp []string
	var notin bool
	var txentry Txentry
	resp, err = j.Cli.Getrawmempool(j.Debug)
	if err != nil {
		return
	}
	notin = false
	var bar *pb.ProgressBar
	if progress {
		bar = pb.StartNew(len(resp))
	}
	for _, txid := range resp {
		if progress {
			bar.Increment()
		}
		if IsMissingTx(j.DB, txid) {
			continue
		}
		if IsMementry(j.DB, txid) {
			continue
		}
		txentry, notin, err = j.Cli.Getmempoolentry(txid, j.Debug)
		if notin {
			if err = SaveMissingTx(j.DB, txid); err != nil {
				return
			}
		}
		if err != nil {
			return
		}
		if err = SaveMempoolhist(j.DB, txid, txentry.Time); err != nil {
			return
		}
		if err = SaveMementry(j.DB, txid, txentry); err != nil {
			return
		}
	}
	if progress {
		bar.Finish()
	}
	return
}
