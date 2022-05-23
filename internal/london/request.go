package london

import "time"

type tireChangeTimesSearchQuery struct {
	From  time.Time `form:"from" time_format:"2006-01-02" binding:"required"`
	Until time.Time `form:"until" time_format:"2006-01-02" binding:"required"`
}

type tireChangeBookingURI struct {
	UUID string `uri:"uuid" binding:"required,max=36,min=36"`
}

type tireChangeBookingRequest struct {
	ContactInformation string `xml:"contactInformation" binding:"required,min=1"`
}
