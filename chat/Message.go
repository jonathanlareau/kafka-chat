package chat

import (
	"io/ioutil"
	"net/http"
	"encoding/json"
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Message struct {
	Avatar string `json:"avatar"`
	Msg  string `json:"msg"`
}

// NewMessage to publish in the chat
func NewMessage(avatar, msg string) Message {
	return Message{avatar, msg}
}

func Uuid() string {
	id := uuid.New();
	return id.String()
}


func joinHandler(publisher Publisher) func(*gin.Context) {
	return func(c *gin.Context) {
		var request Message
		err := json.NewDecoder(c.Request.Body).Decode(&request)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Error": err})
		}

		message := NewMessage("", fmt.Sprintf("The avatar - %s - has joined the chat!", request.Avatar))

		if err := publisher.Publish(context.Background(), message); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Error": err})
		}

		c.JSON(http.StatusAccepted, gin.H{"message": "Message published to chat"})
	}
}

func publishHandler(publisher Publisher) func(*gin.Context) {
	return func(c *gin.Context) {
		var request Message
		err := json.NewDecoder(c.Request.Body).Decode(&request)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Error": err})
		}

		message := NewMessage(request.Avatar, request.Msg)

		if err := publisher.Publish(context.Background(), message); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Error": err})
		}

		c.JSON(http.StatusAccepted, gin.H{"message": "Message published to chat"})
	}
}

// Do call server
func Do(endpoint string, data Message) error {
	requestBody, err := json.Marshal(data)
	if err != nil {
		return err
	}

	res, err := http.Post(endpoint, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusAccepted {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}

		return errors.New(string(body))
	}

	return nil
}

