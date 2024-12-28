package data

import (
	"encoding/xml"
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

func FetchData() (*Data, error) {
	client := &http.Client{}

	resp, err := client.Get("https://api.entur.io/realtime/v1/rest/et?useOriginalId=true&maxSize=100000")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data := &Data{}
	decoder := xml.NewDecoder(resp.Body)
	err = decoder.Decode(data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
