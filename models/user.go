package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	//bd id
	ID        primitive.ObjectID `bson:"_id,omitempty json:"id"`
	Name      string             `bson:"_name json:"name,omitempty"`
	LastName  string             `bson:"_lastName json:"lastName,omitempty"`
	YearBirth time.Time          `bson:"_yearBirth json:"yearBirth,omitempty"`
	Email     string             `bson:"_email json:"email,omitempty"`
	Password  string             `bson:"_password json:"password,omitempty"`
	Avatar    string             `bson:"_avatar json:"avatar,omitempty"`
	Banner    string             `bson:"_banner json:"banner,omitempty"`
	Biography string             `bson:"_biography json:"biography,omitempty"`
	Location  string             `bson:"_location json:"location,omitempty"`
	WebSite   string             `bson:"_webSite json:"webSite,omitempty"`
}
