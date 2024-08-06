package logic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"grpc-common/ucenter/types/wallet"
	"grpc-common/ucenter/uclient"
	"job-center/internal/db"
	"job-center/internal/domain"
	"job-center/internal/model"
	"sync"
	"zero-common/kafka"
	"zero-common/tools"
)

const rdsBtcTxBlockHeight = "rds::btc::tx::block::height"

type BitCoin struct {
	wg          sync.WaitGroup
	url         string
	redisCache  cache.Cache
	walletRpc   uclient.Wallet
	btcDomain   *domain.BtcDomain
	kafkaDomain *domain.QueueDomain
}

func NewBitCoin(url string, redisCache cache.Cache, walletRpc uclient.Wallet, mCli *db.MongoClient, kCli *kafka.KafkaClient) *BitCoin {
	return &BitCoin{
		url:         url,
		redisCache:  redisCache,
		walletRpc:   walletRpc,
		btcDomain:   domain.NewBtcDomain(mCli),
		kafkaDomain: domain.NewQueueDomain(kCli),
	}
}

func (bc *BitCoin) Do() {
	bc.wg.Add(1)
	go bc.scanTx(bc.url)
	bc.wg.Wait()
}

// 查找符合系统url的交易 进行存储
func (bc *BitCoin) scanTx(url string) {
	// 0. 获取系统中的BTC的address地址
	addressListResp, err := bc.walletRpc.GetAddress(context.Background(), &wallet.AssetReq{
		CoinName: "BTC",
	})
	if err != nil {
		logx.Error(err)
		bc.wg.Done()
		return
	}
	addressList := addressListResp.List

	// 1. redis查询是否有记录区块，获取到已处理的区块高度
	var dealBlockHeightStr string
	_ = bc.redisCache.Get(rdsBtcTxBlockHeight, &dealBlockHeightStr)
	var dealBlockHeight int64
	if dealBlockHeightStr == "" {
		dealBlockHeight = 2428713
	} else {
		dealBlockHeight = tools.ToInt64(dealBlockHeightStr)
	}

	// 2. 根据getmininginfo获取到现在的区块高度
	curBlockHeight, err := bc.getMiningInfo(url)
	if err != nil {
		logx.Error(err)
		bc.wg.Done()
		return
	}

	// 3. 计算 curBlockHeight - dealBlockHeight，如果小于等于0，则不需要扫描，表示没有新的交易
	if curBlockHeight <= dealBlockHeight {
		fmt.Println("当前没有新的交易，curBlockHeight：", curBlockHeight, " dealBlockHeight：", dealBlockHeight)
		bc.wg.Done()
		return
	}
	fmt.Printf("bitcoin | %d - %d = %d\n ", curBlockHeight, dealBlockHeight, (curBlockHeight - dealBlockHeight))

	// 4. 循环 根据getblockhash 获取blockhash
	for i := curBlockHeight; i > dealBlockHeight; i-- {
		blockHash, err := bc.getBlockHash(url, i)
		if err != nil {
			logx.Error(err)
			continue
		}
		// 5. 通过getblock，获取交易id列表
		txList, err := bc.getBlock(url, blockHash)
		if err != nil {
			logx.Error(err)
			continue
		}
		// 6. 循环交易id列表，获取到交易详情，得到vount内容
		for _, txId := range txList {
			rawTransaction, err := bc.getRawTransaction(url, txId)
			if err != nil {
				logx.Error(err)
				continue
			}
			// 找到当前交易中的输入地址的币来自哪里
			inputAddressList := make([]string, len(rawTransaction.Vin))
			for index, vin := range rawTransaction.Vin {
				if vin.TxId == "" {
					continue
				}
				transaction, err := bc.getRawTransaction(url, vin.TxId)
				if err != nil {
					logx.Error(err)
					continue
				}
				// 当前交易中的输入地址的钱来自于下面的vout
				vout := transaction.Vout[vin.Vout]
				inputAddressList[index] = vout.ScriptPubKey.Address
			}

			// 7. 判断哪些地址是充值的
			// 先找出vout
			for _, vout := range rawTransaction.Vout {
				voutAddress := vout.ScriptPubKey.Address
				if voutAddress == "" {
					continue
				}
				flag := false
				// 遍历inputAddressList，如果地址一致，说明当前voutAddress不是充值地址
				for _, inputAddress := range inputAddressList {
					if inputAddress != "" && voutAddress == inputAddress {
						flag = true
					}
				}
				// 当前voutAddress地址不是充值地址
				if flag {
					continue
				}

				// 遍历数据库中的地址，找到与voutAddress相同的，然后给这个地址充值
				for _, addr := range addressList {
					if addr != "" && addr == voutAddress {
						// 10. 找到了充值地址后，存入mongo，同时发送kafka进行处理（存入user_transaction表）
						fmt.Println("bitcoin | 找到了当前充值的地址，voutAddress: ", addr)
						err := bc.btcDomain.Recharge(rawTransaction.TxId, addr, rawTransaction.BlockHash, vout.Value, rawTransaction.Time)
						if err != nil {
							logx.Error(err)
							continue
						}
						// kafka
						bc.kafkaDomain.SendRecharge(vout.Value, rawTransaction.Time, addr)
					}
				}
			}
		}
	}

	// 11. 记录当前区块高度
	_ = bc.redisCache.Set(rdsBtcTxBlockHeight, curBlockHeight)
	bc.wg.Done()
}

func (bc *BitCoin) getMiningInfo(url string) (int64, error) {
	//{
	//	"jsonrpc": "1.0",
	//	"method": "getmininginfo",
	//	"params":[123],
	//	"id": "zerocoin"
	//}

	params := make(map[string]any)
	params["jsonrpc"] = "1.0"
	params["method"] = "getmininginfo"
	params["params"] = []int{}
	params["id"] = "zerocoin"

	headers := make(map[string]string)
	// rpcuser:rpcpassword base64编码
	headers["Authorization"] = "Basic Yml0Y29pbjoxMjM0NTY="

	bytes, err := tools.PostWithHeader(url, params, headers, "")
	if err != nil {
		return 0, err
	}
	var miningInfoResp model.MiningInfoResp
	err = json.Unmarshal(bytes, &miningInfoResp)
	if err != nil {
		return 0, err
	}
	if miningInfoResp.Error != "" {
		return 0, errors.New(miningInfoResp.Error)
	}
	return int64(miningInfoResp.Result.Blocks), nil
}

func (bc *BitCoin) getBlockHash(url string, height int64) (string, error) {
	params := make(map[string]any)
	params["jsonrpc"] = "1.0"
	params["method"] = "getblockhash"
	params["params"] = []int64{height}
	params["id"] = "zerocoin"

	headers := make(map[string]string)
	// rpcuser:rpcpassword base64编码
	headers["Authorization"] = "Basic Yml0Y29pbjoxMjM0NTY="

	bytes, err := tools.PostWithHeader(url, params, headers, "")
	if err != nil {
		return "", err
	}
	var blockHashResp model.BlockHashResp
	err = json.Unmarshal(bytes, &blockHashResp)
	if err != nil {
		return "", err
	}
	if blockHashResp.Error != "" {
		return "", errors.New(blockHashResp.Error)
	}
	return blockHashResp.Result, nil
}

func (bc *BitCoin) getBlock(url string, blockHash string) ([]string, error) {
	params := make(map[string]any)
	params["jsonrpc"] = "1.0"
	params["method"] = "getblock"
	params["params"] = []any{blockHash, 1}
	params["id"] = "zerocoin"

	headers := make(map[string]string)
	// rpcuser:rpcpassword base64编码
	headers["Authorization"] = "Basic Yml0Y29pbjoxMjM0NTY="

	bytes, err := tools.PostWithHeader(url, params, headers, "")
	if err != nil {
		return nil, err
	}
	var blockResp model.BlockResp
	err = json.Unmarshal(bytes, &blockResp)
	if err != nil {
		return nil, err
	}
	if blockResp.Error != "" {
		return nil, errors.New(blockResp.Error)
	}
	return blockResp.Result.Tx, nil
}

func (bc *BitCoin) getRawTransaction(url string, txId string) (*model.RawTransactionData, error) {
	params := make(map[string]any)
	params["jsonrpc"] = "1.0"
	params["method"] = "getrawtransaction"
	params["params"] = []any{txId, true}
	params["id"] = "zerocoin"

	headers := make(map[string]string)
	// rpcuser:rpcpassword base64编码
	headers["Authorization"] = "Basic Yml0Y29pbjoxMjM0NTY="

	bytes, err := tools.PostWithHeader(url, params, headers, "")
	if err != nil {
		return nil, err
	}
	var rawTransactionResp model.RawTransactionResp
	err = json.Unmarshal(bytes, &rawTransactionResp)
	if err != nil {
		return nil, err
	}
	if rawTransactionResp.Error != "" {
		return nil, errors.New(rawTransactionResp.Error)
	}
	return &rawTransactionResp.Result, nil
}
