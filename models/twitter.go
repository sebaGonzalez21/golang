package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Twitter struct {
	UserId  string    `bson:"userId" json:"userId,omitempty"`
	Message string    `bson:"message" json:"message,omitempty"`
	Date    time.Time `bson:"date" json:"date,omitempty"`
}

type TwitterList struct {
	ID      primitive.ObjectID `bson:"_id" json:"_id,omitempty"`
	UserId  string             `bson:"userId" json:"userId,omitempty"`
	Message string             `bson:"message" json:"message,omitempty"`
	Date    time.Time          `bson:"date" json:"date,omitempty"`
}

type FollowTweets struct {
	ID             primitive.ObjectID `bson:"_id" json:"_id,omitempty"`
	UserRelationID string             `bson:"userRelationId" json:"userRelationId,omitempty"`
	UserId         string             `bson:"userId" json:"userId,omitempty"`
	Tweet          struct {
		ID      primitive.ObjectID `bson:"_id" json:"_id,omitempty"`
		Message string             `bson:"message" json:"message,omitempty"`
		Date    time.Time          `bson:"date" json:"date,omitempty"`
	}
}
