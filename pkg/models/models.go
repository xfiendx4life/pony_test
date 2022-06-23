package models

import "time"

type Message struct {
	ID        string `json:"id"`
	Data      string `json:"message"`
	TimeStamp time.Time
}
