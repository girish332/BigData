package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/girish332/bigdata/models"
	"log"

	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	log.Println("Starting to consume messages from the queue")
	// Connect to RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	fmt.Println("Connected to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"plan_queue", // name
		false,        // durable
		true,         // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name,       // queue
		"myConsumer", // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	failOnError(err, "Failed to register a consumer")

	// Connect to Elasticsearch
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200",
		},
	}
	es, err := elasticsearch.NewClient(cfg)
	failOnError(err, "Failed to create the Elasticsearch client")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)

			// Deserialize the object
			var plan models.Plan
			err := json.Unmarshal(d.Body, &plan)
			failOnError(err, "Failed to deserialize Plan object")

			// Index the document
			req := esapi.IndexRequest{
				Index:      "plans",
				DocumentID: plan.ObjectId,
				Body:       bytes.NewReader(d.Body),
				Refresh:    "true",
			}

			// Perform the request with the client.
			res, err := req.Do(context.Background(), es)
			failOnError(err, "Error getting response")
			defer res.Body.Close()

			if res.IsError() {
				log.Printf("[%s] Error indexing document ID=%s", res.Status(), plan.ObjectId)
			} else {
				log.Printf("[%s] Successfully indexed document ID=%s", res.Status(), plan.ObjectId)
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
