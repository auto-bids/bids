package models

type Message struct {
	Sender  string `bson:"Sender" json:"Sender"`
	float32 `bson:"Message" json:"Message"`
	Time    int64 `bson:"Time"`
}
type MessageUnwindDB struct {
	Messages Message `bson:"messages" json:"messages"`
}
