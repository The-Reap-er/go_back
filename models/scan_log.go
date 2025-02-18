package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ScanLog struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	URL       string             `bson:"url" json:"url"`
	ScanType  string             `bson:"scan_type" json:"scan_type"` // "spider", "active", or "alert"
	Status    string             `bson:"status" json:"status"`       // "success" or "failure"
	Message   string             `bson:"message" json:"message"`     // Additional message or result
	Timestamp int64              `bson:"timestamp" json:"timestamp"` // Time of the scan
}
