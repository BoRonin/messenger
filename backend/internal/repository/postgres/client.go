package repository

import (
	"context"
	"messanger/infrastructure/postgres"
	"messanger/internal/entity"
)

type PClient struct {
	DB *postgres.Postgres
}

func NewDBClient(db *postgres.Postgres) *PClient {
	return &PClient{
		DB: db,
	}
}

func (DBC *PClient) CreateClient(ctx context.Context, query *entity.Client) error {
	stmt, args, err := DBC.DB.Builder.Insert("clients").Columns("tel", "provider", "tag", "tz").Values(query.Tel, query.Provider, query.Tag, query.Tz).Suffix("returning id").ToSql()
	if err != nil {
		return err
	}
	if err = DBC.DB.Pool.QueryRow(ctx, stmt, args...).Scan(&query.ID); err != nil {
		return err
	}
	return nil
}
func (DBC *PClient) DeleteClient(ctx context.Context, query entity.Client) error {
	stmt := `delete from clients where id = $1`

	_, err := DBC.DB.Pool.Exec(ctx, stmt, query.ID)
	if err != nil {
		return err
	}
	return nil
}
func (DBC *PClient) UpdateClient(ctx context.Context, query entity.Client) error {
	stmt := `update clients set (tel, provider, tag, tz) = ($1, $2, $3, $4)
	WHERE id = $5`

	_, err := DBC.DB.Pool.Exec(ctx, stmt, query.Tel, query.Provider, query.Tag, query.Tz, query.ID)

	if err != nil {
		return err
	}
	return nil
}

func (DBC *PClient) GetClientsByTag(ctx context.Context, tag string) ([]entity.Client, error) {
	stmt := `select id, tel, provider, tz from clients where tag = $1`

	rows, err := DBC.DB.Pool.Query(ctx, stmt, tag)
	if err != nil {
		return nil, err
	}
	var clients []entity.Client
	for rows.Next() {
		var client entity.Client
		client.Tag = tag
		if err := rows.Scan(&client.ID, &client.Tel, &client.Provider, &client.Tz); err != nil {
			println(err)
		}
		clients = append(clients, client)
	}
	return clients, nil
}
