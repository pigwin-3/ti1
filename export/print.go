package export

import (
	"encoding/json"
	"fmt"
	"log"
	"ti1/data"
)

func PrintData(data *data.Data) {
	fmt.Println(data.ServiceDelivery.ResponseTimestamp)
	fmt.Println(data.ServiceDelivery.EstimatedTimetableDelivery[0].EstimatedJourneyVersionFrame.RecordedAtTime)
	for _, journey := range data.ServiceDelivery.EstimatedTimetableDelivery[0].EstimatedJourneyVersionFrame.EstimatedVehicleJourney {
		if journey.RecordedAtTime != "" {
			fmt.Println(journey.RecordedAtTime)
		}
		if journey.LineRef != "" {
			fmt.Println(journey.LineRef)
		}
		if journey.DirectionRef != "" {
			fmt.Println(journey.DirectionRef)
		}
		if journey.DataSource != "" {
			fmt.Println(journey.DataSource)
		}

		if journey.FramedVehicleJourneyRef.DatedVehicleJourneyRef != "" {
			fmt.Println(journey.FramedVehicleJourneyRef.DatedVehicleJourneyRef)
		} else if journey.DatedVehicleJourneyRef != "" {
			fmt.Println(journey.DatedVehicleJourneyRef)
		} else {
			fmt.Println("evj." + journey.EstimatedVehicleJourneyCode)
		}

		if journey.VehicleMode != "" {
			fmt.Println(journey.VehicleMode)
		}
		if journey.FramedVehicleJourneyRef.DataFrameRef != "" {
			fmt.Println(journey.FramedVehicleJourneyRef.DataFrameRef)
		}
		if journey.OriginRef != "" {
			fmt.Println(journey.OriginRef)
		}
		if journey.DestinationRef != "" {
			fmt.Println(journey.DestinationRef)
		}
		if journey.OperatorRef != "" {
			fmt.Println(journey.OperatorRef)
		}
		if journey.VehicleRef != "" {
			fmt.Println(journey.VehicleRef)
		}
		if journey.Cancellation != "" {
			fmt.Println(journey.Cancellation)
		}

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

		// Print the JSON string for the current journey
		fmt.Println(string(jsonString))

		fmt.Println("Calls:")
		for _, recordedCall := range journey.RecordedCalls {
			for _, call := range recordedCall.RecordedCall {
				if call.StopPointRef != "" {
					fmt.Println("RecordedCall StopPointRef:", call.StopPointRef)
				}
				if call.Order != "" {
					fmt.Println("RecordedCall Order:", call.Order)
				}
				if call.Cancellation != "" {
					fmt.Println("RecordedCall Cancellation:", call.Cancellation)
				}
				if call.AimedDepartureTime != "" {
					fmt.Println("RecordedCall AimedDepartureTime:", call.AimedDepartureTime)
				}
				if call.ActualDepartureTime != "" {
					fmt.Println("RecordedCall ActualDepartureTime:", call.ActualDepartureTime)
				}
				if call.AimedArrivalTime != "" {
					fmt.Println("RecordedCall AimedArrivalTime:", call.AimedArrivalTime)
				}
				if call.ActualArrivalTime != "" {
					fmt.Println("RecordedCall ActualArrivalTime:", call.ActualArrivalTime)
				}
				if call.ExpectedArrivalTime != "" {
					fmt.Println("RecordedCall ExpectedArrivalTime:", call.ExpectedArrivalTime)
				}
				if call.ExpectedDepartureTime != "" {
					fmt.Println("RecordedCall ExpectedDepartureTime:", call.ExpectedDepartureTime)
				}

				jsonObjectRC := make(map[string]interface{})

				if call.StopPointName != "" {
					jsonObjectRC["StopPointName"] = call.StopPointName
				}
				if call.ArrivalPlatformName != "" {
					jsonObjectRC["ArrivalPlatformName"] = call.ArrivalPlatformName
				}
				if call.DeparturePlatformName != "" {
					jsonObjectRC["DeparturePlatformName"] = call.DeparturePlatformName
				}
				if call.PredictionInaccurate != "" {
					jsonObjectRC["PredictionInaccurate"] = call.PredictionInaccurate
				}
				if call.Occupancy != "" {
					jsonObjectRC["Occupancy"] = call.Occupancy
				}

				// Convert the JSON object to a JSON string
				jsonString, err := json.Marshal(jsonObjectRC)
				if err != nil {
					log.Fatal(err)
				}

				// Print the JSON string for the current journey
				fmt.Println(string(jsonString))
			}
		}
		for _, estimatedCall := range journey.EstimatedCalls {
			for _, call := range estimatedCall.EstimatedCall {
				if call.StopPointRef != "" {
					fmt.Println("EstimatedCall StopPointRef:", call.StopPointRef)
				}
				if call.Order != "" {
					fmt.Println("EstimatedCall Order:", call.Order)
				}
				if call.Cancellation != "" {
					fmt.Println("EstimatedCall Cancellation:", call.Cancellation)
				}
				if call.AimedDepartureTime != "" {
					fmt.Println("EstimatedCall AimedDepartureTime:", call.AimedDepartureTime)
				}
				if call.AimedArrivalTime != "" {
					fmt.Println("EstimatedCall AimedArrivalTime:", call.AimedArrivalTime)
				}
				if call.ExpectedArrivalTime != "" {
					fmt.Println("EstimatedCall ExpectedArrivalTime:", call.ExpectedArrivalTime)
				}
				if call.ExpectedDepartureTime != "" {
					fmt.Println("EstimatedCall ExpectedDepartureTime:", call.ExpectedDepartureTime)
				}

				jsonObjectEC := make(map[string]interface{})

				if call.StopPointName != "" {
					jsonObjectEC["StopPointName"] = call.StopPointName
				}
				if call.RequestStop != "" {
					jsonObjectEC["RequestStop"] = call.RequestStop
				}
				if call.DepartureStatus != "" {
					jsonObjectEC["DepartureStatus"] = call.DepartureStatus
				}
				if call.DeparturePlatformName != "" {
					jsonObjectEC["DeparturePlatformName"] = call.DeparturePlatformName
				}
				if call.DepartureBoardingActivity != "" {
					jsonObjectEC["DepartureBoardingActivity"] = call.DepartureBoardingActivity
				}
				if call.ArrivalStatus != "" {
					jsonObjectEC["ArrivalStatus"] = call.ArrivalStatus
				}
				if call.ArrivalPlatformName != "" {
					jsonObjectEC["ArrivalPlatformName"] = call.ArrivalPlatformName
				}
				if call.ArrivalBoardingActivity != "" {
					jsonObjectEC["ArrivalBoardingActivity"] = call.ArrivalBoardingActivity
				}
				if call.CallNote != "" {
					jsonObjectEC["CallNote"] = call.CallNote
				}
				if call.DestinationDisplay != "" {
					jsonObjectEC["DestinationDisplay"] = call.DestinationDisplay
				}
				if call.TimingPoint != "" {
					jsonObjectEC["TimingPoint"] = call.TimingPoint
				}
				if call.SituationRef != "" {
					jsonObjectEC["SituationRef"] = call.SituationRef
				}
				if call.PredictionInaccurate != "" {
					jsonObjectEC["PredictionInaccurate"] = call.PredictionInaccurate
				}
				if call.Occupancy != "" {
					jsonObjectEC["Occupancy"] = call.Occupancy
				}
				if call.DepartureStopAssignment.AimedQuayRef != "" {
					jsonObjectEC["DepartureAimedQuayRef"] = call.DepartureStopAssignment.AimedQuayRef
				}
				if call.DepartureStopAssignment.ExpectedQuayRef != "" {
					jsonObjectEC["DepartureExpectedQuayRef"] = call.DepartureStopAssignment.ExpectedQuayRef
				}
				if call.DepartureStopAssignment.ActualQuayRef != "" {
					jsonObjectEC["DepartureActualQuayRef"] = call.DepartureStopAssignment.ActualQuayRef
				}
				if call.Extensions.StopsAtAirport != "" {
					jsonObjectEC["StopsAtAirport"] = call.Extensions.StopsAtAirport
				}
				if call.ArrivalStopAssignment.AimedQuayRef != "" {
					jsonObjectEC["ArrivalAimedQuayRef"] = call.ArrivalStopAssignment.AimedQuayRef
				}
				if call.ArrivalStopAssignment.ExpectedQuayRef != "" {
					jsonObjectEC["ArrivalExpectedQuayRef"] = call.ArrivalStopAssignment.ExpectedQuayRef
				}
				if call.ArrivalStopAssignment.ActualQuayRef != "" {
					jsonObjectEC["ArrivalActualQuayRef"] = call.ArrivalStopAssignment.ActualQuayRef
				}
				if call.ExpectedDeparturePredictionQuality.PredictionLevel != "" {
					jsonObjectEC["ExpectedDeparturePredictionLevel"] = call.ExpectedDeparturePredictionQuality.PredictionLevel
				}

				if call.ExpectedArrivalPredictionQuality.PredictionLevel != "" {
					jsonObjectEC["ExpectedArrivalPredictionLevel"] = call.ExpectedArrivalPredictionQuality.PredictionLevel
				}

				jsonString, err := json.Marshal(jsonObjectEC)

				if err != nil {
					log.Fatal(err)
				}

				fmt.Println(string(jsonString))

			}
		}

	}
}
