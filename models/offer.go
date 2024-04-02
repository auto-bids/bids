package models

type Offer struct {
	Sender string  `bson:"Sender" json:"Sender"`
	Price  float32 `bson:"Offer" json:"Offer"`
	Time   int64   `bson:"Time" json:"Time"`
}
type MessageUnwindDB struct {
	Messages Offer `bson:"messages" json:"messages"`
}
