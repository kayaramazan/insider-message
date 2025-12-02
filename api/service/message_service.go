package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kayaramazan/insider-message/api/cache"
	"github.com/kayaramazan/insider-message/api/model"
	"github.com/kayaramazan/insider-message/api/repository"
)

type MessageService struct {
	MessageRepo     *repository.MessageRepository
	RedisCache      *cache.RedisCache
	webhookUrl      string
	messagePerCycle int
}

func NewMessageService(messageRepo *repository.MessageRepository, cache *cache.RedisCache, url string, messagePerCycle int) *MessageService {
	return &MessageService{
		MessageRepo:     messageRepo,
		RedisCache:      cache,
		webhookUrl:      url,
		messagePerCycle: messagePerCycle,
	}
}

func (s *MessageService) GetAllSentMessages(ctx context.Context) ([]model.Message, error) {
	return s.MessageRepo.GetAllSentMessages(ctx)
}

func (s *MessageService) CreateMessage(ctx context.Context, message *model.Message) error {
	return s.MessageRepo.Create(ctx, message)
}

func (s *MessageService) SendMessage(ctx context.Context) error {
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
		s.MessageRepo.UpdateMessageStatus(ctx, message.ID, int(model.MessageStatusSent))

	}
	return nil
}

func (s *MessageService) sendWebhook(message *model.Message) error {
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
