package repository

import (
	"context"
	"messanger/infrastructure/postgres"
	"messanger/internal/entity"
	"time"
)

type PMessage struct {
	DB *postgres.Postgres
}

func NewDBMessage(db *postgres.Postgres) *PMessage {
	return &PMessage{
		DB: db,
	}
}

func (pm *PMessage) CreateMessage(ctx context.Context, query *entity.Message) error {
	stmt, args, err := pm.DB.Builder.Insert("messages").Columns("time_sent", "status", "messanger_id", "client_id").Values(time.Now(), false, query.Messanger.ID, query.Client.ID).Suffix("returning id").ToSql()
	if err != nil {
		return err
	}
	if err = pm.DB.Pool.QueryRow(ctx, stmt, args...).Scan(&query.ID); err != nil {
		return err
	}
	return nil
}

func (pm *PMessage) DeleteMessage(ctx context.Context, query entity.Message) error {
	stmt := `delete from messages where id = $1`
	_, err := pm.DB.Pool.Exec(ctx, stmt, query.ID)
	if err != nil {
		return err
	}
	return nil
}
func (pm *PMessage) UpdateMessage(ctx context.Context, query entity.Message) error {
	stmt := `update messages set ("message_text", "filter") = ($1, $2)
	WHERE id = $3`

	_, err := pm.DB.Pool.Exec(ctx, stmt, query.Messanger.MessageText, query.Messanger.Filter)

	if err != nil {
		return err
	}
	return nil
}

func (pm *PMessage) SetStatusToTrue(ctx context.Context, query uint) error {
	stmt := `update messages set "status" = true
	WHERE id = $1`
	_, err := pm.DB.Pool.Exec(ctx, stmt, query)
	if err != nil {
		return err
	}
	return nil
}

func (pm *PMessage) GetUnsentMessages(ctx context.Context, query entity.Messanger) ([]entity.Message, error) {
	stms := `select m.id, time_sent, status, client_id, c.tel, c.provider, c.tag, c.tz from messages m
	JOIN clients c on c.id = client_id
	WHERE messanger_id = $1 AND status = false`

	rows, err := pm.DB.Pool.Query(ctx, stms, query.ID)
	if err != nil {
		return nil, err
	}
	var messages []entity.Message
	for rows.Next() {
		var mes entity.Message
		if err := rows.Scan(&mes.ID, &mes.TimeSent,
			&mes.Status, &mes.Client.ID, &mes.Client.Tel,
			&mes.Client.Provider, &mes.Client.Tag, &mes.Client.Tz); err != nil {
			return nil, err
		}
		messages = append(messages, mes)
	}
	return messages, nil
}
func (pm *PMessage) GetMessagesById(ctx context.Context, id uint) ([]entity.Message, error) {
	stms := `select m.id, time_sent, status, client_id, c.tel, c.provider, c.tag, c.tz, 
	r.id, r.time_start, r.time_end, r.message_text, r.filter from messages m
	JOIN clients c on c.id = client_id
	JOIN messangers r on r.id = messanger_id
	WHERE messanger_id = $1`
	rows, err := pm.DB.Pool.Query(ctx, stms, id)
	if err != nil {
		return nil, err
	}
	var messages []entity.Message
	for rows.Next() {
		var mes entity.Message
		if err := rows.Scan(&mes.ID, &mes.TimeSent,
			&mes.Status, &mes.Client.ID, &mes.Client.Tel,
			&mes.Client.Provider, &mes.Client.Tag, &mes.Client.Tz,
			&mes.Messanger.ID, &mes.Messanger.TimeStart, &mes.Messanger.TimeEnd,
			&mes.Messanger.MessageText, &mes.Messanger.Filter); err != nil {
			return nil, err
		}
		messages = append(messages, mes)
	}
	return messages, nil
}
func (pm *PMessage) GetMessages(ctx context.Context) ([]entity.Message, error) {
	stms := `select m.id, time_sent, status, client_id, c.tel, c.provider, c.tag, c.tz, 
	r.id, r.time_start, r.time_end, r.message_text, r.filter from messages m
	JOIN clients c on c.id = client_id
	JOIN messangers r on r.id = messanger_id`
	rows, err := pm.DB.Pool.Query(ctx, stms)
	if err != nil {
		return nil, err
	}
	var messages []entity.Message
	for rows.Next() {
		var mes entity.Message
		if err := rows.Scan(&mes.ID, &mes.TimeSent,
			&mes.Status, &mes.Client.ID, &mes.Client.Tel,
			&mes.Client.Provider, &mes.Client.Tag, &mes.Client.Tz,
			&mes.Messanger.ID, &mes.Messanger.TimeStart, &mes.Messanger.TimeEnd,
			&mes.Messanger.MessageText, &mes.Messanger.Filter); err != nil {
			return nil, err
		}
		messages = append(messages, mes)
	}
	return messages, nil
}
