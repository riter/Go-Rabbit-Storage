package main

import "log"

type ResponseStorage struct {
	Status int
	Message string
	Url   string
	Action string
}

type RequestStorage struct {
	Filename   string
	Action string
}

const exchange = "riter-storage"
const queueRequest = "storage-request"
const queueResponse = "storage-response"

const urlServe = "amqp://guest:guest@localhost:5672/"

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	receiveStorage()
}