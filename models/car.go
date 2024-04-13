package models

type Car struct {
	Title              string   `json:"title" bson:"title" validate:"required,max=40"`
	Make               string   `json:"make" bson:"make" validate:"required,max=30"`
	Model              string   `json:"model" bson:"model" validate:"required,max=30"`
	Price              int      `json:"price" bson:"price" validate:"required"`
	Description        string   `json:"description" bson:"description" validate:"required,max=3000"`
	Photos             []string `json:"photos" bson:"photos" validate:"required"`
	Year               int      `json:"year" bson:"year" validate:"required"`
	Mileage            int      `json:"mileage" bson:"mileage"`
	VinNumber          string   `json:"vin_number" bson:"vin_number"`
	EngineCapacity     int      `json:"engine_capacity" bson:"engine_capacity"`
	Fuel               string   `json:"fuel" bson:"fuel"`
	Transmission       string   `json:"transmission" bson:"transmission"`
	Steering           string   `json:"steering" bson:"steering"`
	Type               string   `json:"type" bson:"type"`
	Power              int      `json:"power" bson:"power"`
	Drive              string   `json:"drive" bson:"drive"`
	Doors              int      `json:"doors" bson:"doors"`
	Seats              int      `json:"seats" bson:"seats"`
	RegistrationNumber string   `json:"registration_number" bson:"registration_number"`
	FirstRegistration  string   `json:"first_registration" bson:"first_registration"`
	Condition          string   `json:"condition" bson:"condition"`
	TelephoneNumber    string   `json:"telephone_number" bson:"telephone_number"`
}
type CarSearch struct {
	Make              string  `json:"make" form:"make" validate:"max=30"`
	Model             string  `json:"model" form:"model" validate:"max=30"`
	PriceMin          int     `json:"price_min" form:"price_min"`
	PriceMax          int     `json:"price_max" form:"price_max"`
	MileageMin        int     `json:"mileage_min" form:"mileage_min"`
	MileageMax        int     `json:"mileage_max" form:"mileage_max"`
	YearMin           int     `json:"year_min" form:"year_min"`
	YearMax           int     `json:"year_max" form:"year_max"`
	Type              string  `json:"type" form:"type"`
	EngineCapacityMin int     `json:"engine_capacity_min" form:"engine_capacity_min"`
	EngineCapacityMax int     `json:"engine_capacity_max" form:"engine_capacity_max"`
	Fuel              string  `json:"fuel" form:"fuel"`
	PowerMin          int     `json:"power_min" form:"power_min"`
	PowerMax          int     `json:"power_max" form:"power_max"`
	Transmission      string  `json:"transmission" form:"transmission"`
	Drive             string  `json:"drive" form:"drive"`
	Steering          string  `json:"steering" form:"steering"`
	Doors             int     `json:"doors" form:"doors"`
	Seats             int     `json:"seats" form:"seats"`
	Condition         string  `json:"condition" form:"condition"`
	CoordinatesX      float32 `json:"lat" form:"lat"`
	CoordinatesY      float32 `json:"lng" form:"lng"`
	Distance          float64 `json:"radius" form:"radius"`
	FilterBy          string  `json:"filter_by" form:"filter_by"`
	SortDirection     int     `json:"sort_direction" form:"sort_direction"`
}
