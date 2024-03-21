package models

type Auction struct {
	Id    string   `bson:"_id" json:"_id"`
	Owner string   `bson:"owner" json:"owner" validate:"required,email"`
	Users []string `bson:"users" json:"users"`
	Car   Car      `json:"car" bson:"car" validate:"required"`
}
