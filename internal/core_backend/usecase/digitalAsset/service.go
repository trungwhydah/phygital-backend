package digitalAsset

import (
	"net/http"
	"strconv"
	"strings"

	"backend-service/internal/core_backend/entity"
)

// Service struct
type Service struct {
	repo Repository
}

// NewService create service
func NewService(r Repository) *Service {
	return &Service{
		repo: r,
	}
}

// GetDigitalAsset get digitalAsset info
func (s *Service) GetDigitalAssetByCollection(collectionID *string) (*[]entity.DigitalAsset, int, error) {
	digitalAssets, err := s.repo.GetDigitalAssetByCollectionID(collectionID)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return digitalAssets, http.StatusOK, nil
}

func (s *Service) CreateDigitalAsset(da *entity.DigitalAsset) (*entity.DigitalAsset, int, error) {
	da, err := s.repo.CreateDigitalAsset(da)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return da, http.StatusOK, nil
}

func (s *Service) UpdateDigitalAsset(da *entity.DigitalAsset) (bool, int, error) {
	ok, err := s.repo.UpdateDigitalAsset(da)
	if err != nil {
		return false, http.StatusInternalServerError, err
	}
	return ok, http.StatusOK, nil
}

func (s *Service) UpdateDigitalAssetMetadata(daID *string, metadata *entity.Metadata) (bool, int, error) {
	ok, err := s.repo.UpdateDigitalAssetMetadata(daID, metadata)
	if err != nil {
		return false, http.StatusInternalServerError, err
	}
	return ok, http.StatusOK, nil
}

// GetDigitalAssetByID
func (s *Service) GetDigitalAssetByTokenID(collectionID *string, tokenID *int) (*entity.DigitalAsset, int, error) {
	digitalAsset, err := s.repo.GetDigitalAssetByTokenID(collectionID, tokenID)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return digitalAsset, http.StatusOK, nil
}

func (s *Service) GetDigitalAssetByID(daID *string) (*entity.DigitalAsset, int, error) {
	digitalAsset, err := s.repo.GetDigitalAssetByID(daID)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return digitalAsset, http.StatusOK, nil
}

func (s *Service) GetAllActiveDigitalAssets() (*[]entity.DigitalAsset, int, error) {
	digitalAssets, err := s.repo.GetAllActiveDigitalAssets()
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return digitalAssets, http.StatusOK, nil
}

func (s *Service) GetActiveDigitalAssetByCollectionID(collectionID *string) (*[]entity.DigitalAsset, int, error) {
	digitalAssets, err := s.repo.GetActiveDigitalAssetByCollectionID(collectionID)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return digitalAssets, http.StatusOK, nil
}

func (s *Service) GetDigitalAssetsProductAggregate() (*[]entity.DigitalAssetProductAggregate, int, error) {
	aggregations, err := s.repo.GetDigitalAssetsProductAggregate()
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	for i := range *aggregations {
		(*aggregations)[i].Product.ParseAttribute()
	}

	return aggregations, http.StatusOK, nil
}

func (s *Service) ConstructMetadata(itemIndex int, orgTagName *string, product *entity.Product) *entity.Metadata {
	var metadata = entity.Metadata{
		ExternalURL: "https://nomion.io/",
		Image:       product.Image.URL,
	}

	switch *orgTagName {
	case "lej":
		metadata.Name = product.ProductName + " #" + strconv.Itoa(itemIndex)
		metadata.Attributes = append(metadata.Attributes, entity.MetadataAttribute{
			TraitType: "Origin",
			Value:     product.Origin,
		})
		attribute := product.Attribute.(entity.AttributeCoffee)
		metadata.Attributes = append(metadata.Attributes, entity.MetadataAttribute{
			TraitType: "Farm",
			Value:     attribute.FarmName,
		})
		metadata.Attributes = append(metadata.Attributes, entity.MetadataAttribute{
			TraitType: "Preprocessed Type",
			Value:     attribute.Process,
		})
		metadata.Attributes = append(metadata.Attributes, entity.MetadataAttribute{
			TraitType: "Acidity",
			Value:     attribute.Acidity,
		})
		metadata.Attributes = append(metadata.Attributes, entity.MetadataAttribute{
			TraitType: "Bitter",
			Value:     attribute.Bitter,
		})
		metadata.Attributes = append(metadata.Attributes, entity.MetadataAttribute{
			TraitType: "Sweet",
			Value:     attribute.Sweet,
		})
	case "da-non-nuoc":
		metadata.Name = product.ProductName
		attribute := product.Attribute.(entity.AttributeSculpture)
		metadata.Attributes = append(metadata.Attributes, entity.MetadataAttribute{
			TraitType: "Tác Giả",
			Value:     attribute.Craftsman.Name,
		})
		metadata.Attributes = append(metadata.Attributes, entity.MetadataAttribute{
			TraitType: "Cuộc Thi",
			Value:     "Đá Non Nước 2023",
		})
		metadata.Attributes = append(metadata.Attributes, entity.MetadataAttribute{
			TraitType: "Chất Liệu",
			Value:     attribute.Stone.Name,
		})
		metadata.Attributes = append(metadata.Attributes, entity.MetadataAttribute{
			TraitType: "Chiều Dài",
			Value:     attribute.SculptureLength,
		})
		metadata.Attributes = append(metadata.Attributes, entity.MetadataAttribute{
			TraitType: "Chiều Rộng",
			Value:     attribute.SculptureWidth,
		})
		metadata.Attributes = append(metadata.Attributes, entity.MetadataAttribute{
			TraitType: "Chiều Cao",
			Value:     attribute.SculptureHeight,
		})
		metadata.Description = attribute.Description
		metadata.AnimationURL = product.ThreeDimension.URL
	case "astronaut":
		metadata.Name = product.ProductName
		metadata.AnimationURL = product.Video.URL
		metadata.Image = product.Video.ThumbnailURL
		metadata.Attributes = append(metadata.Attributes, entity.MetadataAttribute{
			TraitType: "Owner",
			Value:     strings.SplitN(product.ProductName, " ", 2)[1],
		})
	default:
		metadata.Name = product.ProductName + " #" + strconv.Itoa(itemIndex)
	}
	return &metadata
}
