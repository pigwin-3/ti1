package export

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"sync/atomic"
	"ti1/config"
	"ti1/data"
	"ti1/database"
	"time"
)

// DBData is the main entry point for data processing
func DBData(data *data.Data) {
	DBDataOptimized(data)
}

// DBDataOptimized processes data with concurrent workers for better performance
func DBDataOptimized(data *data.Data) {
	fmt.Println(data.ServiceDelivery.ResponseTimestamp)
	fmt.Println(data.ServiceDelivery.EstimatedTimetableDelivery[0].EstimatedJourneyVersionFrame.RecordedAtTime)

	db, err := config.ConnectToPostgreSQL()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Connect to Valkey
	valkeyClient, err := config.ConnectToValkey("config/conf.json")
	if err != nil {
		log.Fatalf("Failed to connect to Valkey: %v", err)
	}
	defer config.DisconnectFromValkey(valkeyClient)

	ctx := context.Background()

	// Get service id aka sid
	sid, err := database.InsertServiceDelivery(db, data.ServiceDelivery.ResponseTimestamp, data.ServiceDelivery.EstimatedTimetableDelivery[0].EstimatedJourneyVersionFrame.RecordedAtTime)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("SID:", sid)

	// Record start time
	startTime := time.Now()

	// Atomic counters for thread-safe counting
	var insertCount, updateCount, estimatedCallInsertCount, estimatedCallUpdateCount, estimatedCallNoneCount, recordedCallInsertCount, recordedCallUpdateCount, recordedCallNoneCount int64

	journeys := data.ServiceDelivery.EstimatedTimetableDelivery[0].EstimatedJourneyVersionFrame.EstimatedVehicleJourney
	totalJourneys := len(journeys)
	fmt.Printf("Processing %d journeys...\n", totalJourneys)

	// Job structures
	type evjJob struct {
		index int
	}

	type callJob struct {
		evjID  int
		values []interface{}
	}

	// Channels
	workerCount := 20 // Adjust based on your database and CPU
	evjJobs := make(chan evjJob, workerCount*2)
	estimatedCallJobs := make(chan callJob, workerCount*10)
	recordedCallJobs := make(chan callJob, workerCount*10)

	var wg sync.WaitGroup
	var callWg sync.WaitGroup

	// Start Estimated Call workers
	for w := 0; w < workerCount; w++ {
		callWg.Add(1)
		go func() {
			defer callWg.Done()
			for job := range estimatedCallJobs {
				id, action, err := database.InsertOrUpdateEstimatedCall(ctx, db, job.values, valkeyClient)
				if err != nil {
					log.Printf("Error inserting/updating estimated call: %v\n", err)
					continue
				}
				if action == "insert" {
					atomic.AddInt64(&estimatedCallInsertCount, 1)
				} else if action == "update" {
					atomic.AddInt64(&estimatedCallUpdateCount, 1)
				} else if action == "none" {
					atomic.AddInt64(&estimatedCallNoneCount, 1)
				}
				_ = id
			}
		}()
	}

	// Start Recorded Call workers
	for w := 0; w < workerCount; w++ {
		callWg.Add(1)
		go func() {
			defer callWg.Done()
			for job := range recordedCallJobs {
				id, action, err := database.InsertOrUpdateRecordedCall(ctx, db, job.values, valkeyClient)
				if err != nil {
					log.Printf("Error inserting/updating recorded call: %v\n", err)
					continue
				}
				if action == "insert" {
					atomic.AddInt64(&recordedCallInsertCount, 1)
				} else if action == "update" {
					atomic.AddInt64(&recordedCallUpdateCount, 1)
				} else if action == "none" {
					atomic.AddInt64(&recordedCallNoneCount, 1)
				}
				_ = id
			}
		}()
	}

	// Start EVJ workers
	for w := 0; w < workerCount; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range evjJobs {
				journey := &journeys[job.index]

				// Prepare values
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

				// Create JSON object
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
					log.Printf("Error marshaling JSON: %v\n", err)
					continue
				}
				otherJson = string(jsonString)
				values = append(values, otherJson)

				// Insert or update EVJ
				id, action, err := database.InsertOrUpdateEstimatedVehicleJourney(db, values)
				if err != nil {
					log.Printf("Error inserting/updating estimated vehicle journey: %v\n", err)
					continue
				}

				if action == "insert" {
					atomic.AddInt64(&insertCount, 1)
				} else if action == "update" {
					atomic.AddInt64(&updateCount, 1)
				}

				// Progress reporting
				total := atomic.AddInt64(&insertCount, 0) + atomic.AddInt64(&updateCount, 0)
				if total%1000 == 0 {
					fmt.Printf(
						"EVJ - I: %d, U: %d, Total: %d; EstCalls - I: %d U: %d N: %d; RecCalls - I: %d U: %d N: %d\n",
						atomic.LoadInt64(&insertCount),
						atomic.LoadInt64(&updateCount),
						total,
						atomic.LoadInt64(&estimatedCallInsertCount),
						atomic.LoadInt64(&estimatedCallUpdateCount),
						atomic.LoadInt64(&estimatedCallNoneCount),
						atomic.LoadInt64(&recordedCallInsertCount),
						atomic.LoadInt64(&recordedCallUpdateCount),
						atomic.LoadInt64(&recordedCallNoneCount),
					)
				}

				// Process Estimated Calls
				for _, estimatedCall := range journey.EstimatedCalls {
					for _, call := range estimatedCall.EstimatedCall {
						var estimatedValues []interface{}

						estimatedValues = append(estimatedValues, id)
						estimatedValues = append(estimatedValues, call.Order)
						estimatedValues = append(estimatedValues, call.StopPointRef)
						estimatedValues = append(estimatedValues, call.AimedDepartureTime)
						estimatedValues = append(estimatedValues, call.ExpectedDepartureTime)
						estimatedValues = append(estimatedValues, call.AimedArrivalTime)
						estimatedValues = append(estimatedValues, call.ExpectedArrivalTime)
						estimatedValues = append(estimatedValues, call.Cancellation)

						// estimated_data JSON
						estimatedJsonObject := make(map[string]interface{})
						if call.ExpectedDepartureTime != "" {
							estimatedJsonObject["ExpectedDepartureTime"] = call.ExpectedDepartureTime
						}
						if call.ExpectedArrivalTime != "" {
							estimatedJsonObject["ExpectedArrivalTime"] = call.ExpectedArrivalTime
						}
						if call.Cancellation != "" {
							estimatedJsonObject["Cancellation"] = call.Cancellation
						}
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

						jsonString, err := json.Marshal(estimatedJsonObject)
						if err != nil {
							log.Printf("Error marshaling estimated call JSON: %v\n", err)
							continue
						}
						estimatedValues = append(estimatedValues, string(jsonString))

						// Convert to string values
						interfaceValues := make([]interface{}, len(estimatedValues))
						for i, v := range estimatedValues {
							interfaceValues[i] = fmt.Sprintf("%v", v)
						}

						// Send to worker pool
						estimatedCallJobs <- callJob{evjID: id, values: interfaceValues}
					}
				}

				// Process Recorded Calls
				for _, recordedCall := range journey.RecordedCalls {
					for _, call := range recordedCall.RecordedCall {
						var recordedValues []interface{}

						recordedValues = append(recordedValues, id)
						recordedValues = append(recordedValues, call.Order)
						recordedValues = append(recordedValues, call.StopPointRef)
						recordedValues = append(recordedValues, call.AimedDepartureTime)
						recordedValues = append(recordedValues, call.ExpectedDepartureTime)
						recordedValues = append(recordedValues, call.AimedArrivalTime)
						recordedValues = append(recordedValues, call.ExpectedArrivalTime)
						recordedValues = append(recordedValues, call.Cancellation)
						recordedValues = append(recordedValues, call.ActualDepartureTime)
						recordedValues = append(recordedValues, call.ActualArrivalTime)

						// recorded_data JSON
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

						jsonString, err := json.Marshal(recordedJsonObject)
						if err != nil {
							log.Printf("Error marshaling recorded call JSON: %v\n", err)
							continue
						}
						recordedValues = append(recordedValues, string(jsonString))

						// Convert to string values
						interfaceValues := make([]interface{}, len(recordedValues))
						for i, v := range recordedValues {
							interfaceValues[i] = fmt.Sprintf("%v", v)
						}

						// Send to worker pool
						recordedCallJobs <- callJob{evjID: id, values: interfaceValues}
					}
				}
			}
		}()
	}

	// Send all EVJ jobs
	for i := range journeys {
		evjJobs <- evjJob{index: i}
	}
	close(evjJobs)

	// Wait for EVJ processing to complete
	wg.Wait()

	// Close call job channels and wait for call processing to complete
	close(estimatedCallJobs)
	close(recordedCallJobs)
	callWg.Wait()

	// Record end time
	endTime := time.Now()

	// Print final stats
	fmt.Printf(
		"\nDONE: EVJ - Inserts: %d, Updates: %d, Total: %d\n"+
			"      EstimatedCalls - I: %d U: %d N: %d\n"+
			"      RecordedCalls  - I: %d U: %d N: %d\n",
		atomic.LoadInt64(&insertCount),
		atomic.LoadInt64(&updateCount),
		atomic.LoadInt64(&insertCount)+atomic.LoadInt64(&updateCount),
		atomic.LoadInt64(&estimatedCallInsertCount),
		atomic.LoadInt64(&estimatedCallUpdateCount),
		atomic.LoadInt64(&estimatedCallNoneCount),
		atomic.LoadInt64(&recordedCallInsertCount),
		atomic.LoadInt64(&recordedCallUpdateCount),
		atomic.LoadInt64(&recordedCallNoneCount),
	)

	// Create map to hold JSON
	serviceDeliveryJsonObject := make(map[string]interface{})
	serviceDeliveryJsonObject["Inserts"] = atomic.LoadInt64(&insertCount)
	serviceDeliveryJsonObject["Updates"] = atomic.LoadInt64(&updateCount)
	serviceDeliveryJsonObject["EstimatedCallInserts"] = atomic.LoadInt64(&estimatedCallInsertCount)
	serviceDeliveryJsonObject["EstimatedCallUpdates"] = atomic.LoadInt64(&estimatedCallUpdateCount)
	serviceDeliveryJsonObject["EstimatedCallNone"] = atomic.LoadInt64(&estimatedCallNoneCount)
	serviceDeliveryJsonObject["RecordedCallInserts"] = atomic.LoadInt64(&recordedCallInsertCount)
	serviceDeliveryJsonObject["RecordedCallUpdates"] = atomic.LoadInt64(&recordedCallUpdateCount)
	serviceDeliveryJsonObject["RecordedCallNone"] = atomic.LoadInt64(&recordedCallNoneCount)
	serviceDeliveryJsonObject["StartTime"] = startTime.Format(time.RFC3339)
	serviceDeliveryJsonObject["EndTime"] = endTime.Format(time.RFC3339)
	serviceDeliveryJsonObject["Duration"] = endTime.Sub(startTime).String()

	// Convert JSON object to JSON string
	serviceDeliveryJsonString, err := json.Marshal(serviceDeliveryJsonObject)
	if err != nil {
		log.Fatal(err)
	}

	// Update ServiceDelivery data in database
	err = database.UpdateServiceDeliveryData(db, sid, string(serviceDeliveryJsonString))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Finished with this ServiceDelivery!")
}
