package models

type Relation struct {
	UserID         string `bson:"userId" json:"userId,omitempty"`
	UserRelationId string `bson:"userRelationId" json:"userRelationId,omitempty"`
}
