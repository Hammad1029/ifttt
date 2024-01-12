package models

import "time"

type MsgKafka struct {
	ClientId  string
	ApiName   string
	Data      map[string]string
	Timestamp time.Time
}
