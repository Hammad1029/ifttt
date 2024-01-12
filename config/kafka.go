package config

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

var kafkaConnection *kafka.Conn

func startKafka() {
	conn, err := kafka.DialLeader(
		context.Background(),
		"tcp",
		fmt.Sprintf("%s:%s", GetConfigProp("kafka.host"), GetConfigProp("kafka.port")),
		GetConfigProp("kafka.topic"),
		0,
	)
	handleError(err)
	kafkaConnection = conn
}

func CreateKafkaMsg(msg any) {
	msgBArr, err := json.Marshal(msg)
	handleError(err)
	kafkaConnection.SetWriteDeadline(time.Now().Add(10 * time.Second))
	_, err = kafkaConnection.Write(msgBArr)
	handleError(err)
}
