package service

import (
	"context"
	"messanger/internal/entity"
)

type MessageBDService interface {
	CreateMessage(ctx context.Context, query *entity.Message) error
	DeleteMessage(ctx context.Context, query entity.Message) error
	UpdateMessage(ctx context.Context, query entity.Message) error
	SetStatusToTrue(ctx context.Context, query uint) error
	GetUnsentMessages(ctx context.Context, query entity.Messanger) ([]entity.Message, error)
	GetMessagesById(ctx context.Context, id uint) ([]entity.Message, error)
	GetMessages(ctx context.Context) ([]entity.Message, error)
}

type MessageService struct {
	DB MessageBDService
}

func NewMessageService(bd MessageBDService) *MessageService {
	return &MessageService{
		DB: bd,
	}
}

func (ms *MessageService) CreateMessage(ctx context.Context, query *entity.Message) error {
	return ms.DB.CreateMessage(ctx, query)
}
func (ms *MessageService) DeleteMessage(ctx context.Context, query entity.Message) error {
	return ms.DB.DeleteMessage(ctx, query)
}
func (ms *MessageService) UpdateMessage(ctx context.Context, query entity.Message) error {
	return ms.DB.UpdateMessage(ctx, query)
}
func (ms *MessageService) SetStatusToTrue(ctx context.Context, query uint) error {
	return ms.DB.SetStatusToTrue(ctx, query)
}
func (ms *MessageService) GetUnsentMessages(ctx context.Context, query entity.Messanger) ([]entity.Message, error) {
	return ms.DB.GetUnsentMessages(ctx, query)
}
func (ms *MessageService) GetMessagesById(ctx context.Context, id uint) ([]entity.Message, error) {
	return ms.DB.GetMessagesById(ctx, id)
}
func (ms *MessageService) GetMessages(ctx context.Context) ([]entity.Message, error) {
	return ms.DB.GetMessages(ctx)
}
