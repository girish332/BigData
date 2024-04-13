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

	req := esapi.IndicesCreateRequest{
		Index: "plans",
	}

	res, err := req.Do(context.Background(), es)
	if err != nil {
		log.Fatalf("Error sending the request: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("Error creating the index: %s", res.String())
	} else {
		log.Printf("Index created successfully")
	}

	jsonData, err := json.Marshal(getMapping())
	failOnError(err, "Failed to serialize the mapping")
	req2 := esapi.IndicesPutMappingRequest{
		Index: []string{"plans"},
		Body:  bytes.NewReader(jsonData),
	}

	res2, err := req2.Do(context.Background(), es)
	failOnError(err, "Failed to create the index")
	defer res2.Body.Close()

	if res2.IsError() {
		log.Printf("Error creating the index: %s", res2.String())
	} else {
		log.Printf("Index created successfully")
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)

			// Deserialize the object
			var plan models.Plan
			err := json.Unmarshal(d.Body, &plan)
			failOnError(err, "Failed to deserialize Plan object")

			// Add the plan_join field to the plan object
			plan.PlanJoin = map[string]interface{}{
				"name": "plan",
			}

			// Serialize the plan object with the added plan_join fields
			planJSON, err := json.Marshal(plan)
			failOnError(err, "Failed to serialize Plan object")

			// Index the plan document
			req := esapi.IndexRequest{
				Index:      "plans",
				DocumentID: plan.ObjectId,
				Body:       bytes.NewReader(planJSON),
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

			// Index the planCostShares document
			plan.PlanCostShares.PlanJoin = map[string]interface{}{
				"name":   "planCostShares",
				"parent": plan.ObjectId,
			}

			// Serialize the planCostShares object with the added plan_join fields
			planCostSharesJSON, err := json.Marshal(plan.PlanCostShares)
			failOnError(err, "Failed to serialize PlanCostShares object")

			// Index the planCostShares document
			req = esapi.IndexRequest{
				Index:      "plans",
				DocumentID: plan.PlanCostShares.ObjectId,
				Body:       bytes.NewReader(planCostSharesJSON),
				Refresh:    "true",
				Routing:    plan.ObjectId,
			}

			// Perform the request with the client.
			res, err = req.Do(context.Background(), es)
			failOnError(err, "Error getting response")
			defer res.Body.Close()

			if res.IsError() {
				log.Printf("[%s] Error indexing document ID=%s", res.Status(), plan.PlanCostShares.ObjectId)
			} else {
				log.Printf("[%s] Successfully indexed document ID=%s", res.Status(), plan.PlanCostShares.ObjectId)
			}

			// Index each linkedPlanServices document
			for _, linkedPlanService := range plan.LinkedPlanServices {
				linkedPlanService.PlanJoin = map[string]interface{}{
					"name":   "linkedPlanServices",
					"parent": plan.ObjectId,
				}

				// Serialize the linkedPlanServices object with the added plan_join fields
				linkedPlanServiceJSON, err := json.Marshal(linkedPlanService)
				failOnError(err, "Failed to serialize LinkedPlanService object")

				// Index the linkedPlanServices document
				req := esapi.IndexRequest{
					Index:      "plans",
					DocumentID: linkedPlanService.ObjectId,
					Body:       bytes.NewReader(linkedPlanServiceJSON),
					Refresh:    "true",
					Routing:    plan.ObjectId,
				}

				// Perform the request with the client.
				res, err := req.Do(context.Background(), es)
				failOnError(err, "Error getting response")
				defer res.Body.Close()

				if res.IsError() {
					log.Printf("[%s] Error indexing document ID=%s", res.Status(), linkedPlanService.ObjectId)
				} else {
					log.Printf("[%s] Successfully indexed document ID=%s", res.Status(), linkedPlanService.ObjectId)
				}

				// Index the linkedService document
				linkedPlanService.LinkedService.PlanJoin = map[string]interface{}{
					"name":   "linkedService",
					"parent": linkedPlanService.ObjectId,
				}

				// Serialize the linkedService object with the added plan_join fields
				linkedServiceJSON, err := json.Marshal(linkedPlanService.LinkedService)
				failOnError(err, "Failed to serialize LinkedService object")

				// Index the linkedService document
				req = esapi.IndexRequest{
					Index:      "plans",
					DocumentID: linkedPlanService.LinkedService.ObjectId,
					Body:       bytes.NewReader(linkedServiceJSON),
					Refresh:    "true",
					Routing:    linkedPlanService.ObjectId,
				}

				// Perform the request with the client.
				res, err = req.Do(context.Background(), es)
				failOnError(err, "Error getting response")
				defer res.Body.Close()

				if res.IsError() {
					log.Printf("[%s] Error indexing document ID=%s", res.Status(), linkedPlanService.LinkedService.ObjectId)
				} else {
					log.Printf("[%s] Successfully indexed document ID=%s", res.Status(), linkedPlanService.LinkedService.ObjectId)
				}

				// Index the planserviceCostShares document
				linkedPlanService.PlanServiceCostShares.PlanJoin = map[string]interface{}{
					"name":   "planserviceCostShares",
					"parent": linkedPlanService.ObjectId,
				}

				// Serialize the planserviceCostShares object with the added plan_join fields
				planServiceCostSharesJSON, err := json.Marshal(linkedPlanService.PlanServiceCostShares)
				failOnError(err, "Failed to serialize PlanServiceCostShares object")

				// Index the planserviceCostShares document
				req = esapi.IndexRequest{
					Index:      "plans",
					DocumentID: linkedPlanService.PlanServiceCostShares.ObjectId,
					Body:       bytes.NewReader(planServiceCostSharesJSON),
					Refresh:    "true",
					Routing:    linkedPlanService.ObjectId,
				}

				// Perform the request with the client.
				res, err = req.Do(context.Background(), es)
				failOnError(err, "Error getting response")
				defer res.Body.Close()

				if res.IsError() {
					log.Printf("[%s] Error indexing document ID=%s", res.Status(), linkedPlanService.PlanServiceCostShares.ObjectId)
				} else {
					log.Printf("[%s] Successfully indexed document ID=%s", res.Status(), linkedPlanService.PlanServiceCostShares.ObjectId)
				}
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func getMapping() map[string]interface{} {
	return map[string]interface{}{
		//"mappings": map[string]interface{}{
		"properties": map[string]interface{}{
			"plan": map[string]interface{}{
				"properties": map[string]interface{}{
					"_org": map[string]interface{}{
						"type": "text",
					},
					"objectId": map[string]interface{}{
						"type": "keyword",
					},
					"objectType": map[string]interface{}{
						"type": "text",
					},
					"planType": map[string]interface{}{
						"type": "text",
					},
					"creationDate": map[string]interface{}{
						"type":   "date",
						"format": "MM-dd-yyyy",
					},
				},
			},
			"planCostShares": map[string]interface{}{
				"properties": map[string]interface{}{
					"copay": map[string]interface{}{
						"type": "long",
					},
					"deductible": map[string]interface{}{
						"type": "long",
					},
					"_org": map[string]interface{}{
						"type": "text",
					},
					"objectId": map[string]interface{}{
						"type": "keyword",
					},
					"objectType": map[string]interface{}{
						"type": "text",
					},
				},
			},
			"linkedPlanServices": map[string]interface{}{
				"properties": map[string]interface{}{
					"_org": map[string]interface{}{
						"type": "text",
					},
					"objectId": map[string]interface{}{
						"type": "keyword",
					},
					"objectType": map[string]interface{}{
						"type": "text",
					},
				},
			},
			"linkedService": map[string]interface{}{
				"properties": map[string]interface{}{
					"_org": map[string]interface{}{
						"type": "text",
					},
					"name": map[string]interface{}{
						"type": "text",
					},
					"objectId": map[string]interface{}{
						"type": "keyword",
					},
					"objectType": map[string]interface{}{
						"type": "text",
					},
				},
			},
			"planserviceCostShares": map[string]interface{}{
				"properties": map[string]interface{}{
					"copay": map[string]interface{}{
						"type": "long",
					},
					"deductible": map[string]interface{}{
						"type": "long",
					},
					"_org": map[string]interface{}{
						"type": "text",
					},
					"objectId": map[string]interface{}{
						"type": "keyword",
					},
					"objectType": map[string]interface{}{
						"type": "text",
					},
				},
			},
			"plan_join": map[string]interface{}{
				"type":                  "join",
				"eager_global_ordinals": "true",
				"relations": map[string]interface{}{
					"plan":               []string{"planCostShares", "linkedPlanServices"},
					"linkedPlanServices": []string{"linkedService", "planserviceCostShares"},
				},
			},
		},
		//},
	}
}
