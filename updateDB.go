package main

import (
	"encoding/binary"
	"strconv"

	log "github.com/google/logger"
	"github.com/syndtr/goleveldb/leveldb"
	ldb_opt "github.com/syndtr/goleveldb/leveldb/opt"
)

const ReadOnly = true

var versionKey = []byte{0x00, 'V', 'E', 'R', 'S', 'I', 'O', 'N'}

func (i *DBInfo) getVerion() (int64, error) {
	n, err := strconv.ParseInt(i.Version, 0, 64)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func readBmrkdDBVers(dbpath string) (int, error) {

	opt := &ldb_opt.Options{
		ErrorIfExist:   false,
		ErrorIfMissing: true,
		ReadOnly:       true,
	}

	db, err := leveldb.OpenFile(dbpath, opt)
	if nil != err {
		return 0, err
	}

	versionValue, err := db.Get(versionKey, nil)
	if leveldb.ErrNotFound == err {
		return 0, nil
	} else if nil != err {
		db.Close()
		return 0, err
	}

	if 4 != len(versionValue) {
		db.Close()
		log.Errorf("incompatible database version length: expected: %d  actual: %d", 4, len(versionValue))
		return 0, ErrorIncompatibleVersionLength
	}
	version := int(binary.BigEndian.Uint32(versionValue))
	return version, nil
}

func downloadDB(url string) error {
	return nil
}
