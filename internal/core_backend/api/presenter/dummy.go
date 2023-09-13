package presenter

import (
	"time"

	"backend-service/internal/core_backend/entity"
)

// DummyResponse data struct
type DummyResponse struct {
	ID           int       `json:"-"`
	Message      string    `json:"message"`
	ReceivedTime time.Time `json:"received"`
}

// presenterDummy struct
type PresenterDummy struct{}

// presenterDummy interface
type ConvertDummy interface {
	ResponseDummy(dummy *entity.Dummy) *DummyResponse
}

// NewPresenterDummy Constructs presenter
func NewPresenterDummy() ConvertDummy {
	return &PresenterDummy{}
}

// Return property data response
func (pp *PresenterDummy) ResponseDummy(dummy *entity.Dummy) *DummyResponse {
	response := &DummyResponse{
		ReceivedTime: time.Now(),
	}

	if dummy == nil {
		response.ID = 0
		response.Message = "Hello world"
	} else {
		response.ID = dummy.ID
		response.Message = dummy.Message
	}

	return response
}
