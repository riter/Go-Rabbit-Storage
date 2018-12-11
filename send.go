package main

import (
	"github.com/streadway/amqp"
)

func sendResponseStorage(response []byte) {
	conn, err := amqp.Dial(urlServe)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		queueResponse, // name
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
		queueResponse,
		queueResponse,
		exchange,
		false,
		nil,
	)

	failOnError(ec, "Failed to declare a exchange")


	err = ch.Publish(
		exchange,     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        response,
		})
	failOnError(err, "Failed to publish a message")

}