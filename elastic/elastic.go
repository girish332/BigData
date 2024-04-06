package elastic

import (
	"github.com/elastic/go-elasticsearch/v8"
	"log"
)

type Client struct {
	ES *elasticsearch.Client
}

type Factory struct {
}

func NewElasticFactory() *Factory {
	return &Factory{}
}

func (f *Factory) NewClient(cfg elasticsearch.Config) (*Client, error) {
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Printf("Error creating the client: %s", err)
		return nil, err
	}

	return &Client{ES: es}, nil
}
