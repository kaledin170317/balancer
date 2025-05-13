package db

import (
	"balancer/internal/domain/models"
	"context"
)

type ClientRepository interface {
	CreateClient(ctx context.Context, client *models.Client) (*models.Client, error)
	GetClient(ctx context.Context, id string) (*models.Client, error)
	UpdateClient(ctx context.Context, client *models.Client) error
	DeleteClient(ctx context.Context, id string) error
	ListClients(ctx context.Context) ([]*models.Client, error)
}
