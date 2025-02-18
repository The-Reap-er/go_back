package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type URL struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID     primitive.ObjectID `bson:"user_id" json:"user_id"`
	Service    string             `bson:"service" json:"service"`
	URLList    string             `bson:"url_list" json:"url_list"`
	Scanner    string             `bson:"scanner,omitempty" json:"scanner,omitempty"`
	ScanResult string             `bson:"scan_result,omitempty" json:"scan_result,omitempty"`
}
