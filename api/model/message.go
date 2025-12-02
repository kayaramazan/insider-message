package model

import (
	"time"

	"github.com/go-playground/validator/v10"
)

// MessageStatus represents the status of a message
type MessageStatus int

const (
	MessageStatusPending MessageStatus = 1
	MessageStatusSent    MessageStatus = 2
)

// String returns the string representation of the status
func (s MessageStatus) String() string {
	switch s {
	case MessageStatusPending:
		return "pending"
	case MessageStatusSent:
		return "sent"
	default:
		return "unknown"
	}
}

type Message struct {
	ID        string        `json:"id" validate:"uuid4"`
	Content   string        `json:"content" validate:"required,max=200"`
	Phone     string        `json:"phone" validate:"required"`
	Status    MessageStatus `json:"status"`
	CreatedAt time.Time     `json:"created_at"`
}

// Message için bir Validate fonksiyonu yazabilirsiniz:

func (m *Message) Validate() error {
	validate := validator.New()
	return validate.Struct(m)
}
