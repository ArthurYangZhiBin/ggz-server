package store

import (
	"flag"
	"fmt"
	"github.com/dgraph-io/badger"
)

var dataStr, valueDir string

var db *badger.DB

func init() {
	flag.StringVar(&dataStr, "dataDir", "/tmp", "数据存储目录")
	flag.StringVar(&valueDir, "valueDir", "/tmp", "值LOG存储目录")
	flag.Parse()

	opts := badger.DefaultOptions
	opts.Dir = dataStr
	opts.ValueDir = valueDir

	var err error
	db, err = badger.Open(opts)
	if err != nil {
		panic(err)
	}
}

func Close() {
	if db != nil {
		db.Close()
	}
}

func View(key string) ([]byte, error) {
	var value []byte
	err := db.View(func(txn *badger.Txn) error {
		it, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		// value, err = it.ValueCopy(nil)
		err = it.Value(func(val []byte) error {
			value = append([]byte{}, val...)
			return nil
		})
		return err
	})
	return value, err
}

func Store(key string, value []byte) error {
	return db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), value)
	})
}

func Remove(key string) error {
	return db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	})
}
func Blurry(prefix string) ([]string, error) {
	var value []string
	err := db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := []byte(prefix)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			k := item.Key()
			err := item.Value(func(val []byte) error {
				fmt.Printf("key=%s, value=%s\n", k, val)
				value = append(value, string(val))
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	return value, err
}
