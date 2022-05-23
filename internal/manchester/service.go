package manchester

import (
	log "github.com/sirupsen/logrus"
)

type tireChangeTimesService struct {
	repository *tireChangeTimeRepository
}

func newTireChangeTimesService(repository *tireChangeTimeRepository) *tireChangeTimesService {
	return &tireChangeTimesService{repository: repository}
}

func (s *tireChangeTimesService) get(query *tireChangeTimesSearchQuery) *tireChangeTimesResponse {
	log.Infof("fetching tire change times for query: %+v", query)
	tireChangeTimes := s.repository.allBySearchQuery(query)
	log.Infof("successfully fetched %d tire change times for query: %+v", len(tireChangeTimes), query)

	return newTireChangeTimesResponse(tireChangeTimes)
}

func (s *tireChangeTimesService) book(id uint, contactInformation string) (*tireChangeTimeBookingResponse, error) {
	log.Infof("trying to book tire change time with id: %d", id)
	tireChangeTime := s.repository.availableByID(id)

	if bookingErr := tireChangeTime.makeBooking(contactInformation); bookingErr != nil {
		return nil, bookingErr
	}

	tireChangeTime = s.repository.save(tireChangeTime)

	log.Infof("successfully booked tire change time with id: %d", id)
	return newTireChangeTimeResponse(tireChangeTime), nil
}
