package db

import (
	"bytes"
	"encoding/gob"
	"github.com/garet2gis/tg_customers_bot/internal/paid_service"
	"github.com/garet2gis/tg_customers_bot/pkg/client/boltdb"
	"github.com/garet2gis/tg_customers_bot/pkg/logging"
	"github.com/garet2gis/tg_customers_bot/pkg/utils"
	"go.etcd.io/bbolt"
	"log"
)

type ServiceTemporaryRepository struct {
	db     boltdb.KeyValueClient
	logger *logging.Logger
}

func (c *ServiceTemporaryRepository) Delete(userID int64, bucket string) error {
	err := c.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		return b.Delete(utils.Itob(userID))
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *ServiceTemporaryRepository) Update(userId int64, service *paid_service.CreatePaidServiceDTO, bucket string) error {
	err := c.db.Update(func(tx *bbolt.Tx) error {

		var byteService bytes.Buffer
		enc := gob.NewEncoder(&byteService)
		err := enc.Encode(*service)
		if err != nil {
			log.Fatal("encode error:", err)
		}

		b := tx.Bucket([]byte(bucket))
		return b.Put(utils.Itob(userId), byteService.Bytes())
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *ServiceTemporaryRepository) Get(userId int64, bucket string) (*paid_service.CreatePaidServiceDTO, error) {
	service := &paid_service.CreatePaidServiceDTO{}
	err := c.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		stateByte := b.Get(utils.Itob(userId))

		dec := gob.NewDecoder(bytes.NewBuffer(stateByte))
		err := dec.Decode(service)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return service, nil
}

func NewServiceTemporaryRepository(db boltdb.KeyValueClient, logger *logging.Logger) paid_service.ServiceTemporaryRepository {
	return &ServiceTemporaryRepository{
		db:     db,
		logger: logger,
	}
}
