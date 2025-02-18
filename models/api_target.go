package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type APITarget struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	Service   string             `bson:"service" json:"service"`
	APISchema string             `bson:"api_schema" json:"api_schema"` // URL or file path to the API schema
}

// APIAlert represents the security alerts for an API target
type APIAlert struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`       // User who owns the API target
	Service   string             `bson:"service" json:"service"`       // API service name
	APISchema string             `bson:"api_schema" json:"api_schema"` // OpenAPI schema URL
	Alerts    map[string]int     `bson:"alerts" json:"alerts"`         // Alert counts by risk level
	Timestamp time.Time          `bson:"timestamp" json:"timestamp"`   // Time of detection
}
