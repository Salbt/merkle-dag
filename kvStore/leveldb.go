package kvstore

import (
	"log"

	"github.com/syndtr/goleveldb/leveldb"
)

type LevelDB struct {
	db *leveldb.DB
}

func NewLevelDB(path string) *LevelDB {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		log.Fatal("create new leveldb failed", err)
	}

	return &LevelDB{
		db: db,
	}
}

func (kv LevelDB) Close() error {
	return kv.db.Close()
}

func (kv LevelDB) Get(key []byte) ([]byte, error) {
	return kv.db.Get(key, nil)
}

func (kv LevelDB) Put(key, value []byte) error {
	return kv.db.Put(key, value, nil)
}

func (kv LevelDB) Has(key []byte) (bool, error) {
	return kv.db.Has(key, nil)
}

func (kv LevelDB) Delete(key []byte) error {
	return kv.db.Delete(key, nil)
}
