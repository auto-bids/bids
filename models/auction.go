package models

type Auction struct {
	Id           string   `bson:"_id" json:"_id"`
	Owner        string   `bson:"owner" json:"owner" validate:"required,email"`
	End          int64    `bson:"end" json:"end"`
	Start        int64    `bson:"start" json:"start"`
	Created      int64    `bson:"created" json:"created"`
	Bidders      []string `bson:"bidders" json:"bidders"`
	Offers       []Offer  `bson:"offers" json:"offers"`
	MinimalRaise float32  `bson:"minimalRaise" json:"minimalRaise"`
	Car          Car      `json:"car" bson:"car" `
}

type PostAuction struct {
	Owner   string   `bson:"owner" json:"owner" validate:"required,email"`
	End     int64    `bson:"end" json:"end" validate:"required"`
	Start   int64    `bson:"start" json:"start" validate:"required"`
	Created int64    `bson:"created" json:"created" `
	Bidders []string `bson:"bidders" json:"bidders"`
	Offers  []Offer  `bson:"offers" json:"offers"`
	Car     Car      `json:"car" bson:"car" validate:"required"`
}
type GetAuctionShort struct {
	Id      string   `bson:"_id" json:"_id"`
	Owner   string   `bson:"owner" json:"owner" validate:"required,email"`
	End     int64    `bson:"end" json:"end"`
	Start   int64    `bson:"start" json:"start"`
	Created int64    `bson:"created" json:"created"`
	Bidders []string `bson:"bidders" json:"bidders"`
	Offers  []Offer  `bson:"offers" json:"offers"`
}

type GetAuctionForRoom struct {
	Owner  string  `bson:"owner" json:"owner" validate:"required,email"`
	Offers []Offer `bson:"offers" json:"offers"`
	End    int64   `bson:"end" json:"end"`
	Start  int64   `bson:"start" json:"start"`
}
type UpdateAuction struct {
	Id  string `bson:"_id" json:"_id"`
	End int64  `bson:"end" json:"end" validate:"required,datetime"`
	Car Car    `json:"car" bson:"car" validate:"required"`
}

type AddBidder struct {
	Id      string   `bson:"_id" json:"_id"`
	Bidders []string `bson:"bidders" json:"bidders"`
}
