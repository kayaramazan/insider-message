package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/kayaramazan/insider-message/api/cache"
	"github.com/kayaramazan/insider-message/api/model"
	"github.com/kayaramazan/insider-message/api/repository"
)

type MessageService interface {
	GetAllSentMessages(ctx context.Context) ([]model.Message, error)
	CreateMessage(ctx context.Context, message *model.Message) error
	SendMessage(ctx context.Context) error
}

type messageServiceImpl struct {
	MessageRepo     repository.MessageRepository
	RedisCache      cache.Cache
	webhookUrl      string
	messagePerCycle int
}

func NewMessageService(messageRepo repository.MessageRepository, cache cache.Cache, url string, messagePerCycle int) MessageService {
	return &messageServiceImpl{
		MessageRepo:     messageRepo,
		RedisCache:      cache,
		webhookUrl:      url,
		messagePerCycle: messagePerCycle,
	}
}

func (s *messageServiceImpl) GetAllSentMessages(ctx context.Context) ([]model.Message, error) {
	return s.MessageRepo.GetAllSentMessages(ctx)
}

func (s *messageServiceImpl) CreateMessage(ctx context.Context, message *model.Message) error {
	return s.MessageRepo.Create(ctx, message)
}

func (s *messageServiceImpl) SendMessage(ctx context.Context) error {
	messages, err := s.MessageRepo.GetUnsendMessages(ctx, s.messagePerCycle)
	if err != nil {
		return err
	}
	for _, message := range messages {
		err := s.sendWebhook(&message)
		if err != nil {
			continue
		}

		fmt.Printf("SENT: MessageID: %s, Content: %s \n", message.ID, message.Content)
		s.RedisCache.Set(ctx, message.ID, time.Now())
		s.MessageRepo.UpdateMessageStatus(ctx, message.ID, int(model.MessageStatusSent))

	}
	return nil
}

func (s *messageServiceImpl) sendWebhook(message *model.Message) error {
	json, err := json.Marshal(map[string]any{
		"to":      message.Phone,
		"content": message.Content,
	})
	if err != nil {
		return err
	}

	resp, err := http.Post(s.webhookUrl, "application/json", bytes.NewReader(json))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
