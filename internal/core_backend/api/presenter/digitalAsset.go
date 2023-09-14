package presenter

import (
	"backend-service/internal/core_backend/entity"
	"fmt"
	"strconv"
	"strings"
)

// DigitalAssetResponse data struct
type DigitalAssetResponse struct {
	ID              string `json:"id"`
	Status          string `json:"status"`
	CollectionID    string `json:"collection_id"`
	TokenID         int64  `json:"token_id"`
	Chain           string `json:"chain"`
	ChainID         int    `json:"chain_id"`
	CollectionName  string `json:"collection_name"`
	ContractAddress string `json:"contract_address"`
	Standard        string `json:"standard"`
	OwnerEmail      string `json:"owner_email"`
	OwnerName       string `json:"owner_full_name"`
	Metadata        any    `json:"metadata"`
}

type ListDigitalAssetsResponse struct {
	DigitalAssets []DigitalAssetResponse `json:"digital_assets"`
}

type MetadataResponse struct {
	Description  string              `json:"description"`
	ExternalURL  string              `json:"external_url"`
	Image        string              `json:"image"`
	Name         string              `json:"name"`
	Attributes   []MetadataAttribute `json:"attributes"`
	AnimationURL string              `json:"animation_url"`
}

type MetadataAttribute struct {
	TraitType   string `json:"trait_type"`
	Value       string `json:"value"`
	DisplayType string `json:"display_type"`
}

type WalletAddressResponse struct {
	WalletAddress string `json:"wallet_address"`
}

type CreateWalletResponse struct {
	Result WalletAddressResponse `json:"result"`
}

type UserWalletResponse struct {
	Email  string                `json:"email"`
	Wallet WalletAddressResponse `json:"wallet"`
}

type ListWalletsResponse struct {
	ListWallet []UserWalletResponse `json:"list_wallets"`
}

type CreateWalletsBulkResponse struct {
	Result ListWalletsResponse `json:"result"`
}

// presenterDigitalAsset struct
type PresenterDigitalAsset struct{}

// presenterDigitalAsset interface
type ConvertDigitalAsset interface {
	ResponseDigitalAssets(digitalAsset *[]entity.DigitalAsset) *ListDigitalAssetsResponse
	ResponseGetMetadata(itemIndex int, org *entity.Organization, product *entity.Product) *MetadataResponse
	ResponseGetDetailDigitalAssets(digitalAssets *[]entity.DigitalAsset, collections *[]entity.DigitalAssetCollection, owners *[]entity.User) *ListDigitalAssetsResponse
}

// NewPresenterDigitalAsset Constructs presenter
func NewPresenterDigitalAsset() ConvertDigitalAsset {
	return &PresenterDigitalAsset{}
}

// Return property data response
func (pp *PresenterDigitalAsset) ResponseDigitalAssets(digitalAssets *[]entity.DigitalAsset) *ListDigitalAssetsResponse {
	var response ListDigitalAssetsResponse
	for _, asset := range *digitalAssets {
		res := DigitalAssetResponse{
			ID:           asset.ID.Hex(),
			CollectionID: asset.CollectionID.Hex(),
			TokenID:      asset.TokenID,
			Metadata:     asset.Metadata,
		}
		response.DigitalAssets = append(response.DigitalAssets, res)
	}

	return &response
}

func (pp *PresenterDigitalAsset) ResponseGetMetadata(itemIndex int, org *entity.Organization, product *entity.Product) *MetadataResponse {
	var metadata = MetadataResponse{
		ExternalURL: "https://nomion.io/",
		Image:       product.Image.URL,
	}

	switch org.NameTag {
	case "lej":
		metadata.Name = product.ProductName + " #" + strconv.Itoa(itemIndex)
		metadata.Attributes = append(metadata.Attributes, MetadataAttribute{
			TraitType: "Origin",
			Value:     product.Origin,
		})
		attribute := product.Attribute.(entity.AttributeCoffee)
		metadata.Attributes = append(metadata.Attributes, MetadataAttribute{
			TraitType: "Farm",
			Value:     attribute.FarmName,
		})
		metadata.Attributes = append(metadata.Attributes, MetadataAttribute{
			TraitType: "Preprocessed Type",
			Value:     attribute.Process,
		})
		metadata.Attributes = append(metadata.Attributes, MetadataAttribute{
			TraitType: "Acidity",
			Value:     attribute.Acidity,
		})
		metadata.Attributes = append(metadata.Attributes, MetadataAttribute{
			TraitType: "Bitter",
			Value:     attribute.Bitter,
		})
		metadata.Attributes = append(metadata.Attributes, MetadataAttribute{
			TraitType: "Sweet",
			Value:     attribute.Sweet,
		})
	case "da-non-nuoc":
		metadata.Name = product.ProductName
		attribute := product.Attribute.(entity.AttributeSculpture)
		metadata.Attributes = append(metadata.Attributes, MetadataAttribute{
			TraitType: "Tác Giả",
			Value:     attribute.Craftsman.Name,
		})
		metadata.Attributes = append(metadata.Attributes, MetadataAttribute{
			TraitType: "Cuộc Thi",
			Value:     "Đá Non Nước 2023",
		})
		stoneTarget := attribute.Stone.Translation
		materialName := (stoneTarget["vi"].(map[string]string))["name"]
		metadata.Attributes = append(metadata.Attributes, MetadataAttribute{
			TraitType: "Chất Liệu",
			Value:     fmt.Sprintf("%v", materialName),
		})
		metadata.Attributes = append(metadata.Attributes, MetadataAttribute{
			TraitType: "Chiều Dài",
			Value:     attribute.SculptureLength,
		})
		metadata.Attributes = append(metadata.Attributes, MetadataAttribute{
			TraitType: "Chiều Rộng",
			Value:     attribute.SculptureWidth,
		})
		metadata.Attributes = append(metadata.Attributes, MetadataAttribute{
			TraitType: "Chiều Cao",
			Value:     attribute.SculptureHeight,
		})
		// TODO - Remove HTML tag
		viDescription := attribute.Translation["vi"].(map[string]string)["description"]
		metadata.Description = fmt.Sprintf("%v", viDescription)
		metadata.AnimationURL = product.ThreeDimension.URL
	case "astronaut":
		metadata.Name = product.ProductName
		metadata.AnimationURL = product.Video.URL
		metadata.Image = product.Video.ThumbnailURL
		metadata.Attributes = append(metadata.Attributes, MetadataAttribute{
			TraitType: "Owner",
			Value:     strings.SplitN(product.ProductName, " ", 2)[1],
		})
	default:
		metadata.Name = product.ProductName + " #" + strconv.Itoa(itemIndex)
		return &metadata
	}

	return &metadata
}

func (pp *PresenterDigitalAsset) ResponseGetDetailDigitalAssets(digitalAssets *[]entity.DigitalAsset, collections *[]entity.DigitalAssetCollection, owners *[]entity.User) *ListDigitalAssetsResponse {
	var response ListDigitalAssetsResponse
	for i, digitalAsset := range *digitalAssets {
		res := DigitalAssetResponse{
			ID:              digitalAsset.ID.Hex(),
			Status:          digitalAsset.Status,
			CollectionID:    digitalAsset.CollectionID.Hex(),
			TokenID:         digitalAsset.TokenID,
			Chain:           (*collections)[i].Chain,
			ChainID:         (*collections)[i].ChainID,
			CollectionName:  (*collections)[i].Name,
			ContractAddress: (*collections)[i].ContractAddress,
			Standard:        (*collections)[i].Standard,
			OwnerEmail:      (*owners)[i].Email,
			OwnerName:       (*owners)[i].Name,
			Metadata:        digitalAsset.Metadata,
		}
		response.DigitalAssets = append(response.DigitalAssets, res)
	}
	return &response
}
