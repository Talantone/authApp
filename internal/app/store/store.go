package store

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type Store struct {
	config *Config
	Client *mongo.Client
}

func New(config *Config) *Store {
	return &Store{
		config: config,
	}
}

func (s *Store) Open() error {
	clientOptions := options.Client().ApplyURI(s.config.DatabaseURL)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	s.Client = client

	return nil
}

func (s *Store) Close() {
	err := s.Client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}
}
