package models

type Offer struct {
	Sender string `bson:"sender" json:"sender"`
	Price  int64  `bson:"offer" json:"offer"`
	Time   int64  `bson:"time" json:"time"`
}
type OfferUnwind struct {
	Offers Offer `json:"offers" bson:"offers"`
}
