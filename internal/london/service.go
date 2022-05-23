package london

import (
	log "github.com/sirupsen/logrus"
	"time"
)

type tireChangeTimesService struct {
	repository *tireChangeTimeRepository
}

func newTireChangeTimesService(repository *tireChangeTimeRepository) *tireChangeTimesService {
	return &tireChangeTimesService{repository: repository}
}

func (s *tireChangeTimesService) getAvailable(from time.Time, until time.Time) (*tireChangeTimesResponse, error) {
	log.Infof("fetching tire change times from %s until %s", from, until)

	if !from.Equal(until) && until.Before(from) {
		return nil, newInvalidTirChangeTimesPeriodError(from, until)
	}

	tireChangeTimes := s.repository.availableByTimeRange(from, until)

	log.Infof("successfully fetched %d tire change times from %s until %s", len(tireChangeTimes), from, until)

	return newTireChangeTimesResponse(tireChangeTimes), nil
}

func (s *tireChangeTimesService) book(uuid string, contactInformation string) (*tireChangeBookingResponse, error) {
	log.Infof("trying to book tire change time with uuid: %s", uuid)
	tireChangeTime := s.repository.oneByUUID(uuid)

	if bookingErr := tireChangeTime.makeBooking(contactInformation); bookingErr != nil {
		return nil, bookingErr
	}

	tireChangeTime = s.repository.save(tireChangeTime)

	log.Infof("successfully booked tire change time with uuid: %s", uuid)
	return newTireChangeTimeResponse(tireChangeTime.UUID, tireChangeTime.Time), nil
}
