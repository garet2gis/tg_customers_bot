package chat_repository

import (
	"errors"
	"github.com/garet2gis/tg_customers_bot/internal/chat_repository"
	"github.com/garet2gis/tg_customers_bot/pkg/client/boltdb"
	"github.com/garet2gis/tg_customers_bot/pkg/logging"
	"github.com/garet2gis/tg_customers_bot/pkg/utils"
	"go.etcd.io/bbolt"
)

var NoChatStateFound = errors.New("no chat state found")

type ChatRepository struct {
	db     boltdb.KeyValueClient
	logger *logging.Logger
}

func (c *ChatRepository) Delete(userID int64, bucket string) error {
	err := c.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		return b.Delete(utils.Itob(userID))
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *ChatRepository) Update(userId int64, state string, bucket string) error {
	err := c.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		return b.Put(utils.Itob(userId), []byte(state))
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *ChatRepository) Get(userId int64, bucket string) (*chat_repository.State, error) {
	var state string
	err := c.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		state = string(b.Get(utils.Itob(userId)))
		return nil
	})
	if err != nil {
		return nil, err
	}

	if state == "" {
		return nil, NoChatStateFound
	}
	stateStruct := chat_repository.NewState(state)

	return &stateStruct, nil
}

func NewChatRepository(db boltdb.KeyValueClient, logger *logging.Logger) chat_repository.ChatRepository {
	return &ChatRepository{
		db:     db,
		logger: logger,
	}
}
