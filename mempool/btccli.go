package mempool

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"strings"
)

type BitcoinCli struct {
	Dir string
	Bin string
}

type Txfees struct {
	Base       float64 `json:"base"`
	Modified   float64 `json:"modified"`
	Ancestor   float64 `json:"ancestor"`
	Descendant float64 `json:"descendant"`
}

type Txentry struct {
	Fee Txfees `json:"fees"`
	// Vsize             int64    `json:"vsize"`
	// Weight            int64    `json:"weight"`
	Time int64 `json:"time"`
	// Height            int64    `json:"height"`
	// Descendantcount   int64    `json:"descendantcount"`
	// Ancestorcount     int64    `json:"ancestorcount"`
	// Ancestorsize      int64    `json:"ancestorsize"`
	Wtxid string `json:"wtxid"`
	// Bip125Replaceable bool     `json:"bip125-replaceable"`
	// Unbroadcast       bool     `json:"unbroadcast"`
	Depends []string `json:"depends"`
	Spentby []string `json:"Spentby"`
}

func rawmempool(dir, bin string, verbose, debug bool) (out []byte, err error) {
	binpath := filepath.Join(dir, bin)
	var fields []string
	if verbose {
		fields = []string{binpath, "getrawmempool", "true"}
	} else {
		fields = []string{binpath, "getrawmempool", "false"}
	}
	cmd := exec.Command(fields[0], fields[1:]...)
	cmd.Dir = dir
	if debug {
		log.Printf("%s\n", strings.Join(fields, " "))
	}

	out, err = cmd.Output()
	if err != nil {
		if debug {
			log.Printf("%V\n", err)
		}
		return
	}
	return
}

func (b BitcoinCli) Getrawmempool(debug bool) (resp []string, err error) {
	out, err := rawmempool(b.Dir, b.Bin, false, debug)
	if err != nil {
		return
	}
	if err = json.Unmarshal(out, &resp); err != nil {
		return
	}
	return
}

func (b BitcoinCli) GetrawmempoolVerbose(debug bool) (resp map[string]Txentry, err error) {
	// fmt.Printf("%+v\n", b)
	out, err := rawmempool(b.Dir, b.Bin, true, debug)
	if err != nil {
		return
	}
	if err = json.Unmarshal(out, &resp); err != nil {
		return
	}
	return
}

func (b BitcoinCli) Getmempoolentry(txid string, debug bool) (resp Txentry, notin bool, err error) {
	notin = false
	binpath := filepath.Join(b.Dir, b.Bin)
	fields := []string{binpath, "getmempoolentry", txid}
	cmd := exec.Command(fields[0], fields[1:]...)
	cmd.Dir = b.Dir
	if debug {
		log.Printf("%s\n", strings.Join(fields, " "))
	}
	var out []byte
	out, err = cmd.Output()
	if err != nil {
		if debug {
			log.Printf("%v\n", err)
		}
		if err.Error() == "exit status 5" {
			err = fmt.Errorf("transaction not in mempool: %v", err)
			notin = true
			return
		}
		//log.Printf("%s\n", strings.Join(fields, " "))
		err = fmt.Errorf("Getmempoolentry err: %v", err)
		return
	}
	if err = json.Unmarshal(out, &resp); err != nil {
		if debug {
			log.Printf("out: %s\n", string(out))
		}
		err = fmt.Errorf("unmarshal err: %v", err)
		return
	}
	return
}
