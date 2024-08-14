package domain

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/logx"
	"ucenter-rpc/internal/dao"
	"ucenter-rpc/internal/model"
	"ucenter-rpc/internal/repo"
	"zero-common/btc"
	"zero-common/operate"
	"zero-common/zerodb"
)

type WithdrawDomain struct {
	withdrawRepo repo.WithdrawRepo
	walletDomain *WalletDomain
}

func NewWithdrawDomain(db *zerodb.ZeroDB) *WithdrawDomain {
	return &WithdrawDomain{
		withdrawRepo: dao.NewWithdrawDao(db),
		walletDomain: NewWalletDomain(db),
	}
}

func (wd *WithdrawDomain) SaveRecord(ctx context.Context, conn zerodb.DbConn, record *model.WithdrawRecord) error {
	return wd.withdrawRepo.SaveRecordWithTx(ctx, conn, record)
}

// Withdraw
// 1. 获取当前用户的address
// 2. 根据地址，获取到UTXO
// 3. 判断UTXO中是否符合 提现金额
// 4. 创建交易
// 5. 签名
// 6. 发送到btc网络
// 7. 更新提现状态
func (wd *WithdrawDomain) Withdraw(record model.WithdrawRecord) error {
	ctx := context.Background()
	wallet, err := wd.walletDomain.FindWalletByCoinId(ctx, record.UserId, record.CoinId)
	if err != nil {
		return err
	}

	address := wallet.Address
	logx.Info("Withdraw | address: ", address)
	unspentResultList, err := btc.ListUnspent(0, 999, []string{address})
	if err != nil {
		return err
	}
	totalAmount := record.TotalAmount
	var inputs []btc.Input
	var utxoAmount float64
	for _, v := range unspentResultList {
		utxoAmount += v.Amount
		input := btc.Input{
			Txid: v.Txid,
			Vout: v.Vout,
		}
		inputs = append(inputs, input)
		if utxoAmount >= totalAmount {
			break
		}
	}
	if utxoAmount < totalAmount {
		return errors.New("余额不足")
	}

	// 创建交易
	// 有两个output，假设utxo为0.1，totalAmount为0.002，arrivedAmount为0.00185，fee为0.00015
	var outputs []map[string]any
	m1 := make(map[string]any)
	m1[record.Address] = record.ArrivedAmount // 提现的钱

	m2 := make(map[string]any)
	m2[address] = operate.SubFloor(utxoAmount, totalAmount, 8) // 自己剩下的钱
	outputs = append(outputs, m1, m2)

	hexTxid, err := btc.CreateRawTransaction(inputs, outputs)
	if err != nil {
		return err
	}

	sign, err := btc.SignRawTransactionWithWallet(hexTxid)
	if err != nil {
		return err
	}

	txid, err := btc.SendRawTransaction(sign.Hex)
	if err != nil {
		return err
	}
	record.TransactionNumber = txid
	// 提现成功
	record.Status = 3
	err = wd.withdrawRepo.UpdateTransactionRecord(ctx, record)
	if err != nil {
		logx.Error(err)
	}
	return nil
}

func (wd *WithdrawDomain) WithdrawRecord(ctx context.Context, userId int64, page int64, pageSize int64) ([]*model.WithdrawRecord, int64, error) {
	list, total, err := wd.withdrawRepo.FindByUserId(ctx, userId, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return list, total, err
}
