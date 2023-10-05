package service

import (
	"context"
	"messanger/internal/entity"
)

type ClientBDService interface {
	CreateClient(ctx context.Context, query *entity.Client) error
	DeleteClient(ctx context.Context, query entity.Client) error
	UpdateClient(ctx context.Context, query entity.Client) error
	GetClientsByTag(ctx context.Context, tag string) ([]entity.Client, error)
}

type ClientService struct {
	DB ClientBDService
}

func NewClientService(db ClientBDService) *ClientService {
	return &ClientService{
		DB: db,
	}
}

func (CS *ClientService) CreateClient(ctx context.Context, query *entity.Client) error {
	return CS.DB.CreateClient(ctx, query)
}
func (CS *ClientService) DeleteClient(ctx context.Context, query entity.Client) error {
	return CS.DB.DeleteClient(ctx, query)
}
func (CS *ClientService) UpdateClient(ctx context.Context, query entity.Client) error {
	return CS.DB.UpdateClient(ctx, query)
}

func (CS *ClientService) GetClientsByTag(ctx context.Context, tag string) ([]entity.Client, error) {
	return CS.DB.GetClientsByTag(ctx, tag)
}
