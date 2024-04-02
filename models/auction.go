package models

type Auction struct {
	Id           string   `bson:"_id" json:"_id"`
	Owner        string   `bson:"owner" json:"owner" validate:"required,email"`
	End          int64    `bson:"end" json:"end"`
	Start        int64    `bson:"start" json:"start"`
	Created      int64    `bson:"created" json:"created" `
	Bidders      []string `bson:"users" json:"users"`
	Status       string   `bson:"status" json:"status"`
	MinimalRaise float32
	//Car          Car `json:"car" bson:"car" `
}
type PostAuction struct {
	Owner   string   `bson:"owner" json:"owner" validate:"required,email"`
	End     int64    `bson:"end" json:"end"`
	Start   int64    `bson:"start" json:"start"`
	Created int64    `bson:"created" json:"created" `
	Bidders []string `bson:"users" json:"users"`
	//Car     Car      `json:"car" bson:"car" `
}
type GetAuctionForRoom struct {
	Owner string `bson:"owner" json:"owner" validate:"required,email"`
	End   int64  `bson:"end" json:"end"`
}
type UpdateAuction struct {
	Id  string `bson:"_id" json:"_id"`
	End int64  `bson:"end" json:"end" validate:"required,datetime"`
	Car Car    `json:"car" bson:"car" validate:"required"`
}
type AddBidder struct {
	Id      string   `bson:"_id" json:"_id"`
	Bidders []string `bson:"users" json:"users"`
}
