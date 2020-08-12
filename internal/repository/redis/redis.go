package redis

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/aintsashqa/go-chat/chat"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

const (
	// Errors
	ErrNewRedisClientWrapper       = "NewRedisClient"
	ErrNewMessageRepositoryWrapper = "NewMessageRepository"
	ErrAddMessageWrapper           = "MessageRepository.AddMessage"
	ErrGetMessageCollectionWrapper = "MessageRepository.GetMessageCollection"

	// Database fields
	FieldID        = "id"
	FieldAuthor    = "author"
	FieldContent   = "content"
	FieldCreatedAt = "created_at"
)

type messageRepository struct {
	client *redis.Client
	logger *logrus.Logger
}

func newRedisClient(redisURL string) (*redis.Client, error) {
	options, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, errors.Wrap(err, ErrNewRedisClientWrapper)
	}

	client := redis.NewClient(options)
	if _, err = client.Ping().Result(); err != nil {
		return nil, errors.Wrap(err, ErrNewRedisClientWrapper)
	}

	return client, nil
}

func NewMessageRepository(redisURL string, logger *logrus.Logger) (chat.MessageRepositoryInterface, error) {
	messageRepository := &messageRepository{
		logger: logger,
	}

	client, err := newRedisClient(redisURL)
	if err != nil {
		return nil, errors.Wrap(err, ErrNewMessageRepositoryWrapper)
	}

	messageRepository.client = client
	return messageRepository, nil
}

func (r *messageRepository) generateKey(message *chat.Message) string {
	id := uuid.NewV4()
	key := fmt.Sprintf("message:%s", id)
	message.ID = id.String()
	r.logger.Debugf("MessageRepository generated a database key [%v] for message", key)
	return key
}

func (r *messageRepository) AddMessage(message *chat.Message) error {
	r.logger.Debugf("MessageRepository trying to add new message [%v] to database", message)
	key := r.generateKey(message)
	data := map[string]interface{}{
		FieldID:        message.ID,
		FieldAuthor:    message.Author,
		FieldContent:   message.Content,
		FieldCreatedAt: message.CreatedAt,
	}
	_, err := r.client.HMSet(key, data).Result()
	if err != nil {
		return errors.Wrap(err, ErrAddMessageWrapper)
	}

	return nil
}

func (r *messageRepository) GetMessageCollection() ([]chat.Message, error) {
	r.logger.Debug("MessageRepository trying to get all messages from database")
	var messages []chat.Message
	r.logger.Debug("MessageRepository fetching all keys from database")
	keys, err := r.client.Keys("message:*").Result()
	r.logger.Debugf("MessageRepository keys count: %d", len(keys))
	if err != nil {
		return messages, errors.Wrap(err, ErrGetMessageCollectionWrapper)
	}

	for _, k := range keys {
		var m chat.Message
		data, err := r.client.HGetAll(k).Result()
		if err != nil {
			r.logger.Warn(errors.Wrap(err, ErrGetMessageCollectionWrapper))
			break
		}

		createdAt, err := strconv.ParseInt(data[FieldCreatedAt], 10, 64)
		if err != nil {
			r.logger.Warn(errors.Wrap(err, ErrGetMessageCollectionWrapper))
			break
		}

		m.ID = data[FieldID]
		m.Author = data[FieldAuthor]
		m.Content = data[FieldContent]
		m.CreatedAt = createdAt

		messages = append(messages, m)
	}

	sort.SliceStable(messages, func(i, j int) bool {
		return messages[i].CreatedAt < messages[j].CreatedAt
	})

	return messages, nil
}
