package nft

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strings"
	"time"

	config "backend-service/config/core_backend"
	constant "backend-service/internal/core_backend/common"
	"backend-service/internal/core_backend/common/logger"
	"backend-service/internal/core_backend/contracts"
	"backend-service/internal/core_backend/contracts/astronaut_nft"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	MINT_GAS_LIMIT   = uint64(300000)
	DEPLOY_GAS_LIMIT = uint64(5000000)
)

// Service struct
type Service struct {
	repo   Repository
	client *ethclient.Client
}

// NewService create service
func NewService(r Repository) *Service {
	s := Service{
		repo: r,
	}
	s.client = Init()
	return &s
}

func Init() *ethclient.Client {
	client, err := ethclient.Dial(config.C.NFT.BASE_NET_URL)
	if err != nil {
		log.Println(err)
	}
	return client
}

func (s *Service) DeployContract() {
	privateKey, err := crypto.HexToECDSA(config.C.NFT.PRIVATE_KEY)
	if err != nil {
		log.Println(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Println("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := s.client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Println(err)
	}

	gasPrice, err := s.client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Println(err)
	}

	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)       // in wei
	auth.GasLimit = DEPLOY_GAS_LIMIT // in units
	auth.GasPrice = gasPrice

	// baseURL := "https://backend-staging.kyokai.vn/digital-asset/lej/"
	// address, tx, _, err := lej_nft.DeployLejNft(auth, s.client, baseURL)
	baseURL := "https://backend-dev.kyokai.vn/digital-asset/astronaut/"
	address, tx, _, err := astronaut_nft.DeployAstronautNft(auth, s.client, baseURL)
	if err != nil {
		log.Println(err)
	}

	log.Println("Contract Address: " + address.Hex())
	log.Println(tx.Hash().Hex())
}

func (s *Service) Mint(contractAdd *string, ownerAdd *string) (string, int, error) {
	privateKey, err := crypto.HexToECDSA(config.C.NFT.PRIVATE_KEY)
	if err != nil {
		logger.LogError(err.Error())
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		err = errors.New("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
		logger.LogError(err.Error())
		return "", http.StatusInternalServerError, err
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := s.client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		logger.LogError(err.Error())
		return "", http.StatusInternalServerError, err
	}

	gasPrice, err := s.client.SuggestGasPrice(context.Background())
	if err != nil {
		logger.LogError(err.Error())
		return "", http.StatusInternalServerError, err
	}

	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = MINT_GAS_LIMIT // in units
	auth.GasPrice = gasPrice

	ownerAddress := common.HexToAddress(*ownerAdd)
	tx, err := contracts.SafeMint(contractAdd, s.client, auth, &ownerAddress)
	if err != nil {
		logger.LogError(err.Error())
		return "", http.StatusInternalServerError, err
	}

	return tx.Hash().Hex(), http.StatusOK, nil
}

func (s *Service) SyncUnreadEvents(lastSyncBlock int64, toBlock int64) {
	addresses, err := s.repo.GetListContractAddresses()
	if err != nil {
		log.Println(err)
	}
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(lastSyncBlock),
		ToBlock:   big.NewInt(toBlock),
		Addresses: *addresses,
	}
	logs, err := s.client.FilterLogs(context.Background(), query)
	if err != nil {
		logger.LogError(err.Error())
		return
	}
	for _, vLog := range logs {
		contractAdd := strings.ToLower(vLog.Address.Hex())
		logTransfer, err := contracts.ParseTransfer(&contractAdd, s.client, vLog)
		if err != nil {
			log.Println(err)
		}
		// Not a Transfer event, so skip
		if logTransfer == nil && err == nil {
			continue
		}
		log.Println("An Unsynced Transfer Event Caught!")
		log.Println("From", logTransfer.FromAddr.Hex())
		log.Println("To", logTransfer.ToAddr.Hex())
		log.Println("TokenID", logTransfer.TokenID.Int64())
		txHash := vLog.TxHash.Hex()
		status := constant.StatusTxSuccess
		tokenID := logTransfer.TokenID.Int64()
		ok, err := s.repo.UpdateMintedDigitalAssets(&txHash, &status, tokenID)
		if err != nil {
			logger.LogError(err.Error())
			return
		}
		if !ok {
			logger.LogError("Failed update digital assets minted with txhash " + vLog.TxHash.Hex())
		}

	}
	ok, err := s.repo.UpdateLastSyncBlock(uint64(toBlock))
	if err != nil || !ok {
		logger.LogError("Cannot update last sync block")
		return
	}
}

func (s *Service) ListenEvent() {
	addresses, err := s.repo.GetListContractAddresses()
	if err != nil {
		logger.LogError("Cannot get list of contract addresses")
		return
	}
	if len(*addresses) != 0 && err == nil {
		query := ethereum.FilterQuery{
			Addresses: *addresses,
		}
		logs := make(chan types.Log)
		sub, err := s.client.SubscribeFilterLogs(context.Background(), query, logs)
		if err != nil {
			log.Println(err)
		}
		isSync := false
		lastBlockNumber, err := s.repo.GetLastSyncBlock()
		if err != nil {
			log.Println(err)
		}

		for {
			select {
			case err := <-sub.Err():
				log.Println(err)
			case vLog := <-logs:
				if !isSync {
					s.SyncUnreadEvents(lastBlockNumber, int64(vLog.BlockNumber)-1)
					isSync = true
				}
				contractAdd := strings.ToLower(vLog.Address.Hex())
				logTransfer, err := contracts.ParseTransfer(&contractAdd, s.client, vLog)
				if err != nil {
					log.Println(err)
				}
				// Not a Transfer event, so skip
				if logTransfer == nil && err == nil {
					continue
				}
				log.Println("A Transfer Event Caught!")
				log.Println("From", logTransfer.FromAddr.Hex())
				log.Println("To", logTransfer.ToAddr.Hex())
				log.Println("TokenID", logTransfer.TokenID.Int64())
				txHash := vLog.TxHash.Hex()
				status := constant.StatusTxSuccess
				tokenID := logTransfer.TokenID.Int64()
				ok, err := s.repo.UpdateMintedDigitalAssets(&txHash, &status, tokenID)
				if err != nil {
					log.Println(err)
				}
				if !ok {
					logger.LogError("Failed update digital assets minted with txhash " + vLog.TxHash.Hex())
				}
				_, _ = s.repo.UpdateLastSyncBlock(uint64(vLog.BlockNumber))
			}
		}
	}
}

func (s *Service) WatchTransaction(txHash *string) {
	status := constant.StatusTxPending
	txReceipt, err := s.GetTransactionReceipt(*txHash)
	if err != nil {
		logger.LogError(fmt.Sprintf("[EVM] - Fail to get transaction receipt for %s: %s", *txHash, err.Error()))
	} else {
		vLog := txReceipt.Logs[0]
		tokenID := int64(0)
		if txReceipt.Status == 1 {
			status = constant.StatusTxSuccess
			contractAdd := strings.ToLower(vLog.Address.Hex())
			logTransfer, err := contracts.ParseTransfer(&contractAdd, s.client, *vLog)
			if err != nil {
				log.Println(err)
			}
			// Not a Transfer event, so skip
			if logTransfer == nil && err == nil {
				return
			}
			log.Println("A Transfer Event Caught!")
			log.Println("From", logTransfer.FromAddr.Hex())
			log.Println("To", logTransfer.ToAddr.Hex())
			log.Println("TokenID", logTransfer.TokenID.Int64())
			tokenID = logTransfer.TokenID.Int64()
		} else {
			status = constant.StatusTxFailure
		}
		ok, err := s.repo.UpdateMintedDigitalAssets(txHash, &status, tokenID)
		if err != nil {
			log.Println(err)
		}
		if !ok {
			logger.LogError("Failed update digital assets minted with txhash " + *txHash)
		}
	}
}

func (s *Service) GetTransactionReceipt(txHash string) (*types.Receipt, error) {
	var (
		receipt *types.Receipt
		err     error
		retry   = 0
		hash    = common.HexToHash(txHash)
	)

	for retry < config.C.NFT.WATCH_TRANSACTION_MAX_RETRY {
		receipt, err = s.client.TransactionReceipt(context.Background(), hash)
		if errors.Is(err, ethereum.NotFound) {
			retry++
			if retry < config.C.NFT.WATCH_TRANSACTION_MAX_RETRY {
				time.Sleep(time.Duration(config.C.NFT.WATCH_TRANSACTION_INTERVAL_TIME) * time.Second)
				continue
			}
		}
		break
	}
	if errors.Is(err, ethereum.NotFound) {
		logger.LogError("Couldn't get transaction status after WATCH_TRANSACTION_MAX_RETRY retries!")
	}

	return receipt, err
}
