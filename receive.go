package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/streadway/amqp"
)

const directoryPath = "./public"

func receiveStorage() {
	conn, err := amqp.Dial(urlServe)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		queueRequest, // name
		true,    // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)

	failOnError(err, "Failed to declare a queue")

	e := ch.ExchangeDeclare(
		exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
		)
	failOnError(e, "Failed to declare a exchange")

	ec := ch.QueueBind(
		queueRequest,
		queueRequest,
		exchange,
		false,
		nil,
		)

	failOnError(ec, "Failed to declare a exchange")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	go func() {
		for d := range msgs {

			var requestStorage RequestStorage
			errDecode := json.Unmarshal(d.Body, &requestStorage)

			failOnError(errDecode, "Failed decode message")

			switch requestStorage.Action {
			case "post":
				processMessageCreateFile(requestStorage)
				break
			case "delete":
				processMessageDeleteFile(requestStorage)
				break
			}

			log.Printf("Received a message: %s", d.Body)

			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	forever := make(chan bool)
	<-forever
}

func processMessageCreateFile(data RequestStorage) {
	path, err := createDocument(data.Filename)

	if err != nil {
		response, _ :=  json.Marshal(ResponseStorage{Status:404, Message:"Error upload file", Action:"create", Url:""})
		go sendResponseStorage(response)
	} else {
		response, _ := json.Marshal(ResponseStorage{Status: 200, Message:"", Action:"create", Url:path})
		go sendResponseStorage(response)
	}

	log.Printf("File create: %s", path)
}

func processMessageDeleteFile(data RequestStorage) {
	path, err := deleteDocument(data.Filename)

	if err != nil {
		response, _ :=  json.Marshal(ResponseStorage{Status:404, Message:"Error: file not exist", Action:"create", Url:""})
		go sendResponseStorage(response)
	} else {
		response, _ := json.Marshal(ResponseStorage{Status: 200, Message:"", Action:"delete", Url:path})
		go sendResponseStorage(response)
	}

	log.Printf("File delete: %s", path)
}

func createDocument(filename string) (string, error) {
	out, err := os.Create(directoryPath + "/" + filename)
	defer out.Close()

	if err != nil {
		return "", err
	}

	path, _ := filepath.Abs(directoryPath + "/" + filename)
	return path, nil
}

func deleteDocument(filename string) (string, error) {

	err := os.Remove(directoryPath + "/" + filename)
	if err != nil {
		return "", err
	}

	path, _ := filepath.Abs(directoryPath + "/" + filename)
	return path, nil
}
