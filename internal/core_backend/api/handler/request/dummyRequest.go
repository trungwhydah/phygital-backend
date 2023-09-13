package request

type DummyRequest struct {
	ID int `form:"id" validate:"required"`
}
