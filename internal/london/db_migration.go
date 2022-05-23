package london

import (
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"gopkg.in/gormigrate.v1"
	"time"
)

var initial = &gormigrate.Migration{
	ID: "201608301400",

	Migrate: func(db *gorm.DB) error {
		// it's a good practise to copy the struct inside the function,
		// so side effects are prevented if the original struct changes during the time
		type tireChangeTimeEntityVersion1 struct {
			ID   uint   `gorm:"primary_key"`
			UUID string `gorm:"size:36;unique_index; not null"`

			Time time.Time

			Available bool

			BookedByContact string

			CreatedAt time.Time
			UpdatedAt time.Time
		}

		err := db.Table(tireChangeTimeEntity{}.TableName()).CreateTable(&tireChangeTimeEntityVersion1{}).Error

		if err == nil {
			nextTime := time.Now().AddDate(0, 0, -7)

			for i := 0; i < 500; i++ {
				nextTime = calculateTireChangeTime(nextTime.Add(time.Hour * time.Duration(1)))
				err = db.Create(newTireChangeTimeEntity(nextTime, calculateAvailability(nextTime))).Error
			}
		}

		if err == nil {
			log.Info("Migrated 201608301400")
		}

		return err
	},

	Rollback: func(tx *gorm.DB) error {
		return tx.DropTable("catalog").Error
	},
}

func calculateTireChangeTime(nextTime time.Time) time.Time {
	for notWorkDay(nextTime) {
		nextTime = nextTime.Add(time.Hour * time.Duration(1))
	}

	return time.Date(
		nextTime.Year(),
		nextTime.Month(),
		nextTime.Day(),
		nextTime.Hour(),
		0,
		0,
		0,
		nextTime.Location(),
	)
}

func notWorkDay(tireChangeTime time.Time) bool {
	return tireChangeTime.Hour() < 8 ||
		tireChangeTime.Hour() > 16 ||
		tireChangeTime.Weekday() == time.Saturday ||
		tireChangeTime.Weekday() == time.Sunday
}

func calculateAvailability(tireChangeTime time.Time) bool {
	dateSum := tireChangeTime.Year() + int(tireChangeTime.Month()) + tireChangeTime.Day() + tireChangeTime.Hour()

	return dateSum%5 > 0
}
