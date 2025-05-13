package usecases

import (
	"balancer/internal/adapters/db"
	"balancer/internal/domain/models"
	"balancer/internal/logger"
	"context"
	"errors"
	"sync"
	"time"
)

type ClientUseCase struct {
	repo  db.ClientRepository
	cache sync.Map // map[string]*models.Client
}

func NewClientUseCase(ctx context.Context, repo db.ClientRepository) (*ClientUseCase, error) {
	uc := &ClientUseCase{repo: repo}

	clients, err := repo.ListClients(ctx)
	if err != nil {
		return nil, err
	}

	for _, c := range clients {
		uc.cache.Store(c.ID, c)
	}

	logger.Info(ctx, "loaded clients into cache", "count", len(clients))
	return uc, nil
}

func (uc *ClientUseCase) Create(ctx context.Context, client *models.Client) (*models.Client, error) {
	if client.ID == "" || client.Capacity <= 0 || client.RatePerSec <= 0 {
		return nil, errors.New("invalid client parameters")
	}

	now := time.Now().UTC()
	client.CreatedAt = now
	client.UpdatedAt = now

	created, err := uc.repo.CreateClient(ctx, client)
	if err != nil {
		return nil, err
	}

	uc.cache.Store(created.ID, created)
	return created, nil
}

func (uc *ClientUseCase) Get(ctx context.Context, id string) (*models.Client, error) {
	value, ok := uc.cache.Load(id)
	if !ok {
		return nil, errors.New("client not found")
	}
	return value.(*models.Client), nil
}

func (uc *ClientUseCase) Update(ctx context.Context, client *models.Client) error {
	if client.ID == "" || client.Capacity <= 0 || client.RatePerSec <= 0 {
		return errors.New("invalid client parameters")
	}

	client.UpdatedAt = time.Now().UTC()

	if err := uc.repo.UpdateClient(ctx, client); err != nil {
		return err
	}

	uc.cache.Store(client.ID, client)
	return nil
}

func (uc *ClientUseCase) Delete(ctx context.Context, id string) error {
	if err := uc.repo.DeleteClient(ctx, id); err != nil {
		return err
	}

	uc.cache.Delete(id)
	return nil
}

func (uc *ClientUseCase) List(ctx context.Context) ([]*models.Client, error) {
	var result []*models.Client
	uc.cache.Range(func(_, value any) bool {
		result = append(result, value.(*models.Client))
		return true
	})
	return result, nil
}
