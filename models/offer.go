package models

type Offer struct {
	Sender string  `bson:"sender" json:"sender"`
	Price  float32 `bson:"offer" json:"offer"`
	Time   int64   `bson:"time" json:"time"`
}
