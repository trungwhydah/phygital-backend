package presenter

import (
	"backend-service/internal/core_backend/entity"
	"sort"
	"sync"
)

// MappingResponse data struct
type GetAllMappingResponse struct {
	Mappings []Mapping `json:"mapping"`
}

type ByTagID []Mapping

func (m ByTagID) Len() int           { return len(m) }
func (m ByTagID) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }
func (m ByTagID) Less(i, j int) bool { return m[i].TagID < m[j].TagID }

type Mapping struct {
	TagID              string `json:"tag_id"`
	ProductItemID      string `json:"product_item_id"`
	OrganizationID     string `json:"org_id"`
	ExternalURL        string `json:"external_url"`
	Claimable          bool   `json:"claimable"`
	OwnerID            string `json:"owner_id"`
	OwnerEmail         string `json:"owner_email"`
	OwnerName          string `json:"owner_name"`
	DigitalAssetID     string `json:"digital_asset_id"`
	CollectionID       string `json:"collection_id"`
	TokenID            int64  `json:"token_id"`
	DigitalAssetStatus string `json:"digital_asset_status"`
}

type ProductsAbleToMapping struct {
	TemplateID   string `json:"template_id"`
	TemplateName string `json:"template_name"`
}

type MultipleMappingResponse struct {
	MappingResults []MappingResult `json:"mapping_results"`
}

type MappingResult struct {
	TagID     string `json:"tag_id"`
	IsSuccess bool   `json:"is_success"`
}

// presenterMapping struct
type PresenterMapping struct{}

// presenterMapping interface
type ConvertMapping interface {
	ResponseGetAllMapping(mappings *[]entity.Mapping, owners *[]entity.User, assets *[]entity.DigitalAsset) *GetAllMappingResponse
	MultipleMappingWithSingleProductResponse(*map[string]bool) *MultipleMappingResponse
}

// NewPresenterMapping Constructs presenter
func NewPresenterMapping() ConvertMapping {
	return &PresenterMapping{}
}

// ResponseGetAllMapping
func (pp *PresenterMapping) ResponseGetAllMapping(mappings *[]entity.Mapping, owners *[]entity.User, assets *[]entity.DigitalAsset) *GetAllMappingResponse {
	var (
		response GetAllMappingResponse
		wg       sync.WaitGroup
		mutex    sync.Mutex
	)

	wg.Add(len(*mappings))
	for _, mapping := range *mappings {
		go func(mapping entity.Mapping, listOwner *[]entity.User, listAssets *[]entity.DigitalAsset) {
			resMap := Mapping{
				TagID:          mapping.TagID,
				ProductItemID:  mapping.ProductItemID.Hex(),
				OrganizationID: mapping.OrganizationID.Hex(),
				Claimable:      mapping.Claimable,
				ExternalURL:    mapping.ExternalURL,
			}

			for _, o := range *listOwner {
				if mapping.OwnerID == o.ID {
					resMap.OwnerID = o.ID
					resMap.OwnerEmail = o.Email
					resMap.OwnerName = o.Name
				}
			}
			if !mapping.DigitalAssetID.IsZero() {
				for _, a := range *listAssets {
					if mapping.DigitalAssetID == a.ID {
						resMap.DigitalAssetID = a.ID.Hex()
						resMap.CollectionID = a.CollectionID.Hex()
						resMap.TokenID = a.TokenID
						resMap.DigitalAssetStatus = a.Status
					}
				}
			}
			// Acquire lock before modifying the response variable
			mutex.Lock()
			response.Mappings = append(response.Mappings, resMap)
			mutex.Unlock()
			defer wg.Done()
		}(mapping, owners, assets)
	}
	wg.Wait()
	sort.Sort(ByTagID(response.Mappings))
	return &response
}

func (pp *PresenterMapping) MultipleMappingWithSingleProductResponse(result *map[string]bool) *MultipleMappingResponse {
	var response MultipleMappingResponse
	for tID, iS := range *result {
		response.MappingResults = append(response.MappingResults, MappingResult{
			TagID:     tID,
			IsSuccess: iS,
		})
	}

	return &response
}
