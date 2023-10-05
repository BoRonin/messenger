package entity

import (
	"time"
)

type Messanger struct {
	ID          uint      `json:"id"`
	TimeStart   time.Time `json:"time_start"`
	TimeEnd     time.Time `json:"time_end"`
	MessageText string    `json:"message_text"`
	Filter      string    `json:"filter"`
}
