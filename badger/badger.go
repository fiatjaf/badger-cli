package badger

import (
	"fmt"
	"time"

	"github.com/dgraph-io/badger/v4"
)

type DB struct {
	*badger.DB
}

type EntryOptions struct {
	TTL time.Duration
}

type ListResult struct {
	Key     string
	Size    int64
	Version uint64
	Meta    byte
}

func (l ListResult) String() string {
	return fmt.Sprintf("% -30s % 10d % 10d % 5s", l.Key, l.Size, l.Version, string(l.Meta))
}

func Open(dir string) (*DB, error) {
	opts := badger.DefaultOptions(dir)
	opts = opts.WithLogger(NewLogger())
	db, err := badger.Open(opts)
	return &DB{DB: db}, err
}

func (db *DB) Get(keys ...string) ([]string, error) {
	var values []string
	err := db.View(func(txn *badger.Txn) error {
		for _, k := range keys {
			item, err := txn.Get(stringToBytes(k))
			if err != nil {
				if err == badger.ErrKeyNotFound {
					return fmt.Errorf("Key %s not found", k)
				}
				return err
			}

			value, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}
			values = append(values, bytesToString(value))
		}

		return nil
	})
	return values, err
}

func (db *DB) List(prefix string, limit, offset int) ([]ListResult, int, error) {
	var keys []ListResult
	var total int

	err := db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false

		if limit > 0 {
			opts.PrefetchSize = limit
		}
		if prefix != "" {
			opts.Prefix = stringToBytes(prefix)
		}
		it := txn.NewIterator(opts)
		defer it.Close()

		currentOffset := 0
		for it.Rewind(); it.ValidForPrefix(stringToBytes(prefix)); it.Next() {
			total++
			currentOffset++
			if currentOffset < offset {
				continue
			}

			if len(keys) < limit {
				item := it.Item()
				keys = append(
					keys,
					ListResult{
						Key:     bytesToString(item.KeyCopy(nil)),
						Size:    item.EstimatedSize(),
						Version: item.Version(),
						Meta:    item.UserMeta(),
					},
				)
			}
		}

		return nil
	})

	return keys, total, err
}

func (db *DB) Set(key, value string, opts *EntryOptions) error {
	return db.Update(func(txn *badger.Txn) error {
		if opts == nil {
			return txn.Set(stringToBytes(key), stringToBytes(value))
		}

		e := badger.NewEntry(stringToBytes(key), stringToBytes(value))
		if opts.TTL > 0 {
			e.WithTTL(opts.TTL)
		}
		return txn.SetEntry(e)
	})
}

func (db *DB) Delete(keys ...string) error {
	return db.Update(func(txn *badger.Txn) error {
		for _, key := range keys {
			if err := txn.Delete(stringToBytes(key)); err != nil {
				return err
			}
		}

		return nil
	})
}
