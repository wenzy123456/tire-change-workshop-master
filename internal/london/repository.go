package london

import (
	"github.com/jinzhu/gorm"
	"time"
)

type tireChangeTimeRepository struct {
	db *gorm.DB
}

func newTireChangeTimeRepository(db *gorm.DB) *tireChangeTimeRepository {
	return &tireChangeTimeRepository{db: db}
}

func (r *tireChangeTimeRepository) availableByTimeRange(from time.Time, until time.Time) []*tireChangeTimeEntity {
	results := make([]*tireChangeTimeEntity, 0)

	query := r.db.Model(&tireChangeTimeEntity{}).
		Where("available = ?", true).
		Where("time >= ?", from).
		Where("time <= ?", until).
		Order("time ASC")

	if err := query.Find(&results).Error; err != nil {
		panic(err)
	}

	return results
}

func (r *tireChangeTimeRepository) oneByUUID(uuid string) *tireChangeTimeEntity {
	var result tireChangeTimeEntity

	query := r.db.Model(&tireChangeTimeEntity{}).Where("uuid = ?", uuid)

	if err := query.Find(&result).Error; gorm.IsRecordNotFoundError(err) {
		return zeroTireChangeTimeEntity
	} else if err != nil {
		panic(err)
	}

	return &result
}

func (r *tireChangeTimeRepository) save(entity *tireChangeTimeEntity) *tireChangeTimeEntity {
	if err := r.db.Save(entity).Error; err != nil {
		panic(err)
	}

	return entity
}
