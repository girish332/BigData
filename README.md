# INFO 7255 BigDataIndexing - Project to demonstrate Indexing of Structured JSON objects.


### Tech Stack
- Go (Gin-Gonic)
- Redis
- Elastic Search
- RabbitMQ

### Features
- Authentication using GCP OAuth2.0
- Validate request JSON object with JSON Schema
- Cache Server Response and validate cache using ETag
- Support POST, PUT, PATCH, GET and DELETE Http Methods for the REST API
- Store JSON Objects in Redis key-value store for data persistence
- Index the JSON Objects in Elastic Server for Search capabilities
- Queueing indexing requests to Elastic Server using RabbitMQ

### Data Flow
1. Generate OAuth token using authorization work flow
2. Validate further API requests using the received ID token
3. Create JSON Object using the `POST` HTTP method
4. Validate incoming JSON Object using the respective JSON Schema
5. De-Structure hierarchial JSON Object while storing in Redis key-value store
6. Enqueue object in RabbitMQ queue to index the object
7. Dequeue from RabbitMQ queue and index data in ElasticServer
8. Implement Search queries using Kibana Console to retrieve indexed data


### Steps to run:
1. Clone the repository
2. Run docker compose up -d (This will start Redis, ElasticSearch, RabbitMQ, Kibana)
3. Run the Go application using `go run main.go`
4. Run the Listener to listen to RabbitMQ queue using `go run listener/main.go`

### API Endpoints

- POST `/v1/plan` - Creates a new plan provided in the request body
- PUT `/v1/plan/{id}` - Updates an existing plan provided by the id
    - A valid Etag for the object should also be provided in the `If-Match` HTTP Request Header
- PATCH `/v1/plan/{id}` - Patches an existing plan provided by the id
    - A valid Etag for the object should also be provided in the `If-Match` HTTP Request Header
- GET `/v1/plan/{id}` - Fetches an existing plan provided by the id
    - An Etag for the object can be provided in the `If-None-Match` HTTP Request Header
    - If the request is successful, a valid Etag for the object is returned in the `ETag` HTTP Response Header
- DELETE `/v1/plan/{id}` - Deletes an existing plan provided by the id
    - A valid Etag for the object should also be provided in the `If-Match` HTTP Request Header
