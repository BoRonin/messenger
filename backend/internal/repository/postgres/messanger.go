package repository

import (
	"context"
	"fmt"
	"messanger/infrastructure/postgres"
	"messanger/internal/entity"
)

type PMessanger struct {
	DB *postgres.Postgres
}

func NewDBMessanger(db *postgres.Postgres) *PMessanger {
	return &PMessanger{
		DB: db,
	}
}

func (pm *PMessanger) CreateMessanger(ctx context.Context, query *entity.Messanger) error {
	stmt, args, err := pm.DB.Builder.Insert("messangers").Columns("time_start", "time_end", "message_text", "filter").Values(query.TimeStart, query.TimeEnd, query.MessageText, query.Filter).Suffix("returning id").ToSql()
	if err != nil {
		return err
	}
	if err = pm.DB.Pool.QueryRow(ctx, stmt, args...).Scan(&query.ID); err != nil {
		return err
	}
	return nil
}

func (pm *PMessanger) DeleteMessanger(ctx context.Context, query entity.Messanger) error {
	stmt := `delete from messangers where id = $1`

	_, err := pm.DB.Pool.Exec(ctx, stmt, query.ID)
	if err != nil {
		return err
	}
	return nil
}
func (pm *PMessanger) UpdateMessanger(ctx context.Context, query entity.Messanger) error {
	stmt := `update clients set ("message_text", "filter") = ($1, $2)
	WHERE id = $3`

	_, err := pm.DB.Pool.Exec(ctx, stmt, query.MessageText, query.Filter, query.ID)

	if err != nil {
		return err
	}
	return nil
}

func (pm *PMessanger) GetActiveMessangers(ctx context.Context) ([]entity.Messanger, error) {
	stmt := `select id, filter from messangers WHERE NOW() BETWEEN time_start AND time_end`

	rows, err := pm.DB.Pool.Query(ctx, stmt)
	if err != nil {
		return nil, err
	}
	ms := []entity.Messanger{}
	for rows.Next() {
		m := entity.Messanger{}
		if err := rows.Scan(&m.ID, &m.Filter); err != nil {
			fmt.Printf("err: %v\n", err)
		}
		ms = append(ms, m)
	}
	return ms, nil
}
func (pm *PMessanger) GetMessangers(ctx context.Context) ([]entity.Messanger, error) {
	stmt := `select id, filter, time_start, time_end from messangers`

	rows, err := pm.DB.Pool.Query(ctx, stmt)
	if err != nil {
		return nil, err
	}
	ms := []entity.Messanger{}
	for rows.Next() {
		m := entity.Messanger{}
		if err := rows.Scan(&m.ID, &m.Filter, &m.TimeStart, &m.TimeEnd); err != nil {
			fmt.Printf("err: %v\n", err)
		}
		ms = append(ms, m)
	}
	return ms, nil
}
