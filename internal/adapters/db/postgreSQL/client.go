package postgreSQL

import (
	"balancer/internal/domain/models"
	"balancer/internal/logger"
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type ClientRepository struct {
	db          *sqlx.DB
	stmtBuilder sq.StatementBuilderType
}

func NewClientRepository(db *sqlx.DB) *ClientRepository {
	return &ClientRepository{
		db:          db,
		stmtBuilder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *ClientRepository) CreateClient(ctx context.Context, client *models.Client) (*models.Client, error) {
	logger.Info(ctx, "creating client in db", "id", client.ID)

	now := time.Now().UTC()
	client.CreatedAt = now
	client.UpdatedAt = now

	query, args, err := r.stmtBuilder.
		Insert("clients").
		Columns("id", "capacity", "rate_per_sec", "created_at", "updated_at").
		Values(client.ID, client.Capacity, client.RatePerSec, client.CreatedAt, client.UpdatedAt).
		ToSql()
	if err != nil {
		logger.Error(ctx, "sql build error on create", "err", err)
		return nil, err
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		logger.Error(ctx, "db insert failed", "err", err)
		return nil, err
	}

	logger.Info(ctx, "client created in db", "id", client.ID)
	return client, nil
}

func (r *ClientRepository) GetClient(ctx context.Context, id string) (*models.Client, error) {
	logger.Debug(ctx, "fetching client from db", "id", id)

	query, args, err := r.stmtBuilder.
		Select("id", "capacity", "rate_per_sec", "created_at", "updated_at").
		From("clients").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		logger.Error(ctx, "sql build error on get", "err", err)
		return nil, err
	}

	var client models.Client
	err = r.db.GetContext(ctx, &client, query, args...)
	if err != nil {
		logger.Error(ctx, "client fetch failed", "id", id, "err", err)
		return nil, err
	}

	logger.Debug(ctx, "client fetched", "id", id)
	return &client, nil
}

func (r *ClientRepository) UpdateClient(ctx context.Context, client *models.Client) error {
	logger.Info(ctx, "updating client in db", "id", client.ID)

	client.UpdatedAt = time.Now().UTC()

	query, args, err := r.stmtBuilder.
		Update("clients").
		Set("capacity", client.Capacity).
		Set("rate_per_sec", client.RatePerSec).
		Set("updated_at", client.UpdatedAt).
		Where(sq.Eq{"id": client.ID}).
		ToSql()
	if err != nil {
		logger.Error(ctx, "sql build error on update", "err", err)
		return err
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		logger.Error(ctx, "client update failed", "id", client.ID, "err", err)
		return err
	}

	logger.Info(ctx, "client updated", "id", client.ID)
	return nil
}

func (r *ClientRepository) DeleteClient(ctx context.Context, id string) error {
	logger.Warn(ctx, "deleting client", "id", id)

	query, args, err := r.stmtBuilder.
		Delete("clients").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		logger.Error(ctx, "sql build error on delete", "err", err)
		return err
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		logger.Error(ctx, "client delete failed", "id", id, "err", err)
		return err
	}

	logger.Warn(ctx, "client deleted", "id", id)
	return nil
}

func (r *ClientRepository) ListClients(ctx context.Context) ([]*models.Client, error) {
	logger.Info(ctx, "listing all clients")

	query, args, err := r.stmtBuilder.
		Select("id", "capacity", "rate_per_sec", "created_at", "updated_at").
		From("clients").
		OrderBy("created_at").
		ToSql()
	if err != nil {
		logger.Error(ctx, "sql build error on list", "err", err)
		return nil, err
	}

	var clients []*models.Client
	err = r.db.SelectContext(ctx, &clients, query, args...)
	if err != nil {
		logger.Error(ctx, "client list fetch failed", "err", err)
		return nil, err
	}

	logger.Info(ctx, "clients fetched", "count", len(clients))
	return clients, nil
}
