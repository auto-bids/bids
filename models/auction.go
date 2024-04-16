package models

type Auction struct {
	Id           string   `bson:"_id" json:"_id"`
	Owner        string   `bson:"owner" json:"owner" validate:"required,email"`
	End          int64    `bson:"end" json:"end"`
	Start        int64    `bson:"start" json:"start"`
	StartFrom    int64    `bson:"startFrom" json:"startFrom"`
	Created      int64    `bson:"created" json:"created"`
	Bidders      []string `bson:"bidders" json:"bidders"`
	Offers       []Offer  `bson:"offers" json:"offers"`
	MinimalRaise int64    `bson:"minimalRaise" json:"minimalRaise"`
	Car          Car      `json:"car" bson:"car" `
}

type PostAuction struct {
	Owner        string   `bson:"owner" json:"owner" validate:"required,email"`
	End          int64    `bson:"end" json:"end" validate:"required"`
	Start        int64    `bson:"start" json:"start" validate:"required"`
	Created      int64    `bson:"created" json:"created" `
	Bidders      []string `bson:"bidders" json:"bidders"`
	Offers       []Offer  `bson:"offers" json:"offers"`
	MinimalRaise int64    `bson:"minimalRaise" json:"minimalRaise"`
	Car          Car      `json:"car" bson:"car" validate:"required"`
}
type GetAuctionShort struct {
	Id           string   `bson:"_id" json:"_id"`
	Owner        string   `bson:"owner" json:"owner" validate:"required,email"`
	End          int64    `bson:"end" json:"end"`
	Start        int64    `bson:"start" json:"start"`
	StartFrom    int64    `bson:"startFrom" json:"startFrom"`
	Created      int64    `bson:"created" json:"created"`
	Bidders      []string `bson:"bidders" json:"bidders"`
	MinimalRaise int64    `bson:"minimalRaise" json:"minimalRaise"`
	Car          Car      `json:"car" bson:"car" `
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
type EndAuction struct {
	Offer  Offer  `bson:"offer" json:"offer"`
	Status string `bson:"status" json:"status"`
}
type Status struct {
	Status string `bson:"status" json:"status" validate:"omitempty,oneof=ended ongoing"`
}
