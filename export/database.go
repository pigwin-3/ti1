package export

import (
	"encoding/json"
	"fmt"
	"log"
	"ti1/config"
	"ti1/data"
	"ti1/database"
)

func DBData(data *data.Data) {
	fmt.Println(data.ServiceDelivery.ResponseTimestamp)
	fmt.Println(data.ServiceDelivery.EstimatedTimetableDelivery[0].EstimatedJourneyVersionFrame.RecordedAtTime)

	db, err := config.ConnectToPostgreSQL()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Get service id aka sid
	sid, err := database.InsertServiceDelivery(db, data.ServiceDelivery.ResponseTimestamp, data.ServiceDelivery.EstimatedTimetableDelivery[0].EstimatedJourneyVersionFrame.RecordedAtTime)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("SID:", sid)

	for _, journey := range data.ServiceDelivery.EstimatedTimetableDelivery[0].EstimatedJourneyVersionFrame.EstimatedVehicleJourney {
		var values []interface{}
		var datedVehicleJourneyRef, otherJson string

		values = append(values, sid)
		values = append(values, journey.RecordedAtTime)
		values = append(values, journey.LineRef)
		values = append(values, journey.DirectionRef)
		values = append(values, journey.DataSource)

		if journey.FramedVehicleJourneyRef.DatedVehicleJourneyRef != "" {
			datedVehicleJourneyRef = journey.FramedVehicleJourneyRef.DatedVehicleJourneyRef
		} else if journey.DatedVehicleJourneyRef != "" {
			datedVehicleJourneyRef = journey.DatedVehicleJourneyRef
		} else {
			datedVehicleJourneyRef = "evj." + journey.EstimatedVehicleJourneyCode
		}
		values = append(values, datedVehicleJourneyRef)

		values = append(values, journey.VehicleMode)
		values = append(values, journey.FramedVehicleJourneyRef.DataFrameRef)
		values = append(values, journey.OriginRef)
		values = append(values, journey.DestinationRef)
		values = append(values, journey.OperatorRef)
		values = append(values, journey.VehicleRef)
		values = append(values, journey.Cancellation)

		// Create a map to hold the JSON object for the current journey
		jsonObject := make(map[string]interface{})

		// Add relevant fields to the JSON object
		if journey.OriginName != "" {
			jsonObject["OriginName"] = journey.OriginName
		}
		if journey.DestinationName != "" {
			jsonObject["DestinationName"] = journey.DestinationName
		}
		if journey.ProductCategoryRef != "" {
			jsonObject["ProductCategoryRef"] = journey.ProductCategoryRef
		}
		if journey.ServiceFeatureRef != "" {
			jsonObject["ServiceFeatureRef"] = journey.ServiceFeatureRef
		}
		if journey.Monitored != "" {
			jsonObject["Monitored"] = journey.Monitored
		}
		if journey.JourneyPatternRef != "" {
			jsonObject["JourneyPatternRef"] = journey.JourneyPatternRef
		}
		if journey.JourneyPatternName != "" {
			jsonObject["JourneyPatternName"] = journey.JourneyPatternName
		}
		if journey.PublishedLineName != "" {
			jsonObject["PublishedLineName"] = journey.PublishedLineName
		}
		if journey.DirectionName != "" {
			jsonObject["DirectionName"] = journey.DirectionName
		}
		if journey.OriginAimedDepartureTime != "" {
			jsonObject["OriginAimedDepartureTime"] = journey.OriginAimedDepartureTime
		}
		if journey.DestinationAimedArrivalTime != "" {
			jsonObject["DestinationAimedArrivalTime"] = journey.DestinationAimedArrivalTime
		}
		if journey.BlockRef != "" {
			jsonObject["BlockRef"] = journey.BlockRef
		}
		if journey.VehicleJourneyRef != "" {
			jsonObject["VehicleJourneyRef"] = journey.VehicleJourneyRef
		}
		if journey.Occupancy != "" {
			jsonObject["Occupancy"] = journey.Occupancy
		}
		if journey.DestinationDisplayAtOrigin != "" {
			jsonObject["DestinationDisplayAtOrigin"] = journey.DestinationDisplayAtOrigin
		}
		if journey.ExtraJourney != "" {
			jsonObject["ExtraJourney"] = journey.ExtraJourney
		}
		if journey.RouteRef != "" {
			jsonObject["RouteRef"] = journey.RouteRef
		}
		if journey.GroupOfLinesRef != "" {
			jsonObject["GroupOfLinesRef"] = journey.GroupOfLinesRef
		}
		if journey.ExternalLineRef != "" {
			jsonObject["ExternalLineRef"] = journey.ExternalLineRef
		}
		if journey.InCongestion != "" {
			jsonObject["InCongestion"] = journey.InCongestion
		}
		if journey.PredictionInaccurate != "" {
			jsonObject["PredictionInaccurate"] = journey.PredictionInaccurate
		}
		if journey.JourneyNote != "" {
			jsonObject["JourneyNote"] = journey.JourneyNote
		}
		if journey.Via.PlaceName != "" {
			jsonObject["Via"] = journey.Via.PlaceName
		}

		// Convert the JSON object to a JSON string
		jsonString, err := json.Marshal(jsonObject)
		if err != nil {
			log.Fatal(err)
		}
		otherJson = string(jsonString)
		values = append(values, otherJson)

		// Insert or update the record
		err = database.InsertOrUpdateEstimatedVehicleJourney(db, values)
		if err != nil {
			fmt.Printf("Error inserting/updating estimated vehicle journey: %v\n", err)
		}
	}
}
