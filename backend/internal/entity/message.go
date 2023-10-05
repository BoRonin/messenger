package entity

import "time"

type Message struct {
	ID        uint      `json:"id"`
	TimeSent  time.Time `json:"time_sent"`
	Status    bool      `json:"status"`
	Messanger Messanger `json:"messanger"`
	Client    Client    `json:"client"`
}
type MessageToGo struct {
	ID    uint   `json:"id"`
	Phone int    `json:"phone"`
	Text  string `json:"text"`
}
