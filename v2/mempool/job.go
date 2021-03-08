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
	Cli   BitcoinConf
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

func SaveMementry(db *leveldb.DB, txid string, mempoolResp MempoolResp) (err error) {
	key := fmt.Sprintf("mempool::%s", txid)
	b, err := json.Marshal(mempoolResp)
	if err != nil {
		return
	}
	err = db.Put([]byte(key), b, nil)
	return
}

func GetMementry(db *leveldb.DB, txid string) (mempoolResp MempoolResp, err error) {

	key := fmt.Sprintf("mempool::%s", txid)
	b, err := db.Get([]byte(key), nil)
	if err != nil {

		return
	}
	if err = json.Unmarshal(b, &mempoolResp); err != nil {

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
	var mempoolResp MempoolResp
	if err = json.Unmarshal(b, &mempoolResp); err != nil {
		return
	}
	entryTime = mempoolResp.Time
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
	var resp map[string]MempoolResp
	resp, err = j.Cli.GetrawmempoolVerbose()

	if err != nil {
		return
	}
	var bar *pb.ProgressBar
	if progress {
		bar = pb.StartNew(len(resp))
	}
	for txid, mempoolResp := range resp {
		if progress {
			bar.Increment()
		}
		if IsMissingTx(j.DB, txid) {
			continue
		}
		if IsMementry(j.DB, txid) {
			continue
		}
		if err = SaveMempoolhist(j.DB, txid, mempoolResp.Time); err != nil {
			return
		}
		if err = SaveMementry(j.DB, txid, mempoolResp); err != nil {
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
	var mempoolResp MempoolResp
	resp, err = j.Cli.Getrawmempool()
	if err != nil {
		return
	}
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
		mempoolResp, err = j.Cli.GetMempoolEntry(txid)

		if err != nil {
			if err = SaveMissingTx(j.DB, txid); err != nil {
				if j.Debug {
					log.Printf("SaveMissingTx : %s : %v\n", txid, err)
				}
				continue
			}
			if j.Debug {
				log.Printf("GetMempoolEntry : %s : %v\n", txid, err)
			}
			continue
		}
		if err = SaveMempoolhist(j.DB, txid, mempoolResp.Time); err != nil {
			if j.Debug {
				log.Printf("SaveMempoolhist : %s : %v\n", txid, err)
			}
			continue
		}
		if err = SaveMementry(j.DB, txid, mempoolResp); err != nil {
			if j.Debug {
				log.Println("SaveMempoolhist : %s : %v\n", txid, err)
			}
			continue
		}
	}
	if progress {
		bar.Finish()
	}
	return
}
