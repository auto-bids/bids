package models

import "time"

type Auction struct {
	Id      string    `bson:"_id" json:"_id"`
	Owner   string    `bson:"owner" json:"owner" validate:"required,email"`
	End     time.Time `bson:"end" json:"end"`
	Created time.Time `bson:"created" json:"created" `
	Bidders []string  `bson:"users" json:"users"`
	Car     Car       `json:"car" bson:"car" `
}
type PostAuction struct {
	Owner   string    `bson:"owner" json:"owner" validate:"required,email"`
	End     time.Time `bson:"end" json:"end"`
	Created time.Time `bson:"created" json:"created" `
	Bidders []string  `bson:"users" json:"users"`
	Car     Car       `json:"car" bson:"car" `
}
type UpdateAuction struct {
	Id  string    `bson:"_id" json:"_id"`
	End time.Time `bson:"end" json:"end" validate:"required,datetime"`
	Car Car       `json:"car" bson:"car" validate:"required"`
}
