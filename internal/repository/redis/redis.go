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
	WrapperNewRedisClientMethod                        = "NewRedisClient"
	WrapperNewMessageRepositoryMethod                  = "NewMessageRepository"
	WrapperMessageRepositoryAddMessageMethod           = "MessageRepository.AddMessage"
	WrapperMessageRepositoryGetMessageCollectionMethod = "MessageRepository.GetMessageCollection"
	WrapperMessageRepositoryGenerateKeyMethod          = "MessageRepository.GenerateKey"

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
		return nil, errors.Wrap(err, WrapperNewRedisClientMethod)
	}

	client := redis.NewClient(options)
	if _, err = client.Ping().Result(); err != nil {
		return nil, errors.Wrap(err, WrapperNewRedisClientMethod)
	}

	return client, nil
}

func NewMessageRepository(redisURL string, logger *logrus.Logger) (chat.MessageRepositoryInterface, error) {
	messageRepository := &messageRepository{
		logger: logger,
	}

	client, err := newRedisClient(redisURL)
	if err != nil {
		return nil, errors.Wrap(err, WrapperNewMessageRepositoryMethod)
	}

	messageRepository.client = client
	return messageRepository, nil
}

func (r *messageRepository) generateKey(hubId string, message *chat.Message) string {
	id := uuid.NewV4()
	key := fmt.Sprintf("hub-%s-message:%s", hubId, id)
	message.ID = id.String()
	r.logger.Debugf("[%s]: Generated a database key [%s] for message", WrapperMessageRepositoryGenerateKeyMethod, key)
	return key
}

func (r *messageRepository) AddMessage(hubId string, message *chat.Message) error {
	key := r.generateKey(hubId, message)
	r.logger.Debugf("[%s]: Trying to add new message [%+v] to database", WrapperMessageRepositoryAddMessageMethod, message)
	data := map[string]interface{}{
		FieldID:        message.ID,
		FieldAuthor:    message.Author,
		FieldContent:   message.Content,
		FieldCreatedAt: message.CreatedAt,
	}
	_, err := r.client.HMSet(key, data).Result()
	if err != nil {
		return errors.Wrap(err, WrapperMessageRepositoryAddMessageMethod)
	}

	return nil
}

func (r *messageRepository) GetMessageCollection(hubId string) ([]chat.Message, error) {
	var messages []chat.Message
	pattern := fmt.Sprintf("hub-%s-message:*", hubId)
	keys, err := r.client.Keys(pattern).Result()
	if err != nil {
		return messages, errors.Wrap(err, WrapperMessageRepositoryGetMessageCollectionMethod)
	}

	for _, k := range keys {
		var m chat.Message
		data, err := r.client.HGetAll(k).Result()
		if err != nil {
			r.logger.Warn(errors.Wrap(err, WrapperMessageRepositoryGetMessageCollectionMethod))
			break
		}

		createdAt, err := strconv.ParseInt(data[FieldCreatedAt], 10, 64)
		if err != nil {
			r.logger.Warn(errors.Wrap(err, WrapperMessageRepositoryGetMessageCollectionMethod))
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
