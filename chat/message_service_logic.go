package chat

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	WrapperMessageServiceProcessMethod              = "MessageService.Process"
	WrapperMessageServiceGetMessageCollectionMethod = "MessageService.GetMessageCollection"
)

type messageService struct {
	messageRepository MessageRepositoryInterface
	logger            *logrus.Logger
}

func NewMessageService(messageRepository MessageRepositoryInterface, logger *logrus.Logger) MessageServiceInterface {
	return &messageService{
		messageRepository: messageRepository,
		logger:            logger,
	}
}

func (h *messageService) Process(input []byte) []byte {
	var message Message

	if err := json.Unmarshal(input, &message); err != nil {
		h.logger.Error(errors.Wrap(err, WrapperMessageServiceProcessMethod))
		return []byte{}
	}

	message.CreatedAt = time.Now().UTC().Unix()
	if err := h.messageRepository.AddMessage(&message); err != nil {
		h.logger.Error(errors.Wrap(err, WrapperMessageServiceProcessMethod))
		return []byte{}
	}

	result, err := json.Marshal(message)
	if err != nil {
		h.logger.Error(errors.Wrap(err, WrapperMessageServiceProcessMethod))
		return []byte{}
	}

	return result
}

func (h *messageService) GetMessageCollection() [][]byte {
	var result [][]byte

	messages, err := h.messageRepository.GetMessageCollection()
	if err != nil {
		h.logger.Warn(errors.Wrap(err, WrapperMessageServiceGetMessageCollectionMethod))
	}

	for _, m := range messages {
		mtba, err := json.Marshal(m)
		if err != nil {
			h.logger.Error(errors.Wrap(err, WrapperMessageServiceGetMessageCollectionMethod))
			break
		}
		result = append(result, mtba)
	}
	return result
}
