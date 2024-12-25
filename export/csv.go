package export

import (
	"encoding/csv"
	"log"
	"os"
	"ti1/data"
)

func ExportToCSV(data *data.Data) {
	// Open the file for writing
	file, err := os.Create("output.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Create a CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write the header row
	header := []string{
		"RecordedAtTime",
		"LineRef",
		"DirectionRef",
		"DataSource",
		"DatedVehicleJourneyRef",
		"VehicleMode",
		"DataFrameRef",
		"OriginRef",
		"DestinationRef",
		"OperatorRef",
		"VehicleRef",
		"Cancellation",
		"OriginName",
		"DestinationName",
		"ProductCategoryRef",
		"ServiceFeatureRef",
		"Monitored",
		"JourneyPatternRef",
		"JourneyPatternName",
		"PublishedLineName",
		"DirectionName",
		"OriginAimedDepartureTime",
		"DestinationAimedArrivalTime",
		"BlockRef",
		"VehicleJourneyRef",
		"Occupancy",
		"DestinationDisplayAtOrigin",
		"ExtraJourney",
		"RouteRef",
		"GroupOfLinesRef",
		"ExternalLineRef",
		"InCongestion",
		"PredictionInaccurate",
		"JourneyNote",
		"Via",
	}
	writer.Write(header)

	// Write the data rows
	for _, journey := range data.ServiceDelivery.EstimatedTimetableDelivery[0].EstimatedJourneyVersionFrame.EstimatedVehicleJourney {
		row := []string{
			journey.RecordedAtTime,
			journey.LineRef,
			journey.DirectionRef,
			journey.DataSource,
			journey.FramedVehicleJourneyRef.DatedVehicleJourneyRef,
			journey.VehicleMode,
			journey.FramedVehicleJourneyRef.DataFrameRef,
			journey.OriginRef,
			journey.DestinationRef,
			journey.OperatorRef,
			journey.VehicleRef,
			journey.Cancellation,
			journey.OriginName,
			journey.DestinationName,
			journey.ProductCategoryRef,
			journey.ServiceFeatureRef,
			journey.Monitored,
			journey.JourneyPatternRef,
			journey.JourneyPatternName,
			journey.PublishedLineName,
			journey.DirectionName,
			journey.OriginAimedDepartureTime,
			journey.DestinationAimedArrivalTime,
			journey.BlockRef,
			journey.VehicleJourneyRef,
			journey.Occupancy,
			journey.DestinationDisplayAtOrigin,
			journey.ExtraJourney,
			journey.RouteRef,
			journey.GroupOfLinesRef,
			journey.ExternalLineRef,
			journey.InCongestion,
			journey.PredictionInaccurate,
			journey.JourneyNote,
			journey.Via.PlaceName,
		}
		writer.Write(row)
	}
}
