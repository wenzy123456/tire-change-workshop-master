package manchester

import "time"

type tireChangeTimesSearchQuery struct {
	Amount uint      `form:"amount"`
	Page   uint      `form:"page" binding:"required_with=Amount"`
	From   time.Time `form:"from" time_format:"2006-01-02"`
}

func (q *tireChangeTimesSearchQuery) offset() uint {
	return q.Page * q.Amount
}

func (q *tireChangeTimesSearchQuery) isPaginated() bool {
	return q.Amount > 0
}

type tireChangeBookingURI struct {
	ID uint `uri:"id" binding:"required"`
}

type tireChangeBookingRequest struct {
	ContactInformation string `json:"contactInformation" binding:"required,min=1"`
}
