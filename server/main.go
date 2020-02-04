/**
 * jonathan.lareau@gmail.com
 *
 * Very simple client example
 *
 * Receive rest call
 * Call the server in grpc
 * Return response in json
 *
 **/

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/jonathanlareau/kafka-chat/chat"
	"golang.org/x/net/context"
)

func main() {

	var (
		brokers = "kafka-zookeeper:9092"
		topic   = "chat"
	)

	publisher := chat.NewPublisher(strings.Split(brokers, ","), topic)

	routes := mux.NewRouter()
	routes.HandleFunc("/", indexHandler).Methods("GET")
	routes.HandleFunc("/api/sayhello/{msg}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UFT-8")

		vars := mux.Vars(r)

		message := chat.NewMessage("Angular", vars["msg"])

		if err := publisher.Publish(context.Background(), message); err != nil {
			log.Fatal(err)
		}

	}).Methods("GET")

	routes.HandleFunc("/api/readhello/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UFT-8")
		chMsg := make(chan chat.Message)
		chErr := make(chan error)
		consumer := chat.NewConsumer(strings.Split(brokers, ","), topic)

		go func() {
			consumer.Read(context.Background(), chMsg, chErr)
		}()
		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		response := ""
		for {
			select {
			case <-quit:
				response = "quit"
				goto end
			case m := <-chMsg:
				response = m.Msg
				goto end
			case err := <-chErr:
				log.Println(err)
				goto end
			}
		}
	end:
		json.NewEncoder(w).Encode(response)
	}).Methods("GET")

	//Listening for Rest calls
	fmt.Println("Application is running on : 8080 .....")
	http.ListenAndServe(":8080", routes)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UFT-8")
	json.NewEncoder(w).Encode("Server is running")
}
