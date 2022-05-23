package manchester

import (
	"github.com/jinzhu/gorm"
)

type tireChangeTimeRepository struct {
	db *gorm.DB
}

func newTireChangeTimeRepository(db *gorm.DB) *tireChangeTimeRepository {
	return &tireChangeTimeRepository{db: db}
}

func (r *tireChangeTimeRepository) allBySearchQuery(searchQuery *tireChangeTimesSearchQuery) []*tireChangeTimeEntity {
	results := make([]*tireChangeTimeEntity, 0)

	query := r.db.Model(&tireChangeTimeEntity{}).Order("time ASC")

	if searchQuery.isPaginated() {
		query = query.Offset(searchQuery.offset()).Limit(searchQuery.Amount)
	}

	if !searchQuery.From.IsZero() {
		query = query.Where("time >= ?", searchQuery.From)
	}

	if err := query.Find(&results).Error; err != nil {
		panic(err)
	}

	return results
}

func (r *tireChangeTimeRepository) availableByID(id uint) *tireChangeTimeEntity {
	var result tireChangeTimeEntity

	query := r.db.Model(&tireChangeTimeEntity{}).Where("id = ?", id)

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
