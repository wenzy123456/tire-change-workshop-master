package london

import (
	"bytes"
	"encoding/xml"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const rfc3339DateFormat = "2006-01-02"

func TestGetAvailableTireChangeTimes(t *testing.T) {
	router := Init(true)

	t.Run("successfully get all available for today and tomorrow in correct order", func(t *testing.T) {
		today := time.Now().Format(rfc3339DateFormat)
		nextWeek := time.Now().AddDate(0, 0, 7).Format(rfc3339DateFormat)
		reqURL := fmt.Sprintf(v1Path+"/tire-change-times/available?from=%s&until=%s", today, nextWeek)

		requestWriter := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, reqURL, nil)
		router.ServeHTTP(requestWriter, req)

		result := &tireChangeTimesResponse{}
		unMarshal(t, requestWriter.Body.Bytes(), result)

		assert.Equal(t, http.StatusOK, requestWriter.Code)
		verifyTireChangeTimesResponse(t, result)
	})

	t.Run("fail with invalid date format", func(t *testing.T) {
		today := time.Now().Format(rfc3339DateFormat)
		reqURL := fmt.Sprintf(v1Path+"/tire-change-times/available?from=%s&until=INVALID", today)

		requestWriter := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, reqURL, nil)
		router.ServeHTTP(requestWriter, req)

		result := &errorResponse{}
		unMarshal(t, requestWriter.Body.Bytes(), result)

		assert.Equal(t, http.StatusBadRequest, requestWriter.Code)
		assert.Equal(t, http.StatusBadRequest, result.StatusCode)
		assert.NotEmpty(t, result.Error)
	})

	t.Run("fail with invalid date period", func(t *testing.T) {
		today := time.Now().Format(rfc3339DateFormat)
		yesterday := time.Now().AddDate(0, 0, -1).Format(rfc3339DateFormat)
		reqURL := fmt.Sprintf(v1Path+"/tire-change-times/available?from=%s&until=%s", today, yesterday)

		requestWriter := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, reqURL, nil)
		router.ServeHTTP(requestWriter, req)

		result := &errorResponse{}
		unMarshal(t, requestWriter.Body.Bytes(), result)

		assert.Equal(t, http.StatusBadRequest, requestWriter.Code)
		assert.Equal(t, http.StatusBadRequest, result.StatusCode)
		assert.NotEmpty(t, result.Error)
	})
}

func TestTireChangeTimeBooking(t *testing.T) {
	router := Init(true)

	t.Run("successfully book available tire change time", func(t *testing.T) {
		availableTireChangeTime := newTireChangeTimeEntity(time.Now(), true)
		must(t, db.Create(availableTireChangeTime).Error)

		reqURL := fmt.Sprintf(v1Path+"/tire-change-times/%s/booking", availableTireChangeTime.UUID)
		request := &tireChangeBookingRequest{ContactInformation: "TEST"}

		requestWriter := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPut, reqURL, marshal(t, request))
		router.ServeHTTP(requestWriter, req)

		result := &tireChangeBookingResponse{}
		unMarshal(t, requestWriter.Body.Bytes(), result)

		assert.Equal(t, http.StatusOK, requestWriter.Code)
		assert.Equal(t, availableTireChangeTime.UUID, result.UUID)
		assert.Equal(t, availableTireChangeTime.Time.UTC(), result.Time)
		assert.Equal(t, request.ContactInformation, getTireChangeTime(t, availableTireChangeTime.UUID).BookedByContact)
		assert.False(t, getTireChangeTime(t, availableTireChangeTime.UUID).Available)
	})

	t.Run("successfully update already booked change time for same contact", func(t *testing.T) {
		contactInformation := "TEST"
		bookedTireChangeTime := newTireChangeTimeEntity(time.Now(), false)
		bookedTireChangeTime.BookedByContact = contactInformation
		must(t, db.Create(bookedTireChangeTime).Error)

		reqURL := fmt.Sprintf(v1Path+"/tire-change-times/%s/booking", bookedTireChangeTime.UUID)
		request := &tireChangeBookingRequest{ContactInformation: contactInformation}

		requestWriter := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPut, reqURL, marshal(t, request))
		router.ServeHTTP(requestWriter, req)

		assert.Equal(t, http.StatusOK, requestWriter.Code)
	})

	t.Run("fail to book unavailable tire change time", func(t *testing.T) {
		unAvailableTireChangeTime := newTireChangeTimeEntity(time.Now(), false)
		unAvailableTireChangeTime.BookedByContact = "some guy"
		must(t, db.Create(unAvailableTireChangeTime).Error)

		reqURL := fmt.Sprintf(v1Path+"/tire-change-times/%s/booking", unAvailableTireChangeTime.UUID)
		request := &tireChangeBookingRequest{ContactInformation: "another guy"}

		requestWriter := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPut, reqURL, marshal(t, request))
		router.ServeHTTP(requestWriter, req)

		assert.Equal(t, http.StatusUnprocessableEntity, requestWriter.Code)
	})

	t.Run("fail to book unknown tire change time", func(t *testing.T) {
		reqURL := fmt.Sprintf(v1Path+"/tire-change-times/%s/booking", uuid.NewV4().String())
		request := &tireChangeBookingRequest{ContactInformation: "TEST"}

		requestWriter := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPut, reqURL, marshal(t, request))
		router.ServeHTTP(requestWriter, req)

		assert.Equal(t, http.StatusUnprocessableEntity, requestWriter.Code)
	})

	t.Run("fail to book with invalid request uri", func(t *testing.T) {
		reqURL := fmt.Sprintf(v1Path+"/tire-change-times/%s/booking", "INVALID")

		requestWriter := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPut, reqURL, marshal(t, &tireChangeBookingRequest{}))
		router.ServeHTTP(requestWriter, req)

		assert.Equal(t, http.StatusBadRequest, requestWriter.Code)
	})

	t.Run("fail to book with invalid request", func(t *testing.T) {
		availableTireChangeTime := newTireChangeTimeEntity(time.Now(), true)
		must(t, db.Create(availableTireChangeTime).Error)

		reqURL := fmt.Sprintf(v1Path+"/tire-change-times/%s/booking", availableTireChangeTime.UUID)
		invalidRequest := &tireChangeBookingRequest{}

		requestWriter := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPut, reqURL, marshal(t, invalidRequest))
		router.ServeHTTP(requestWriter, req)

		result := &errorResponse{}
		unMarshal(t, requestWriter.Body.Bytes(), result)

		assert.Equal(t, http.StatusBadRequest, requestWriter.Code)
		assert.NotEmpty(t, result.Error)
		assert.Equal(t, http.StatusBadRequest, result.StatusCode)
	})
}

func getTireChangeTime(t *testing.T, uuid string) *tireChangeTimeEntity {
	var result tireChangeTimeEntity

	if err := db.Model(tireChangeTimeEntity{}).Where("uuid = ?", uuid).Find(&result).Error; err != nil {
		t.Fatalf("failed to fetch tire change time, error: %v", err)
	}

	return &result
}

func marshal(t *testing.T, value interface{}) io.Reader {
	valueXML, e := xml.Marshal(value)
	must(t, e)

	return bytes.NewBuffer(valueXML)
}

func unMarshal(t *testing.T, rawXML []byte, result interface{}) {
	if err := xml.Unmarshal(rawXML, result); err != nil {
		t.Fatalf("failed to unmarshal XML, error: %v", err)
	}
}

func verifyTireChangeTimesResponse(t *testing.T, result *tireChangeTimesResponse) {
	assert.NotEmpty(t, result)
	firstTime := result.AvailableTimes[0].Time

	for _, availableTime := range result.AvailableTimes {
		tireChangeTimeEntity := getTireChangeTime(t, availableTime.UUID)
		assert.NotNil(t, tireChangeTimeEntity)
		assert.Equal(t, tireChangeTimeEntity.Time.UTC(), availableTime.Time)
		assert.True(
			t,
			availableTime.Time.Equal(firstTime) || availableTime.Time.After(firstTime),
			"response items should be in ascending order",
		)
	}
}

func must(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("failed to run test task, error: %v", err)
	}
}
