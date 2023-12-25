package migration

import "go.mongodb.org/mongo-driver/mongo"

// Model represents the migration model
type Model interface {
	Indices() []mongo.IndexModel
	Name() string
}
