package export

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
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

	// counters
	var insertCount, updateCount, totalCount, estimatedCallInsertCount, estimatedCallUpdateCount, recordedCallInsertCount, recordedCallUpdateCount int
	var mu sync.Mutex

	// Create a worker pool
	numWorkers := 10
	jobs := make(chan []interface{}, numWorkers)
	results := make(chan struct {
		action string
		id     int
		err    error
	}, numWorkers)

	var wg sync.WaitGroup

	// Worker function
	worker := func() {
		defer wg.Done()
		for values := range jobs {
			id, action, err := database.InsertOrUpdateEstimatedVehicleJourney(db, values)
			results <- struct {
				action string
				id     int
				err    error
			}{action, id, err}
		}
	}

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker()
	}

	// Send jobs to workers
	go func() {
		for _, journey := range data.ServiceDelivery.EstimatedTimetableDelivery[0].EstimatedVehicleJourney {
			var values []interface{}
			var datedVehicleJourneyRef, otherJson string

			values = append(values, sid)
			values = append(values, journey.RecordedAtTime)
			values = append(values, journey.LineRef)
			values = append(values, strings.ToLower(journey.DirectionRef))
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

			jsonObject := make(map[string]interface{})
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

			jsonString, err := json.Marshal(jsonObject)
			if err != nil {
				log.Fatal(err)
			}
			otherJson = string(jsonString)
			values = append(values, otherJson)

			jobs <- values

			// Add the missing code here
			for _, estimatedCall := range journey.EstimatedCalls {
				for _, call := range estimatedCall.EstimatedCall {
					var estimatedValues []interface{}

					//1 estimatedvehiclejourney
					estimatedValues = append(estimatedValues, id)
					//2 order
					estimatedValues = append(estimatedValues, call.Order)
					//3 stoppointref
					estimatedValues = append(estimatedValues, call.StopPointRef)
					//4 aimeddeparturetime
					estimatedValues = append(estimatedValues, call.AimedDepartureTime)
					//5 expecteddeparturetime
					estimatedValues = append(estimatedValues, call.ExpectedDepartureTime)
					//6 aimedarrivaltime
					estimatedValues = append(estimatedValues, call.AimedArrivalTime)
					//7 expectedarrivaltime
					estimatedValues = append(estimatedValues, call.ExpectedArrivalTime)
					//8 cancellation
					estimatedValues = append(estimatedValues, call.Cancellation)

					//9 estimated_data (JSON)
					estimatedJsonObject := make(map[string]interface{})
					// data already logged
					if call.ExpectedDepartureTime != "" {
						estimatedJsonObject["ExpectedDepartureTime"] = call.ExpectedDepartureTime
					}
					if call.ExpectedArrivalTime != "" {
						estimatedJsonObject["ExpectedArrivalTime"] = call.ExpectedArrivalTime
					}
					if call.Cancellation != "" {
						estimatedJsonObject["Cancellation"] = call.Cancellation
					}
					// The rest
					if call.StopPointName != "" {
						estimatedJsonObject["StopPointName"] = call.StopPointName
					}
					if call.RequestStop != "" {
						estimatedJsonObject["RequestStop"] = call.RequestStop
					}
					if call.DepartureStatus != "" {
						estimatedJsonObject["DepartureStatus"] = call.DepartureStatus
					}
					if call.DeparturePlatformName != "" {
						estimatedJsonObject["DeparturePlatformName"] = call.DeparturePlatformName
					}
					if call.DepartureBoardingActivity != "" {
						estimatedJsonObject["DepartureBoardingActivity"] = call.DepartureBoardingActivity
					}
					if call.DepartureStopAssignment.AimedQuayRef != "" {
						estimatedJsonObject["DepartureStopAssignment.AimedQuayRef"] = call.DepartureStopAssignment.AimedQuayRef
					}
					if call.DepartureStopAssignment.ExpectedQuayRef != "" {
						estimatedJsonObject["DepartureStopAssignment.ExpectedQuayRef"] = call.DepartureStopAssignment.ExpectedQuayRef
					}
					if call.DepartureStopAssignment.ActualQuayRef != "" {
						estimatedJsonObject["DepartureStopAssignment.ActualQuayRef"] = call.DepartureStopAssignment.ActualQuayRef
					}
					if call.Extensions.StopsAtAirport != "" {
						estimatedJsonObject["Extensions.StopsAtAirport"] = call.Extensions.StopsAtAirport
					}
					if call.ArrivalStatus != "" {
						estimatedJsonObject["ArrivalStatus"] = call.ArrivalStatus
					}
					if call.ArrivalPlatformName != "" {
						estimatedJsonObject["ArrivalPlatformName"] = call.ArrivalPlatformName
					}
					if call.ArrivalBoardingActivity != "" {
						estimatedJsonObject["ArrivalBoardingActivity"] = call.ArrivalBoardingActivity
					}
					if call.ArrivalStopAssignment.AimedQuayRef != "" {
						estimatedJsonObject["ArrivalStopAssignment.AimedQuayRef"] = call.ArrivalStopAssignment.AimedQuayRef
					}
					if call.ArrivalStopAssignment.ExpectedQuayRef != "" {
						estimatedJsonObject["ArrivalStopAssignment.ExpectedQuayRef"] = call.ArrivalStopAssignment.ExpectedQuayRef
					}
					if call.ArrivalStopAssignment.ActualQuayRef != "" {
						estimatedJsonObject["ArrivalStopAssignment.ActualQuayRef"] = call.ArrivalStopAssignment.ActualQuayRef
					}
					if call.CallNote != "" {
						estimatedJsonObject["CallNote"] = call.CallNote
					}
					if call.DestinationDisplay != "" {
						estimatedJsonObject["DestinationDisplay"] = call.DestinationDisplay
					}
					if call.ExpectedDeparturePredictionQuality.PredictionLevel != "" {
						estimatedJsonObject["ExpectedDeparturePredictionQuality.PredictionLevel"] = call.ExpectedDeparturePredictionQuality.PredictionLevel
					}
					if call.ExpectedArrivalPredictionQuality.PredictionLevel != "" {
						estimatedJsonObject["ExpectedArrivalPredictionQuality.PredictionLevel"] = call.ExpectedArrivalPredictionQuality.PredictionLevel
					}
					if call.TimingPoint != "" {
						estimatedJsonObject["TimingPoint"] = call.TimingPoint
					}
					if call.SituationRef != "" {
						estimatedJsonObject["SituationRef"] = call.SituationRef
					}
					if call.PredictionInaccurate != "" {
						estimatedJsonObject["PredictionInaccurate"] = call.PredictionInaccurate
					}
					if call.Occupancy != "" {
						estimatedJsonObject["Occupancy"] = call.Occupancy
					}

					// Convert the JSON object to a JSON string
					jsonString, err := json.Marshal(estimatedJsonObject)
					if err != nil {
						log.Fatal(err)
					}
					estimatedValues = append(estimatedValues, string(jsonString))

					// Insert or update the record
					stringValues := make([]string, len(estimatedValues))
					for i, v := range stringValues {
						stringValues[i] = fmt.Sprintf("%v", v)
					}
					interfaceValues := make([]interface{}, len(stringValues))
					for i, v := range stringValues {
						interfaceValues[i] = v
					}
					id, action, err := database.InsertOrUpdateEstimatedCall(db, interfaceValues)
					if err != nil {
						fmt.Printf("Error inserting/updating estimated call: %v\n", err)
					} else {
						if 1 == 0 {
							fmt.Printf("Action: %s, ID: %d\n", action, id)
						}

						if action == "insert" {
							estimatedCallInsertCount++
						} else if action == "update" {
							estimatedCallUpdateCount++
						}
					}
				}
			}
			for _, recordedCall := range journey.RecordedCalls {
				for _, call := range recordedCall.RecordedCall {
					var recordedValues []interface{}

					//1 estimatedvehiclejourney
					recordedValues = append(recordedValues, id)
					//2 order
					recordedValues = append(recordedValues, call.Order)
					//3 stoppointref
					recordedValues = append(recordedValues, call.StopPointRef)
					//4 aimeddeparturetime
					recordedValues = append(recordedValues, call.AimedDepartureTime)
					//5 expecteddeparturetime
					recordedValues = append(recordedValues, call.ExpectedDepartureTime)
					//6 aimedarrivaltime
					recordedValues = append(recordedValues, call.AimedArrivalTime)
					//7 expectedarrivaltime
					recordedValues = append(recordedValues, call.ExpectedArrivalTime)
					//8 cancellation
					recordedValues = append(recordedValues, call.Cancellation)
					//9 actualdeparturetime
					recordedValues = append(recordedValues, call.ActualDepartureTime)
					//10 actualarrivaltime
					recordedValues = append(recordedValues, call.ActualArrivalTime)

					//11 recorded_data (JSON)
					recordedJsonObject := make(map[string]interface{})
					if call.StopPointName != "" {
						recordedJsonObject["StopPointName"] = call.StopPointName
					}
					if call.ArrivalPlatformName != "" {
						recordedJsonObject["ArrivalPlatformName"] = call.ArrivalPlatformName
					}
					if call.DeparturePlatformName != "" {
						recordedJsonObject["DeparturePlatformName"] = call.DeparturePlatformName
					}
					if call.PredictionInaccurate != "" {
						recordedJsonObject["PredictionInaccurate"] = call.PredictionInaccurate
					}
					if call.Occupancy != "" {
						recordedJsonObject["Occupancy"] = call.Occupancy
					}

					// Convert the JSON object to a JSON string
					jsonString, err := json.Marshal(recordedJsonObject)
					if err != nil {
						log.Fatal(err)
					}
					recordedValues = append(recordedValues, string(jsonString))

					// Insert or update the record
					stringValues := make([]string, len(recordedValues))
					for i, v := range stringValues {
						stringValues[i] = fmt.Sprintf("%v", v)
					}
					interfaceValues := make([]interface{}, len(stringValues))
					for i, v := range stringValues {
						interfaceValues[i] = v
					}

					id, action, err := database.InsertOrUpdateRecordedCall(db, interfaceValues)
					if err != nil {
						fmt.Printf("Error inserting/updating recorded call: %v\n", err)
					} else {
						if 1 == 0 {
							fmt.Printf("Action: %s, ID: %d\n", action, id)
						}

						if action == "insert" {
							recordedCallInsertCount++
							//fmt.Printf("Action: %s, ID: %d\n", action, id)
						} else if action == "update" {
							recordedCallUpdateCount++
						}
					}
				}
			}
		}
		close(jobs)
	}()

	// Collect results
	go func() {
		for result := range results {
			if result.err != nil {
				fmt.Printf("Error inserting/updating estimated vehicle journey: %v\n", result.err)
			} else {
				mu.Lock()
				if result.action == "insert" {
					insertCount++
				} else if result.action == "update" {
					updateCount++
				}
				totalCount = insertCount + updateCount
				mu.Unlock()

				if totalCount%1000 == 0 {
					fmt.Printf(
						"Inserts: %d, Updates: %d, Total: %d; estimatedCalls = I: %d U: %d; recordedCalls = I: %d U: %d\n",
						insertCount,
						updateCount,
						totalCount,
						estimatedCallInsertCount,
						estimatedCallUpdateCount,
						recordedCallInsertCount,
						recordedCallUpdateCount,
					)
				}
			}
		}
	}()

	wg.Wait()
	close(results)

	fmt.Printf(
		"DONE: Inserts: %d, Updates: %d, Total: %d; estimatedCalls = I: %d U: %d; recordedCalls = I: %d U: %d\n",
		insertCount,
		updateCount,
		totalCount,
		estimatedCallInsertCount,
		estimatedCallUpdateCount,
		recordedCallInsertCount,
		recordedCallUpdateCount,
	)

	serviceDeliveryJsonObject := make(map[string]interface{})
	serviceDeliveryJsonObject["Inserts"] = insertCount
	serviceDeliveryJsonObject["Updates"] = updateCount
	serviceDeliveryJsonObject["EstimatedCallInserts"] = estimatedCallInsertCount
	serviceDeliveryJsonObject["EstimatedCallUpdates"] = estimatedCallUpdateCount
	serviceDeliveryJsonObject["RecordedCallInserts"] = recordedCallInsertCount
	serviceDeliveryJsonObject["RecordedCallUpdates"] = recordedCallUpdateCount

	serviceDeliveryJsonString, err := json.Marshal(serviceDeliveryJsonObject)
	if err != nil {
		log.Fatal(err)
	}

	err = database.UpdateServiceDeliveryData(db, sid, string(serviceDeliveryJsonString))
	if err != nil {
		log.Fatal(err)
	}
}
