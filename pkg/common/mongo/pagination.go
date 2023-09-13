package mongo

import (
	"errors"

	paginationpkg "backend-service/pkg/common/pagination"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrorInvalidModel = errors.New("required model")
	ErrorModelStruct  = errors.New("model must be struct")
)

type SortOperationType int8

const (
	AscSortMongo  SortOperationType = 1
	DescSortMongo SortOperationType = -1
)

type CompareOperationType string

const (
	LtOperationMongo CompareOperationType = "$lt"
	GtOperationMongo CompareOperationType = "$gt"
)

// BuildPagePaginationPipeline Pagination using Page
func BuildPagePaginationPipeline(pagination *paginationpkg.Pagination) mongo.Pipeline {
	sortOperator := GetSortOperator(pagination)
	sortStage := bson.D{{Key: "$sort", Value: bson.M{pagination.OrderBy: sortOperator}}}

	skipStage := bson.D{{Key: "$skip", Value: (pagination.Page - 1) * pagination.Limit}}

	limitStage := bson.D{{Key: "$limit", Value: pagination.Limit}}

	pipeline := mongo.Pipeline{sortStage, skipStage, limitStage}

	return pipeline
}

func GetSortOperator(pagination *paginationpkg.Pagination) SortOperationType {
	sortOperator := DescSortMongo

	if pagination.IsAsc() {
		sortOperator = AscSortMongo
	}

	return sortOperator
}
