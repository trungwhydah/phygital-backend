package presenter

import (
	"math/big"
	"sort"

	"backend-service/internal/core_backend/entity"
)

// ProductItemResponse data struct
type ProductItemDetailResponse struct {
	ID          string  `json:"id"`
	OwnerID     string  `json:"owner_id"`
	ProductID   string  `json:"product_id"`
	FarmName    string  `json:"farm_name"`
	Varietal    string  `json:"varietal"`
	Process     string  `json:"process"`
	ProductName string  `json:"product_name"`
	RatingScore float64 `json:"rating_score"`
	TotalLike   int     `json:"total_like"`
	ItemIndex   int     `json:"item_index"`
}

// ProductItemResponse data struct
type ProductItemDetailWithOwnerResponse struct {
	ID            string `json:"id"`
	ProductID     string `json:"product_id"`
	ProductName   string `json:"product_name"`
	ItemIndex     int    `json:"item_index"`
	TotalLike     int    `json:"total_like"`
	OwnerID       string `json:"owner_id"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	Organization  string `json:"organization"`
	Role          string `json:"role"`
	Status        string `json:"status"`
}

type AllProductItemResponse struct {
	ProductItems []ProductItemDetailWithOwnerResponse `json:"product_items_list"`
}

type StoryDetailResponse struct {
	MappingDetail      StoryMappingResponse      `json:"mapping_detail"`
	ProductDetail      entity.Product            `json:"product_detail"`
	ProductItemDetail  StoryProductItemResponse  `json:"product_item_detail"`
	OwnerDetail        StoryOwnerResponse        `json:"owner_detail"`
	TemplateDetail     StoryTemplateResponse     `json:"template_detail"`
	OrganizationDetail StoryOrganizationResponse `json:"organization_detail"`
	DigitalAssetDetail StoryDigitalAssetResponse `json:"digital_asset_detail"`
	HomepageDetail     StoryHomepageResponse     `json:"homepage_detail"`
}

type StoryMappingResponse struct {
	ExternalURL string `json:"external_url"`
	Claimable   bool   `json:"claimable"`
}

type StoryProductItemResponse struct {
	ID        string `json:"product_item_id"`
	TotalLike int    `json:"total_like"`
	ItemIndex int    `json:"item_index"`
}

type StoryOwnerResponse struct {
	ID       string `json:"owner_id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

type StoryTemplateResponse struct {
	ID        string                       `json:"template_id"`
	Name      string                       `json:"name"`
	Category  string                       `json:"category"`
	Languages []string                     `json:"languages"`
	Pages     []StoryTemplatePagesResponse `json:"pages"`
	Menu      []StoryTemplateMenuResponse  `json:"menu"`
}

type StoryTemplatePagesResponse struct {
	Name       string                 `json:"name"`
	Type       string                 `json:"type"`
	URL        string                 `json:"url"`
	PageID     string                 `json:"page_id"`
	Attributes map[string]interface{} `json:"attributes"`
}

type StoryTemplateMenuResponse struct {
	Title StoryTemplateMenuTitleResponse `json:"title"`
	URL   string                         `json:"url"`
}

type StoryTemplateMenuTitleResponse struct {
	VI string `json:"vi"`
	EN string `json:"en"`
}

type StoryHomepageResponse struct {
	ID         string                 `json:"id"`
	Name       string                 `json:"name"`
	URLLink    string                 `json:"url_link"`
	Type       string                 `json:"type"`
	Attributes map[string]interface{} `json:"attributes"`
}

type StoryOrganizationResponse struct {
	ID                  string `json:"org_id"`
	OrganizationName    string `json:"org_name"`
	OrganizationTagName string `json:"org_tag_name"`
	OrganizationLogoURL string `json:"org_logo_url"`
}

type StoryDigitalAssetResponse struct {
	ID              string `json:"digital_asset_id"`
	TokenID         int64  `json:"token_id"`
	Status          string `json:"status"`
	CollectionID    string `json:"collection_id"`
	CollectionName  string `json:"collection_name"`
	Chain           string `json:"chain"`
	ContractAddress string `json:"contract_address"`
	Standard        string `json:"standard"`
}

type GalleryProductItemsListResponse struct {
	GalleryItems []GalleryProductItemsResponse `json:"gallery_items"`
}

type GalleryProductItemsListResponseV2 struct {
	GalleryItems []StoryDetailResponse `json:"gallery_items"`
}

type GalleryProductItemsResponse struct {
	TagID              string                       `json:"tag_id"`
	ProductItemID      string                       `json:"product_item_id"`
	Product            entity.Product               `json:"product_detail"`
	TotalLike          int                          `json:"total_like"`
	Homepage           StoryHomepageResponse        `json:"homepage"`
	Craftsmen          StoryHomepageResponse        `json:"craftsmen"`
	Pages              []StoryTemplatePagesResponse `json:"pages"`
	DigitalAssetID     string                       `json:"digital_asset_id"`
	CollectionID       string                       `json:"collection_id"`
	TokenID            int64                        `json:"token_id"`
	DigitalAssetStatus string                       `json:"digital_asset_status"`
	CollectionName     string                       `json:"collection_name"`
	Chain              string                       `json:"chain"`
	ContractAddress    string                       `json:"contract_address"`
	Standard           string                       `json:"standard"`
}

// Test NFT
type ProductItemMetadataResponse struct {
	Description   string      `json:"description"`
	ExternalURL   string      `json:"external_url"`
	Image         string      `json:"image"`
	Name          string      `json:"name"`
	Attributes    []Attribute `json:"attributes"`
	Animation_url string      `json:"animation_url"`
}

type Attribute struct {
	TraitType   string      `json:"trait_type"`
	Value       interface{} `json:"value"`
	DisplayType string      `json:"display_type"`
}

type ByItemIndex []ProductItemDetailWithOwnerResponse

func (m ByItemIndex) Len() int           { return len(m) }
func (m ByItemIndex) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }
func (m ByItemIndex) Less(i, j int) bool { return m[i].ItemIndex < m[j].ItemIndex }

type ByProductAndItemIndex []ProductItemDetailWithOwnerResponse

func (m ByProductAndItemIndex) Len() int      { return len(m) }
func (m ByProductAndItemIndex) Swap(i, j int) { m[i], m[j] = m[j], m[i] }
func (m ByProductAndItemIndex) Less(i, j int) bool {
	return m[i].ProductName < m[j].ProductID || (m[i].ProductName == m[j].ProductName && m[i].ItemIndex < m[j].ItemIndex)
}

// presenterProductItem struct
type PresenterProductItem struct{}

// presenterProductItem interface
type ConvertProductItem interface {
	ResponseProductItemDetail(productItem *entity.ProductItem, product *entity.Product) *ProductItemDetailResponse
	ResponseGetAllProductItemsOfProduct(productItems *[]entity.ProductItem, user *[]entity.User, productName *string) *AllProductItemResponse
	ResponseGetAllProductItemsInOrg(productItems *[]entity.ProductItem, user *[]entity.User, products *[]entity.Product) *AllProductItemResponse
	ResponseGetStoryDetail(mapping *entity.Mapping, product *entity.Product, productItem *entity.ProductItem, owner *entity.User, template *entity.TemplateWebpages, homepage *entity.WebPage, organization *entity.Organization, da *entity.DigitalAsset, dac *entity.DigitalAssetCollection) *StoryDetailResponse
	ResponseGalleryProductItems(*string, []int, *[]entity.Mapping, *[]entity.Product, *[]entity.WebPage, *[]entity.WebPage, *[]entity.TemplateWebpages, *[]entity.DigitalAsset, *[]entity.DigitalAssetCollection) *GalleryProductItemsListResponse
	ResponseGalleryProductItemsV2(mappings *[]*entity.Mapping, products *[]*entity.Product, productItems *[]*entity.ProductItem, owners *[]*entity.User, templates *[]*entity.TemplateWebpages, organizations *[]*entity.Organization, das *[]*entity.DigitalAsset, dacs *[]*entity.DigitalAssetCollection) *GalleryProductItemsListResponseV2
	ResponseGetMetadata(*big.Int) *ProductItemMetadataResponse
}

// NewPresenterProductItem Constructs presenter
func NewPresenterProductItem() ConvertProductItem {
	return &PresenterProductItem{}
}

func (pp *PresenterProductItem) ResponseGetMetadata(tokenId *big.Int) *ProductItemMetadataResponse {
	metadata := &ProductItemMetadataResponse{
		Description: "Chiếc nón lá gắn liền với đời sống tinh thần của người dân, của người phụ nữ Việt. Trải dài mọi miền đất nước, hình ảnh nón lá luôn hiện diện, đó chính là nét đẹp, nét duyên, là sự bình dị mộc mạc của người phụ nữ Việt Nam nói chung  và của người phụ nữ Đà Nẵng nói riêng",
		ExternalURL: "https://kyokai.vn",
		Image:       "https://storage.googleapis.com/assets-fygito/images/user_id/1688617560-Chiếc Nón - Đàm Đông.png",
		Name:        "Chiếc Nón - Đàm Đông",
		Attributes: []Attribute{
			{
				TraitType: "Tác Giả",
				Value:     "Đàm Đông",
			},
			{
				TraitType: "Cuộc thi",
				Value:     "Đá Non Nước 2023",
			},
			{
				TraitType: "Chất Liệu",
				Value:     "Đá Marble Trắng",
			},
			{
				TraitType: "Chiều Dài",
				Value:     "65 cm",
			},
			{
				TraitType: "Chiều Rộng",
				Value:     "40 cm",
			},
			{
				TraitType: "Chiều Cao",
				Value:     "150 cm",
			},
			{
				TraitType: "Category",
				Value:     "Sculpture",
			},
			{
				TraitType: "Origin",
				Value:     "Viet Nam",
			},
		},
		Animation_url: "https://storage.googleapis.com/assets-fygito/images/user_id/1688612305-Chiếc Nón - Đàm Đông.glb",
	}
	if tokenId.Cmp(big.NewInt(2)) == 0 {
		metadata.Animation_url = "https://storage.googleapis.com/assets-fygito/images/user_id/1689675844-%C4%90a%CC%80%20Na%CC%86%CC%83ng%20Ti%CC%80nh%20Ngu%CC%9Bo%CC%9B%CC%80i%20-%20%C4%90o%CC%82%CC%83%20Va%CC%86n%20Hie%CC%A3%CC%82u.glb"
	}
	return metadata
}

// Return property data response
func (pp *PresenterProductItem) ResponseProductItemDetail(productItem *entity.ProductItem, product *entity.Product) *ProductItemDetailResponse {
	response := &ProductItemDetailResponse{
		ID:        productItem.ID.Hex(),
		OwnerID:   productItem.OwnerID,
		ProductID: productItem.ProductID.Hex(),
		// FarmName:    product.FarmName,
		// Varietal:    product.Varietal,
		// Process:     product.Process,
		ProductName: product.ProductName,
		RatingScore: product.RatingScore,
		TotalLike:   productItem.TotalLike,
		ItemIndex:   productItem.ItemIndex,
	}

	return response
}

func (pp *PresenterProductItem) ResponseGetAllProductItemsOfProduct(productItems *[]entity.ProductItem, users *[]entity.User, productName *string) *AllProductItemResponse {
	var response AllProductItemResponse
	for i := 0; i < len(*productItems); i++ {
		productItem := (*productItems)[i]
		user := (*users)[i]
		response.ProductItems = append(response.ProductItems, ProductItemDetailWithOwnerResponse{
			ID:            productItem.ID.Hex(),
			ProductID:     productItem.ProductID.Hex(),
			ProductName:   *productName,
			TotalLike:     productItem.TotalLike,
			ItemIndex:     productItem.ItemIndex,
			OwnerID:       productItem.OwnerID,
			Email:         user.Email,
			EmailVerified: user.EmailVerified,
			Name:          user.Name,
			Picture:       user.Picture,
			Organization:  user.Organization,
			Role:          user.Role,
			Status:        user.Status,
		})
	}
	sort.Sort(ByItemIndex(response.ProductItems))
	return &response
}

func (pp *PresenterProductItem) ResponseGetAllProductItemsInOrg(productItems *[]entity.ProductItem, users *[]entity.User, products *[]entity.Product) *AllProductItemResponse {
	var response AllProductItemResponse
	for i := 0; i < len(*productItems); i++ {
		productItem := (*productItems)[i]
		user := (*users)[i]
		response.ProductItems = append(response.ProductItems, ProductItemDetailWithOwnerResponse{
			ID:            productItem.ID.Hex(),
			ProductID:     productItem.ProductID.Hex(),
			ProductName:   (*products)[i].ProductName,
			TotalLike:     productItem.TotalLike,
			ItemIndex:     productItem.ItemIndex,
			OwnerID:       productItem.OwnerID,
			Email:         user.Email,
			EmailVerified: user.EmailVerified,
			Name:          user.Name,
			Picture:       user.Picture,
			Organization:  user.Organization,
			Role:          user.Role,
			Status:        user.Status,
		})
	}
	sort.Sort(ByProductAndItemIndex(response.ProductItems))
	return &response
}

func (pp *PresenterProductItem) ResponseGetStoryDetail(mapping *entity.Mapping, product *entity.Product, productItem *entity.ProductItem, owner *entity.User, template *entity.TemplateWebpages, homepage *entity.WebPage, organization *entity.Organization, da *entity.DigitalAsset, dac *entity.DigitalAssetCollection) *StoryDetailResponse {
	response := &StoryDetailResponse{}

	if mapping != nil {
		response.MappingDetail = StoryMappingResponse{
			ExternalURL: mapping.ExternalURL,
			Claimable:   mapping.Claimable,
		}
	}

	if product != nil {
		response.ProductDetail = *product
	}

	if productItem != nil {
		response.ProductItemDetail = StoryProductItemResponse{
			ID:        productItem.ID.Hex(),
			TotalLike: productItem.TotalLike,
			ItemIndex: productItem.ItemIndex,
		}
	}

	if owner != nil {
		response.OwnerDetail = StoryOwnerResponse{
			ID:       owner.ID,
			FullName: owner.Name,
			Email:    owner.Email,
		}
	}

	if template != nil {
		var pages []StoryTemplatePagesResponse
		var menus []StoryTemplateMenuResponse
		for _, page := range template.Pages {
			pages = append(pages, StoryTemplatePagesResponse{
				Name:       page.Name,
				Type:       page.Type,
				URL:        page.URLLink,
				PageID:     page.ID.Hex(),
				Attributes: page.Attributes,
			})
		}
		for _, menu := range template.Menu {
			menus = append(menus, StoryTemplateMenuResponse{
				Title: StoryTemplateMenuTitleResponse{
					VI: menu.Title.VI,
					EN: menu.Title.EN,
				},
				URL: menu.URLLink,
			})
		}
		response.TemplateDetail = StoryTemplateResponse{
			ID:        template.ID.Hex(),
			Name:      template.Name,
			Category:  template.Category,
			Languages: template.Languages,
			Pages:     pages,
			Menu:      menus,
		}
	}

	if homepage != nil {
		response.HomepageDetail = StoryHomepageResponse{
			ID:         homepage.ID.Hex(),
			Name:       homepage.Name,
			URLLink:    homepage.URLLink,
			Type:       homepage.Type,
			Attributes: homepage.Attributes,
		}
	}

	if organization != nil {
		response.OrganizationDetail = StoryOrganizationResponse{
			ID:                  organization.ID.Hex(),
			OrganizationName:    organization.OrganizationName,
			OrganizationTagName: organization.NameTag,
			OrganizationLogoURL: organization.LogoURL,
		}
	}

	if da != nil {
		response.DigitalAssetDetail = StoryDigitalAssetResponse{
			ID:              da.ID.Hex(),
			CollectionID:    da.CollectionID.Hex(),
			TokenID:         da.TokenID,
			Status:          da.Status,
			CollectionName:  dac.Name,
			Chain:           dac.Chain,
			Standard:        dac.Standard,
			ContractAddress: dac.ContractAddress,
		}
	}
	return response
}

func (pp *PresenterProductItem) ResponseGalleryProductItems(orgName *string, totalLike []int, mappings *[]entity.Mapping, products *[]entity.Product, homepages *[]entity.WebPage, craftsmen *[]entity.WebPage, templates *[]entity.TemplateWebpages, das *[]entity.DigitalAsset, dacs *[]entity.DigitalAssetCollection) *GalleryProductItemsListResponse {
	response := &GalleryProductItemsListResponse{}
	for i, mapping := range *mappings {
		var pages []StoryTemplatePagesResponse
		for _, page := range (*templates)[i].Pages {
			pages = append(pages, StoryTemplatePagesResponse{
				Type:   page.Type,
				URL:    page.URLLink,
				PageID: page.ID.Hex(),
			})
		}
		info := GalleryProductItemsResponse{
			TagID:         mapping.TagID,
			ProductItemID: mapping.ProductItemID.Hex(),
			Product:       (*products)[i],
			TotalLike:     totalLike[i],
			Homepage: StoryHomepageResponse{
				ID:         (*homepages)[i].ID.Hex(),
				Name:       (*homepages)[i].Name,
				URLLink:    (*homepages)[i].URLLink,
				Type:       (*homepages)[i].Type,
				Attributes: (*homepages)[i].Attributes,
			},
			Craftsmen: StoryHomepageResponse{
				ID:         (*craftsmen)[i].ID.Hex(),
				Name:       (*craftsmen)[i].Name,
				URLLink:    (*craftsmen)[i].URLLink,
				Type:       (*craftsmen)[i].Type,
				Attributes: (*craftsmen)[i].Attributes,
			},
			Pages: pages,
		}
		if !(*das)[i].ID.IsZero() {
			info.DigitalAssetID = (*das)[i].ID.Hex()
			info.DigitalAssetStatus = (*das)[i].Status
			info.TokenID = (*das)[i].TokenID
			info.CollectionID = (*das)[i].CollectionID.Hex()
			info.CollectionName = (*dacs)[i].Name
			info.Chain = (*dacs)[i].Chain
			info.Standard = (*dacs)[i].Standard
			info.ContractAddress = (*dacs)[i].ContractAddress
		}
		response.GalleryItems = append(response.GalleryItems, info)
	}
	return response
}

func (pp *PresenterProductItem) ResponseGalleryProductItemsV2(mappings *[]*entity.Mapping, products *[]*entity.Product, productItems *[]*entity.ProductItem, owners *[]*entity.User, templates *[]*entity.TemplateWebpages, organizations *[]*entity.Organization, das *[]*entity.DigitalAsset, dacs *[]*entity.DigitalAssetCollection) *GalleryProductItemsListResponseV2 {
	response := &GalleryProductItemsListResponseV2{}
	for i, mapping := range *mappings {
		product := (*products)[i]
		productItem := (*productItems)[i]
		owner := (*owners)[i]
		template := (*templates)[i]
		organization := (*organizations)[i]
		da := (*das)[i]
		dac := (*dacs)[i]
		info := StoryDetailResponse{}
		if mapping != nil {
			info.MappingDetail = StoryMappingResponse{
				ExternalURL: mapping.ExternalURL,
				Claimable:   mapping.Claimable,
			}
		}

		if product != nil {
			info.ProductDetail = *product
		}

		if productItem != nil {
			info.ProductItemDetail = StoryProductItemResponse{
				ID:        productItem.ID.Hex(),
				TotalLike: productItem.TotalLike,
				ItemIndex: productItem.ItemIndex,
			}
		}

		if owner != nil {
			info.OwnerDetail = StoryOwnerResponse{
				ID:       owner.ID,
				FullName: owner.Name,
				Email:    owner.Email,
			}
		}

		if template != nil {
			var pages []StoryTemplatePagesResponse
			var menus []StoryTemplateMenuResponse
			for _, page := range template.Pages {
				pages = append(pages, StoryTemplatePagesResponse{
					Name:       page.Name,
					Type:       page.Type,
					URL:        page.URLLink,
					PageID:     page.ID.Hex(),
					Attributes: page.Attributes,
				})
			}
			for _, menu := range template.Menu {
				menus = append(menus, StoryTemplateMenuResponse{
					Title: StoryTemplateMenuTitleResponse{
						VI: menu.Title.VI,
						EN: menu.Title.EN,
					},
					URL: menu.URLLink,
				})
			}
			info.TemplateDetail = StoryTemplateResponse{
				ID:        template.ID.Hex(),
				Name:      template.Name,
				Category:  template.Category,
				Languages: template.Languages,
				Pages:     pages,
				Menu:      menus,
			}
		}

		if organization != nil {
			info.OrganizationDetail = StoryOrganizationResponse{
				ID:                  organization.ID.Hex(),
				OrganizationName:    organization.OrganizationName,
				OrganizationTagName: organization.NameTag,
				OrganizationLogoURL: organization.LogoURL,
			}
		}

		if da != nil {
			info.DigitalAssetDetail = StoryDigitalAssetResponse{
				ID:              da.ID.Hex(),
				CollectionID:    da.CollectionID.Hex(),
				TokenID:         da.TokenID,
				Status:          da.Status,
				CollectionName:  dac.Name,
				Chain:           dac.Chain,
				Standard:        dac.Standard,
				ContractAddress: dac.ContractAddress,
			}
		}
		response.GalleryItems = append(response.GalleryItems, info)
	}
	return response
}
