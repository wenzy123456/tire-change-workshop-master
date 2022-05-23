package london

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

var zeroTireChangeTimeEntity = &tireChangeTimeEntity{}

type tireChangeTimeEntity struct {
	ID   uint   `gorm:"primary_key"`
	UUID string `gorm:"size:36;unique_index; not null"`

	Time time.Time

	Available bool

	BookedByContact string

	CreatedAt time.Time
	UpdatedAt time.Time
}

func newTireChangeTimeEntity(changeTime time.Time, available bool) *tireChangeTimeEntity {
	return &tireChangeTimeEntity{
		UUID:      uuid.NewV4().String(),
		Time:      changeTime,
		Available: available,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (e *tireChangeTimeEntity) makeBooking(contactInformation string) error {
	if e == zeroTireChangeTimeEntity || (!e.Available && e.BookedByContact != contactInformation) {
		return newUnAvailableBookingError(e)
	}

	e.Available = false
	e.UpdatedAt = time.Now()
	e.BookedByContact = contactInformation

	return nil
}

func (e tireChangeTimeEntity) TableName() string {
	return "tire_change_time"
}
