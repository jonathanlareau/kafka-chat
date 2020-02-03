package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/jonathanlareau/kafka-chat/chat"
)

func main() {

	var (
		brokers = "kafka-zookeeper:9092"
		topic   = "chat"
	)

	publisher := chat.NewPublisher(strings.Split(brokers, ","), topic)

	r := gin.Default()
	r.POST("/join", joinHandler(publisher))
	r.POST("/publish", publishHandler(publisher))

	_ = r.Run()
}
