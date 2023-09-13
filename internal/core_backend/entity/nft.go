package entity

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type EventTransfer struct {
	FromAddr common.Address
	ToAddr   common.Address
	TokenID  *big.Int
}

type Variable struct {
	ID                  string `bson:"_id"`
	LastSyncBlockNumber int64  `bson:"last_sync_block_number"`
}

func (Variable) CollectionName() string {
	return "variables"
}
