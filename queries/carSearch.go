package queries

import (
	"bids/models"
	"go.mongodb.org/mongo-driver/bson"
)

func GetOfferQuery(offer models.CarSearch) bson.M {
	query := bson.M{}

	if offer.Make != "" {
		query["car.make"] = offer.Make
	}

	if offer.Model != "" {
		query["car.model"] = offer.Model
	}

	if offer.PriceMin != 0 && offer.PriceMax != 0 {
		query["car.price"] = bson.M{"$gte": offer.PriceMin, "$lte": offer.PriceMax}
	} else if offer.PriceMin != 0 {
		query["car.price"] = bson.M{"$gte": offer.PriceMin}
	} else if offer.PriceMax != 0 {
		query["car.price"] = bson.M{"$lte": offer.PriceMax}
	}

	if offer.MileageMin != 0 && offer.MileageMax != 0 {
		query["car.mileage"] = bson.M{"$gte": offer.MileageMin, "$lte": offer.MileageMax}
	} else if offer.MileageMin != 0 {
		query["car.mileage"] = bson.M{"$gte": offer.MileageMin}
	} else if offer.MileageMax != 0 {
		query["car.mileage"] = bson.M{"$lte": offer.MileageMax}
	}

	if offer.YearMin != 0 && offer.YearMax != 0 {
		query["car.year"] = bson.M{"$gte": offer.YearMin, "$lte": offer.YearMax}
	} else if offer.YearMin != 0 {
		query["car.year"] = bson.M{"$gte": offer.YearMin}
	} else if offer.YearMax != 0 {
		query["car.year"] = bson.M{"$lte": offer.YearMax}
	}

	if offer.Type != "" {
		query["car.type"] = offer.Type
	}

	if offer.EngineCapacityMin != 0 && offer.EngineCapacityMax != 0 {
		query["car.engine_capacity"] = bson.M{"$gte": offer.EngineCapacityMin, "$lte": offer.EngineCapacityMax}
	} else if offer.EngineCapacityMin != 0 {
		query["car.engine_capacity"] = bson.M{"$gte": offer.EngineCapacityMin}
	} else if offer.EngineCapacityMax != 0 {
		query["car.engine_capacity"] = bson.M{"$lte": offer.EngineCapacityMax}
	}

	if offer.Fuel != "" {
		query["car.fuel"] = offer.Fuel
	}

	if offer.PowerMin != 0 && offer.PowerMax != 0 {
		query["car.power"] = bson.M{"$gte": offer.PowerMin, "$lte": offer.PowerMax}
	} else if offer.PowerMin != 0 {
		query["car.power"] = bson.M{"$gte": offer.PowerMin}
	} else if offer.PowerMax != 0 {
		query["car.power"] = bson.M{"$lte": offer.PowerMax}
	}

	if offer.Transmission != "" {
		query["car.transmission"] = offer.Transmission
	}

	if offer.Drive != "" {
		query["car.drive"] = offer.Drive
	}

	if offer.Steering != "" {
		query["car.steering"] = offer.Steering
	}

	if offer.Doors != 0 {
		query["car.doors"] = offer.Doors
	}

	if offer.Seats != 0 {
		query["car.seats"] = offer.Seats
	}

	if offer.Condition != "" {
		query["car.condition"] = offer.Condition
	}

	if offer.CoordinatesX != 0 && offer.CoordinatesY != 0 && offer.Distance != 0 {
		query["car.location"] = bson.M{
			"$geoWithin": bson.M{
				"$centerSphere": []interface{}{
					[]interface{}{offer.CoordinatesX, offer.CoordinatesY},
					offer.Distance / 6378100, // distance is in kilometers
				},
			},
		}
	}

	return query
}
