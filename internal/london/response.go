package london

import "time"

type errorResponse struct {
	StatusCode int    `xml:"statusCode"`
	Error      string `xml:"error"`
}

type tireChangeBookingResponse struct {
	UUID string    `xml:"uuid"`
	Time time.Time `xml:"time"`
}

func newTireChangeTimeResponse(UUID string, time time.Time) *tireChangeBookingResponse {
	return &tireChangeBookingResponse{UUID: UUID, Time: time.UTC()}
}

type tireChangeTimesResponse struct {
	AvailableTimes []*tireChangeBookingResponse `xml:"availableTime"`
}

func newTireChangeTimesResponse(entities []*tireChangeTimeEntity) *tireChangeTimesResponse {
	var availableTimes []*tireChangeBookingResponse

	for _, entity := range entities {
		availableTimes = append(availableTimes, newTireChangeTimeResponse(entity.UUID, entity.Time))
	}

	return &tireChangeTimesResponse{AvailableTimes: availableTimes}
}
