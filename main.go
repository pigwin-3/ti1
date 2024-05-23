package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
)

type Data struct {
	XMLName         xml.Name `xml:"Siri"`
	ServiceDelivery struct {
		ResponseTimestamp          string `xml:"ResponseTimestamp"`
		ProducerRef                string `xml:"ProducerRef"`
		EstimatedTimetableDelivery []struct {
			ResponseTimestamp            string `xml:"ResponseTimestamp"`
			EstimatedJourneyVersionFrame struct {
				RecordedAtTime          string `xml:"RecordedAtTime"`
				EstimatedVehicleJourney []struct {
					RecordedAtTime         string `xml:"RecordedAtTime"`
					LineRef                string `xml:"LineRef"`
					DirectionRef           string `xml:"DirectionRef"`
					DatedVehicleJourneyRef string `xml:"DatedVehicleJourneyRef"`
					VehicleMode            string `xml:"VehicleMode"`
					OriginRef              string `xml:"OriginRef"`
					OriginName             string `xml:"OriginName"`
					DestinationRef         string `xml:"DestinationRef"`
					DestinationName        string `xml:"DestinationName"`
					OperatorRef            string `xml:"OperatorRef"`
					ProductCategoryRef     string `xml:"ProductCategoryRef"`
					ServiceFeatureRef      string `xml:"ServiceFeatureRef"`
					Monitored              string `xml:"Monitored"`
					DataSource             string `xml:"DataSource"`
					VehicleRef             string `xml:"VehicleRef"`
					EstimatedCalls         []struct {
						EstimatedCall []struct {
							StopPointRef              string `xml:"StopPointRef"`
							Order                     string `xml:"Order"`
							StopPointName             string `xml:"StopPointName"`
							RequestStop               string `xml:"RequestStop"`
							AimedDepartureTime        string `xml:"AimedDepartureTime"`
							ExpectedDepartureTime     string `xml:"ExpectedDepartureTime"`
							DepartureStatus           string `xml:"DepartureStatus"`
							DeparturePlatformName     string `xml:"DeparturePlatformName"`
							DepartureBoardingActivity string `xml:"DepartureBoardingActivity"`
							DepartureStopAssignment   struct {
								AimedQuayRef    string `xml:"AimedQuayRef"`
								ExpectedQuayRef string `xml:"ExpectedQuayRef"`
								ActualQuayRef   string `xml:"ActualQuayRef"`
							} `xml:"DepartureStopAssignment"`
							Extensions struct {
								StopsAtAirport string `xml:"StopsAtAirport"`
							} `xml:"Extensions"`
							AimedArrivalTime        string `xml:"AimedArrivalTime"`
							ExpectedArrivalTime     string `xml:"ExpectedArrivalTime"`
							ArrivalStatus           string `xml:"ArrivalStatus"`
							ArrivalPlatformName     string `xml:"ArrivalPlatformName"`
							ArrivalBoardingActivity string `xml:"ArrivalBoardingActivity"`
							ArrivalStopAssignment   struct {
								AimedQuayRef    string `xml:"AimedQuayRef"`
								ExpectedQuayRef string `xml:"ExpectedQuayRef"`
								ActualQuayRef   string `xml:"ActualQuayRef"`
							} `xml:"ArrivalStopAssignment"`
							CallNote                           string `xml:"CallNote"`
							Cancellation                       string `xml:"Cancellation"`
							DestinationDisplay                 string `xml:"DestinationDisplay"`
							ExpectedDeparturePredictionQuality struct {
								PredictionLevel string `xml:"PredictionLevel"`
							} `xml:"ExpectedDeparturePredictionQuality"`
							ExpectedArrivalPredictionQuality struct {
								PredictionLevel string `xml:"PredictionLevel"`
							} `xml:"ExpectedArrivalPredictionQuality"`
							TimingPoint          string `xml:"TimingPoint"`
							SituationRef         string `xml:"SituationRef"`
							PredictionInaccurate string `xml:"PredictionInaccurate"`
							Occupancy            string `xml:"Occupancy"`
						} `xml:"EstimatedCall"`
					} `xml:"EstimatedCalls"`
					IsCompleteStopSequence  string `xml:"IsCompleteStopSequence"`
					FramedVehicleJourneyRef struct {
						DataFrameRef           string `xml:"DataFrameRef"`
						DatedVehicleJourneyRef string `xml:"DatedVehicleJourneyRef"`
					} `xml:"FramedVehicleJourneyRef"`
					Cancellation                string `xml:"Cancellation"`
					JourneyPatternRef           string `xml:"JourneyPatternRef"`
					JourneyPatternName          string `xml:"JourneyPatternName"`
					PublishedLineName           string `xml:"PublishedLineName"`
					DirectionName               string `xml:"DirectionName"`
					OriginAimedDepartureTime    string `xml:"OriginAimedDepartureTime"`
					DestinationAimedArrivalTime string `xml:"DestinationAimedArrivalTime"`
					BlockRef                    string `xml:"BlockRef"`
					VehicleJourneyRef           string `xml:"VehicleJourneyRef"`
					RecordedCalls               []struct {
						RecordedCall []struct {
							StopPointRef          string `xml:"StopPointRef"`
							Order                 string `xml:"Order"`
							Cancellation          string `xml:"Cancellation"`
							AimedDepartureTime    string `xml:"AimedDepartureTime"`
							ActualDepartureTime   string `xml:"ActualDepartureTime"`
							AimedArrivalTime      string `xml:"AimedArrivalTime"`
							ActualArrivalTime     string `xml:"ActualArrivalTime"`
							StopPointName         string `xml:"StopPointName"`
							ArrivalPlatformName   string `xml:"ArrivalPlatformName"`
							ExpectedArrivalTime   string `xml:"ExpectedArrivalTime"`
							ExpectedDepartureTime string `xml:"ExpectedDepartureTime"`
							DeparturePlatformName string `xml:"DeparturePlatformName"`
							PredictionInaccurate  string `xml:"PredictionInaccurate"`
							Occupancy             string `xml:"Occupancy"`
						} `xml:"RecordedCall"`
					} `xml:"RecordedCalls"`
					Occupancy                   string `xml:"Occupancy"`
					DestinationDisplayAtOrigin  string `xml:"DestinationDisplayAtOrigin"`
					PredictionInaccurate        string `xml:"PredictionInaccurate"`
					EstimatedVehicleJourneyCode string `xml:"EstimatedVehicleJourneyCode"`
					ExtraJourney                string `xml:"ExtraJourney"`
					RouteRef                    string `xml:"RouteRef"`
					GroupOfLinesRef             string `xml:"GroupOfLinesRef"`
					ExternalLineRef             string `xml:"ExternalLineRef"`
					InCongestion                string `xml:"InCongestion"`
					JourneyNote                 string `xml:"JourneyNote"`
					Via                         struct {
						PlaceName string `xml:"PlaceName"`
					} `xml:"Via"`
				} `xml:"EstimatedVehicleJourney"`
			} `xml:"EstimatedJourneyVersionFrame"`
		} `xml:"EstimatedTimetableDelivery"`
	} `xml:"ServiceDelivery"`
}

func main() {
	// Fetch data from entur
	resp, err := http.Get("https://api.entur.io/realtime/v1/rest/et")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	data := &Data{} // Initialize your struct

	decoder := xml.NewDecoder(resp.Body) // Create a new XML decoder
	err = decoder.Decode(data)           // Decode the XML data into your struct
	if err != nil {
		log.Fatal(err)
	}

	if 1 == 0 {
		printData(data)
	}

}

func printData(data *Data) {
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
