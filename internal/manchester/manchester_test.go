package manchester

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const rfc3339DateFormat = "2006-01-02"

func TestGetTireChangeTimes(t *testing.T) {
	router := Init(true)

	t.Run("successfully get all in correct order", func(t *testing.T) {
		reqURL := v2Path + "/tire-change-times"

		requestWriter := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, reqURL, nil)
		router.ServeHTTP(requestWriter, req)

		result := &tireChangeTimesResponse{}
		unMarshal(t, requestWriter.Body.Bytes(), result)

		assert.Equal(t, http.StatusOK, requestWriter.Code)
		verifyTireChangeTimesResponse(t, *result)
	})

	t.Run("successfully get subset in correct order", func(t *testing.T) {
		reqURL := v2Path + "/tire-change-times?amount=101&page=14"

		requestWriter := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, reqURL, nil)
		router.ServeHTTP(requestWriter, req)

		result := &tireChangeTimesResponse{}
		unMarshal(t, requestWriter.Body.Bytes(), result)

		assert.Equal(t, http.StatusOK, requestWriter.Code)
		assert.Len(t, *result, 86) // Db contains 1500 rows, 101 * 14 = 1414 and subset count 1500 - 1414 = 86
		verifyTireChangeTimesResponse(t, *result)
	})

	t.Run("successfully get all from today in correct order", func(t *testing.T) {
		today := time.Now().Format(rfc3339DateFormat)
		reqURL := fmt.Sprintf(v2Path+"/tire-change-times?from=%s", today)

		requestWriter := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, reqURL, nil)
		router.ServeHTTP(requestWriter, req)

		result := &tireChangeTimesResponse{}
		unMarshal(t, requestWriter.Body.Bytes(), result)

		assert.Equal(t, http.StatusOK, requestWriter.Code)
		// Db contains 1500 rows, created since week ago, week contains 45 tire change times, therefore  1500 - 45 = 155
		assert.Len(t, *result, 1455)
		verifyTireChangeTimesResponse(t, *result)
	})

	t.Run("fail to get all from invalid date", func(t *testing.T) {
		reqURL := v2Path + "/tire-change-times?from=INVALID"

		requestWriter := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, reqURL, nil)
		router.ServeHTTP(requestWriter, req)

		result := &errorResponse{}
		unMarshal(t, requestWriter.Body.Bytes(), result)

		assert.Equal(t, http.StatusBadRequest, requestWriter.Code)
		assert.Equal(t, validationErrorCode, result.Code)
		assert.NotEmpty(t, result.Message)
	})

	t.Run("fail to get subset with invalid offset", func(t *testing.T) {
		reqURL := v2Path + "/tire-change-times?amount=1&page=0"

		requestWriter := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, reqURL, nil)
		router.ServeHTTP(requestWriter, req)

		result := &errorResponse{}
		unMarshal(t, requestWriter.Body.Bytes(), result)

		assert.Equal(t, http.StatusBadRequest, requestWriter.Code)
		assert.Equal(t, validationErrorCode, result.Code)
		assert.NotEmpty(t, result.Message)
	})
}

func TestTireChangeTimeBooking(t *testing.T) {
	router := Init(true)

	t.Run("successfully book available tire change time", func(t *testing.T) {
		availableTireChangeTime := newTireChangeTimeEntity(time.Now(), true)
		must(t, db.Create(availableTireChangeTime).Error)

		reqURL := fmt.Sprintf(v2Path+"/tire-change-times/%d/booking", availableTireChangeTime.ID)
		request := &tireChangeBookingRequest{ContactInformation: "TEST"}

		requestWriter := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, reqURL, marshal(t, request))
		router.ServeHTTP(requestWriter, req)

		result := &tireChangeTimeBookingResponse{}
		unMarshal(t, requestWriter.Body.Bytes(), result)

		assert.Equal(t, http.StatusOK, requestWriter.Code)
		assert.Equal(t, availableTireChangeTime.ID, result.ID)
		assert.Equal(t, availableTireChangeTime.Time.UTC(), result.Time)
		assert.False(t, getTireChangeTime(t, availableTireChangeTime.ID).Available)
	})

	t.Run("fail to book unavailable tire change time", func(t *testing.T) {
		availableTireChangeTime := newTireChangeTimeEntity(time.Now(), false)
		must(t, db.Create(availableTireChangeTime).Error)

		reqURL := fmt.Sprintf(v2Path+"/tire-change-times/%d/booking", availableTireChangeTime.ID)
		request := &tireChangeBookingRequest{ContactInformation: "TEST"}

		requestWriter := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, reqURL, marshal(t, request))
		router.ServeHTTP(requestWriter, req)

		result := &errorResponse{}
		unMarshal(t, requestWriter.Body.Bytes(), result)

		assert.Equal(t, http.StatusUnprocessableEntity, requestWriter.Code)
		assert.Equal(t, unAvailableTimeErrorCode, result.Code)
		assert.NotEmpty(t, result.Message)
	})

	t.Run("fail to book unknown tire change time", func(t *testing.T) {
		reqURL := fmt.Sprintf(v2Path+"/tire-change-times/%d/booking", 34534523423)
		request := &tireChangeBookingRequest{ContactInformation: "TEST"}

		requestWriter := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, reqURL, marshal(t, request))
		router.ServeHTTP(requestWriter, req)

		result := &errorResponse{}
		unMarshal(t, requestWriter.Body.Bytes(), result)

		assert.Equal(t, http.StatusUnprocessableEntity, requestWriter.Code)
		assert.Equal(t, unAvailableTimeErrorCode, result.Code)
		assert.NotEmpty(t, result.Message)
	})

	t.Run("fail to book with invalid request uri", func(t *testing.T) {
		reqURL := fmt.Sprintf(v2Path+"/tire-change-times/%s/booking", "INVALID")

		requestWriter := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, reqURL, marshal(t, &tireChangeBookingRequest{}))
		router.ServeHTTP(requestWriter, req)

		result := &errorResponse{}
		unMarshal(t, requestWriter.Body.Bytes(), result)

		assert.Equal(t, http.StatusBadRequest, requestWriter.Code)
		assert.Equal(t, validationErrorCode, result.Code)
		assert.NotEmpty(t, result.Message)
	})

	t.Run("fail to book with invalid request", func(t *testing.T) {
		availableTireChangeTime := newTireChangeTimeEntity(time.Now(), true)
		must(t, db.Create(availableTireChangeTime).Error)

		reqURL := fmt.Sprintf(v2Path+"/tire-change-times/%d/booking", availableTireChangeTime.ID)
		invalidRequest := &tireChangeBookingRequest{}

		requestWriter := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, reqURL, marshal(t, invalidRequest))
		router.ServeHTTP(requestWriter, req)

		result := &errorResponse{}
		unMarshal(t, requestWriter.Body.Bytes(), result)

		assert.Equal(t, http.StatusBadRequest, requestWriter.Code)
		assert.Equal(t, validationErrorCode, result.Code)
		assert.NotEmpty(t, result.Message)
	})
}

func getTireChangeTime(t *testing.T, id uint) *tireChangeTimeEntity {
	var result tireChangeTimeEntity

	if err := db.Model(tireChangeTimeEntity{}).Where("id = ?", id).Find(&result).Error; err != nil {
		t.Fatalf("failed to fetch tire change time, error: %v", err)
	}

	return &result
}

func marshal(t *testing.T, value interface{}) io.Reader {
	valueJSON, e := json.Marshal(value)
	must(t, e)

	return bytes.NewBuffer(valueJSON)
}

func unMarshal(t *testing.T, rawJSON []byte, result interface{}) {
	if err := json.Unmarshal(rawJSON, result); err != nil {
		t.Fatalf("failed to unmarshal JSON, error: %v", err)
	}
}

func verifyTireChangeTimesResponse(t *testing.T, result tireChangeTimesResponse) {
	assert.NotEmpty(t, result)
	firstTime := result[0].Time

	for _, tireChangeTime := range result {
		tireChangeTimeEntity := getTireChangeTime(t, tireChangeTime.ID)
		assert.NotNil(t, tireChangeTimeEntity)
		assert.Equal(t, tireChangeTimeEntity.Time.UTC(), tireChangeTime.Time)
		assert.Equal(t, tireChangeTimeEntity.Available, tireChangeTime.Available)
		assert.True(
			t,
			tireChangeTime.Time.Equal(firstTime) || tireChangeTime.Time.After(firstTime),
			"response items should be in ascending order",
		)
	}
}

func must(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("failed to run test task, error: %v", err)
	}
}
