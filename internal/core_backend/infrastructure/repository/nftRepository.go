package repository

import (
	"context"

	"backend-service/internal/core_backend/entity"

	"github.com/ethereum/go-ethereum/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// NFTRepository struct
type NFTRepository struct {
	dbMongo *mongo.Database
}

// NewNFTRepository create repository
func NewNFTRepository(dbMongo *mongo.Database) *NFTRepository {
	return &NFTRepository{dbMongo: dbMongo}
}

func (r *NFTRepository) GetLastSyncBlock() (int64, error) {
	var variable entity.Variable
	err := r.dbMongo.Collection(entity.Variable{}.CollectionName()).FindOne(context.TODO(), bson.M{}).Decode(&variable)
	if err != nil {
		return 0, err
	}
	return variable.LastSyncBlockNumber, nil
}

func (r *NFTRepository) UpdateLastSyncBlock(blockNumber uint64) (bool, error) {
	filter := bson.D{{Key: "_id", Value: "blockchain"}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "last_sync_block_number", Value: blockNumber},
		}}}

	result, err := r.dbMongo.Collection(entity.Variable{}.CollectionName()).UpdateOne(
		context.TODO(),
		&filter,
		&update)

	if err != nil {
		return false, err
	}

	return result.MatchedCount != 0, nil
}

func (r *NFTRepository) GetListContractAddresses() (*[]common.Address, error) {
	cursor, err := r.dbMongo.Collection(entity.DigitalAssetCollection{}.CollectionName()).Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	var addresses []common.Address
	// Iterate over the retrieved documents
	for cursor.Next(context.Background()) {
		var col entity.DigitalAssetCollection
		err := cursor.Decode(&col)
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, common.HexToAddress(col.ContractAddress))
	}
	return &addresses, nil
}

func (r *NFTRepository) UpdateMintedDigitalAssets(txHash *string, status *string, tokenID int64) (bool, error) {
	filter := bson.D{{Key: "tx_hash", Value: *txHash}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "status", Value: *status},
			{Key: "token_id", Value: tokenID},
		}}}

	result, err := r.dbMongo.Collection(entity.DigitalAsset{}.CollectionName()).UpdateOne(
		context.TODO(),
		&filter,
		&update)

	if err != nil {
		return false, err
	}

	return result.MatchedCount != 0, nil
}
