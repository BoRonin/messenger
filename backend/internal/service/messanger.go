package service

import (
	"context"
	"messanger/internal/entity"
)

type MessangerBDService interface {
	CreateMessanger(ctx context.Context, query *entity.Messanger) error
	DeleteMessanger(ctx context.Context, query entity.Messanger) error
	UpdateMessanger(ctx context.Context, query entity.Messanger) error
	GetActiveMessangers(ctx context.Context) ([]entity.Messanger, error)
	GetMessangers(ctx context.Context) ([]entity.Messanger, error)
}
type MessangerStarter interface {
	SendMessages(ctx context.Context, messages []entity.Message) error
}
type MessangerService struct {
	DB MessangerBDService
	MS MessangerStarter
}

func NewMessangerService(db MessangerBDService, MS MessangerStarter) *MessangerService {
	return &MessangerService{
		DB: db,
		MS: MS,
	}
}

func (ms *MessangerService) CreateMessanger(ctx context.Context, query *entity.Messanger) error {
	return ms.DB.CreateMessanger(ctx, query)
}
func (ms *MessangerService) DeleteMessanger(ctx context.Context, query entity.Messanger) error {
	return ms.DB.DeleteMessanger(ctx, query)
}
func (ms *MessangerService) UpdateMessanger(ctx context.Context, query entity.Messanger) error {
	return ms.DB.UpdateMessanger(ctx, query)
}
func (ms *MessangerService) GetActiveMessangers(ctx context.Context) ([]entity.Messanger, error) {
	return ms.DB.GetActiveMessangers(ctx)
}
func (ms *MessangerService) GetMessangers(ctx context.Context) ([]entity.Messanger, error) {
	return ms.DB.GetMessangers(ctx)
}

func (ms *MessangerService) SendMessages(ctx context.Context, messages []entity.Message) error {
	return ms.MS.SendMessages(ctx, messages)
}
