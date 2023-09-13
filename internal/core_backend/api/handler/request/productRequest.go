package request

type InteractProductDetailRequest struct {
	ProductID string `json:"product_id" validate:"required"`
}
