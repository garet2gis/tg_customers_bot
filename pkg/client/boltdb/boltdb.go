package boltdb

import (
	"github.com/garet2gis/tg_customers_bot/pkg/logging"
	"go.etcd.io/bbolt"
)

type Bucket string

type KeyValueClient interface {
	Close() error
	View(fn func(*bbolt.Tx) error) error
	Batch(fn func(*bbolt.Tx) error) error
	Update(fn func(*bbolt.Tx) error) error
}

func NewKeyValueClient(dbName string, buckets []string, logger *logging.Logger) (KeyValueClient, error) {
	db, err := bbolt.Open(dbName, 0600, nil)
	if err != nil {
		return nil, err
	}

	if err = db.Batch(func(tx *bbolt.Tx) error {
		for i := 0; i < len(buckets); i++ {
			_, err = tx.CreateBucketIfNotExists([]byte(buckets[i]))
			if err != nil {
				return err
			}

		}

		return err
	}); err != nil {
		return nil, err
	}

	return db, nil
}
