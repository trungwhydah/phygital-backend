package presenter

// SessionResponse data struct
type VerifySessionResponse struct {
	SesssionID     string `json:"session_id"`
	Message        string `json:"message"`
	IsValidSession bool   `json:"is_valid_session"`
}

// presenterSession struct
type PresenterSession struct{}

// presenterSession interface
type ConvertSession interface {
	ResponseVerifySession(sessionID, message *string) *VerifySessionResponse
}

// NewPresenterSession Constructs presenter
func NewPresenterSession() ConvertSession {
	return &PresenterSession{}
}

// Return property data response
func (pp *PresenterSession) ResponseVerifySession(sessionID, message *string) *VerifySessionResponse {
	var response = &VerifySessionResponse{
		SesssionID:     *sessionID,
		Message:        *message,
		IsValidSession: *message == "",
	}

	return response
}
