package presenter

import (
	"backend-service/internal/core_backend/entity"
)

// AuthorResponse data struct
type AuthorDetailResponse struct {
	Author      *entity.Author    `json:"author"`
	ProductList *[]entity.Product `json:"product_list"`
}

// presenterAuthor struct
type PresenterAuthor struct{}

// presenterAuthor interface
type ConvertAuthor interface {
	ResponseAuthorDetail(author *entity.Author, productList *[]entity.Product) *AuthorDetailResponse
}

// NewPresenterAuthor Constructs presenter
func NewPresenterAuthor() ConvertAuthor {
	return &PresenterAuthor{}
}

// Return property data response
func (pp *PresenterAuthor) ResponseAuthorDetail(author *entity.Author, productList *[]entity.Product) *AuthorDetailResponse {
	return &AuthorDetailResponse{
		Author:      author,
		ProductList: productList,
	}
}
