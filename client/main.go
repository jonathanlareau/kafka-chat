package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"


	"github.com/jonathanlareau/kafka-chat/chat"
)

var myAvatar string

func main() {

	var (
		brokers = "kafka-zookeeper:9092"
		topic   = "chat"
	)

	publisher := chat.NewPublisher(strings.Split(brokers, ","), topic)

	fmt.Println("Choose the name of your avatar:")
	_, _ = fmt.Scanf("%s\n", &myAvatar)
	message := chat.NewMessage("", fmt.Sprintf("The avatar - %s - has joined the chat!", myAvatar))
	fmt.Printf("You have joined this chat as name: %s \n", myAvatar)
	if err := publisher.Publish(context.Background(), message); err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			reader := bufio.NewReader(os.Stdin)
			msg, _ := reader.ReadString('\n')
			message := chat.Message{Avatar: myAvatar, Msg: msg}
			if err := publisher.Publish(context.Background(), message); err != nil {
				log.Fatal(err)
			}
		}
	}()

	chMsg := make(chan chat.Message)
	chErr := make(chan error)
	consumer := chat.NewConsumer(strings.Split(brokers, ","), topic)

	go func() {
		consumer.Read(context.Background(), chMsg, chErr)
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-quit:
			goto end
		case m := <-chMsg:
			printMessage(m)
		case err := <-chErr:
			log.Println(err)
		}
	}
end:

	fmt.Println("\nYou have quit the chat!")

}

func printMessage(m chat.Message) {

	switch {
	case m.Avatar == myAvatar:
		return
	case m.Avatar == "":
		fmt.Println(m.Msg)
	default:
		fmt.Printf("%s say: %s", m.Avatar, m.Msg)
	}
}

// for mode server only
func joinTheRoom(host string) error {
	d := chat.Message{Avatar: myAvatar}
	endpoint := fmt.Sprintf("%s/join", host)
	return chat.Do(endpoint, d)
}

// for mode server only
func publishMessage(host, msg string) error {
	d := chat.Message{Avatar: myAvatar, Msg: msg}
	endpoint := fmt.Sprintf("%s/publish", host)
	return chat.Do(endpoint, d)
}
