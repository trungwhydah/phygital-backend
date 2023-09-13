package nft

import (
	"github.com/ethereum/go-ethereum/common"
)

// NFT interface
type NFT interface {
	GetLastSyncBlock() (int64, error)
	UpdateLastSyncBlock(uint64) (bool, error)
	GetListContractAddresses() (*[]common.Address, error)
	UpdateMintedDigitalAssets(*string, *string, int64) (bool, error)
}

// Repository interface
type Repository interface {
	NFT
}

// UseCase interface
type UseCase interface {
	// Interface for usecase - service
	DeployContract()
	Mint(*string, *string) (string, int, error)
	ListenEvent()
	SyncUnreadEvents(lastSyncBlock int64, toBlock int64)
	WatchTransaction(txHash *string)
}
