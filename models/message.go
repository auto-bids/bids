package models

type Message struct {
	Message     string `json:"message" bson:"message"`
	Sender      string `json:"sender" bson:"sender"`
	Destination string `json:"destination" bson:"destination"`
	Options     string `json:"options" bson:"options"`
	Offer       Offer  `json:"offer" bson:"offer"`
}
