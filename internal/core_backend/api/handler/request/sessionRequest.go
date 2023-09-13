package request

type SessionInteractRequest struct {
	SessionID string `json:"session_id" validate:"required"`
}
