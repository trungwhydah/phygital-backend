package request

type CreateProductItemRequest struct {
	ProductID string `form:"product_id" validate:"required"`
}

type CreateMultipleProductItemsRequest struct {
	NumItems  int    `form:"num_item" validate:"required"`
	ProductID string `form:"product_id" validate:"required"`
}

type ProductItemInteractionRequest struct {
	ProductItemID string `form:"product_item_id,omitempty" validate:"required"`
}

type SetOwnerRequest struct {
	ProductItemID string `form:"product_item_id,omitempty" validate:"required"`
	Token         string `form:"token" validate:"required"`
	OwnerID       string
}

type ProductItemLikeRequest struct {
	ProductItemID string `form:"product_item_id" validate:"required"`
	// UserID        string `form:"user_id" validate:"required"`
}
