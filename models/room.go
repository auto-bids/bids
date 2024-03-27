package models

type RoomDB struct {
	Id       string   `bson:"_id" json:"_id"`
	Name     string   `bson:"name" json:"name"`
	Users    []string `bson:"users" json:"users"`
	Messages []Offer  `bson:"messages" json:"messages"`
}

type PostRoomDB struct {
	Name     string   `bson:"name" json:"name"`
	Users    []string `bson:"users" json:"users" validate:"required,min=1,max=1"`
	Messages []Offer  `bson:"messages" json:"messages"`
}
