package repository

import (
	"context"

	"github.com/kayaramazan/insider-message/api/database"
	"github.com/kayaramazan/insider-message/api/model"
)

type MessageRepository interface {
	Create(ctx context.Context, message *model.Message) error
	GetAllSentMessages(ctx context.Context) ([]model.Message, error)
	GetUnsendMessages(ctx context.Context, limit int) ([]model.Message, error)
	UpdateMessageStatus(ctx context.Context, id string, status int) error
}

type messageRepositoryImpl struct {
	db database.Database
}

func NewMessageRepository(db database.Database) MessageRepository {
	return &messageRepositoryImpl{
		db: db,
	}
}

func (r *messageRepositoryImpl) Create(ctx context.Context, message *model.Message) error {
	query := `
        INSERT INTO messages (content, phone)
        VALUES ($1, $2)
        RETURNING id, created_at
    `
	_, err := r.db.Exec(ctx, query, message.Content, message.Phone)

	return err
}

func (r *messageRepositoryImpl) GetAllSentMessages(ctx context.Context) ([]model.Message, error) {
	var messages []model.Message
	query := `SELECT id, content, phone, created_at, status FROM messages where status = $1 ORDER BY created_at`

	row, err := r.db.Query(ctx, query, int(model.MessageStatusSent))
	if err != nil {
		return nil, err
	}
	defer row.Close()

	for row.Next() {
		var msg model.Message
		err := row.Scan(&msg.ID, &msg.Content, &msg.Phone, &msg.CreatedAt, &msg.Status)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	if row.Err() != nil {
		return nil, row.Err()
	}

	return messages, nil
}

func (r *messageRepositoryImpl) GetUnsendMessages(ctx context.Context, limit int) ([]model.Message, error) {
	var messages []model.Message
	query := `SELECT id, content, phone, created_at FROM messages where status != $1 ORDER BY created_at limit $2`

	row, err := r.db.Query(ctx, query, int(model.MessageStatusSent), limit)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	for row.Next() {
		var msg model.Message
		err := row.Scan(&msg.ID, &msg.Content, &msg.Phone, &msg.CreatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	if row.Err() != nil {
		return nil, row.Err()
	}

	return messages, nil
}

func (r *messageRepositoryImpl) UpdateMessageStatus(ctx context.Context, id string, status int) error {

	_, err := r.db.Exec(ctx, `UPDATE messages SET status = $1 WHERE id=$2`, status, id)

	return err
}
