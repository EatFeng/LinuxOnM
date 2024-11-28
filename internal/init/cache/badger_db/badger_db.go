package badger_db

import (
	"fmt"
	"github.com/dgraph-io/badger/v4"
	"time"
)

type Cache struct {
	db *badger.DB
}

func (c *Cache) Get(key string) ([]byte, error) {
	var result []byte
	err := c.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			result = append([]byte{}, val...)
			return nil
		})
		return err
	})
	return result, err
}

func (c *Cache) SetWithTTL(key string, value interface{}, duration time.Duration) error {
	err := c.db.Update(func(txn *badger.Txn) error {
		v := []byte(fmt.Sprintf("%v", value))
		e := badger.NewEntry([]byte(key), v).WithTTL(duration)
		return txn.SetEntry(e)
	})
	return err
}
