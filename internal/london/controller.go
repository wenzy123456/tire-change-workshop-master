package london

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const v1Path = "/api/v1"

type controller struct {
	service *tireChangeTimesService
}

func registerController(router *gin.Engine, service *tireChangeTimesService) {
	c := &controller{service: service}

	router.GET(v1Path+"/tire-change-times/available", c.getTireChangeTimes)
	router.PUT(v1Path+"/tire-change-times/:uuid/booking", c.putTireChangeBooking)
}

// getTireChangeTimes godoc
// @Summary List of available tire change times
// @Accept xml
// @Produce xml
// @Param from query string true "search available times from date" Format(date) default(2006-01-02)
// @Param until query string true "search available times until date" Format(date) default(2030-01-02)
// @Success 200 {object} tireChangeTimesResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /tire-change-times/available [get]
func (c *controller) getTireChangeTimes(ctx *gin.Context) {
	var query tireChangeTimesSearchQuery

	if err := ctx.ShouldBind(&query); err != nil {
		panic(validationError{err})
	}

	availableTimes, err := c.service.getAvailable(query.From, query.Until)

	if err != nil {
		panic(err)
	}

	ctx.XML(http.StatusOK, availableTimes)
}

// putTireChangeBooking godoc
// @Summary Book tire change time
// @Accept xml
// @Produce xml
// @Param uuid path string true "available tire change time UUID" minlength(36) maxlength(36)
// @Param body body tireChangeBookingRequest true "Request body"
// @Success 200 {object} tireChangeBookingResponse
// @Failure 400 {object} errorResponse
// @Failure 422 {object} errorResponse "The tire change time has already been booked by another contact"
// @Failure 500 {object} errorResponse
// @Router /tire-change-times/{uuid}/booking [put]
func (c *controller) putTireChangeBooking(ctx *gin.Context) {
	var uri tireChangeBookingURI
	var request tireChangeBookingRequest

	if err := ctx.ShouldBindUri(&uri); err != nil {
		panic(validationError{err})
	} else if err := ctx.ShouldBindXML(&request); err != nil {
		panic(validationError{err})
	}

	booking, err := c.service.book(uri.UUID, request.ContactInformation)

	if err != nil {
		panic(err)
	}

	ctx.XML(http.StatusOK, booking)
}
