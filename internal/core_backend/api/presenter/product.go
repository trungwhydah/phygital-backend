package presenter

import (
	"time"

	"backend-service/internal/core_backend/entity"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ProductResponse data struct
type ProductResponse struct {
	ProductDetail      entity.Product             `json:"product"`
	OrganizationDetail OrganizationDetailResponse `json:"organization"`
}

type OrganizationDetailResponse struct {
	ID               primitive.ObjectID `json:"id"`
	Status           string             `json:"status"`
	CreatedAt        time.Time          `json:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at"`
	OrganizationName string             `json:"org_name"`
	NameTag          string             `json:"org_tag_name"`
	LogoURL          string             `json:"org_logo_url"`
	OwnerID          string             `json:"owner_id"`
}

type ListProductResponse struct {
	ProductList []ProductResponse `json:"product_list"`
}

type CreateProductResponse struct {
	ProductID string `json:"product_id"`
}

// presenterProduct struct
type PresenterProduct struct{}

// presenterProduct interface
type ConvertProduct interface {
	ResponseGetProductDetail(product *entity.Product, organization *entity.Organization) *ProductResponse
	ResponseAllProducts(products *[]entity.Product, organizations *[]entity.Organization) *ListProductResponse
	ResponseCreateProduct(productID *string) *CreateProductResponse
}

// NewPresenterProduct Constructs presenter
func NewPresenterProduct() ConvertProduct {
	return &PresenterProduct{}
}

// Return property data response
func (pp *PresenterProduct) ResponseGetProductDetail(product *entity.Product, organization *entity.Organization) *ProductResponse {
	product.ParseAttribute()
	response := &ProductResponse{
		ProductDetail: *product,
		OrganizationDetail: OrganizationDetailResponse{
			ID:               organization.ID,
			Status:           organization.Status,
			CreatedAt:        organization.CreatedAt,
			UpdatedAt:        organization.UpdatedAt,
			OrganizationName: organization.OrganizationName,
			NameTag:          organization.NameTag,
			LogoURL:          organization.LogoURL,
			OwnerID:          organization.OwnerID,
		},
	}

	return response
}

func (pp *PresenterProduct) ResponseAllProducts(products *[]entity.Product, organizations *[]entity.Organization) *ListProductResponse {
	var response ListProductResponse
	for i, product := range *products {
		response.ProductList = append(response.ProductList, ProductResponse{
			ProductDetail: product,
			OrganizationDetail: OrganizationDetailResponse{
				ID:               (*organizations)[i].ID,
				Status:           (*organizations)[i].Status,
				CreatedAt:        (*organizations)[i].CreatedAt,
				UpdatedAt:        (*organizations)[i].UpdatedAt,
				OrganizationName: (*organizations)[i].OrganizationName,
				NameTag:          (*organizations)[i].NameTag,
				LogoURL:          (*organizations)[i].LogoURL,
				OwnerID:          (*organizations)[i].OwnerID,
			},
		})
	}

	return &response
}

func (pp *PresenterProduct) ResponseCreateProduct(productID *string) *CreateProductResponse {
	return &CreateProductResponse{
		ProductID: *productID,
	}
}
