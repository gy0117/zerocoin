package domain

import (
	"context"
	"ucenter-rpc/internal/dao"
	"ucenter-rpc/internal/model"
	"ucenter-rpc/internal/repo"
	"zero-common/tools"
	"zero-common/zerodb"
)

type TransactionDomain struct {
	mtRepo       repo.UserTransactionRepo
	walletDomain *WalletDomain
}

func NewTransactionDomain(db *zerodb.ZeroDB) *TransactionDomain {
	return &TransactionDomain{
		mtRepo:       dao.NewMemberTransactionDao(db),
		walletDomain: NewWalletDomain(db),
	}
}

func (d *TransactionDomain) GetTransactions(ctx context.Context, pageNo int64, pageSize int64, mId int64, symbol string, startTime string, endTime string, transactionType string) ([]*model.UserTransactionVo, int64, error) {
	tt := tools.ToInt64(transactionType)
	list, total, err := d.mtRepo.GetTransactions(ctx, pageNo, pageSize, mId, startTime, endTime, symbol, int(tt))
	if err != nil {
		return nil, 0, err
	}

	voList := make([]*model.UserTransactionVo, len(list))

	for i, v := range list {
		voList[i] = v.ToVo()
	}
	return voList, total, err
}

// SaveRecharge 存入数据库
func (d *TransactionDomain) SaveRecharge(address string, value float64, time int64, transactionType string, symbol string) error {
	ctx := context.Background()
	memberTransaction, err := d.mtRepo.FindByAddressAndAmountAndTime(ctx, address, value, time)
	if err != nil {
		return err
	}

	// 根据address找到钱包
	memberWallet, err := d.walletDomain.GetByAddress(ctx, address)
	if err != nil {
		return err
	}
	if memberTransaction == nil {
		mt := &model.UserTransaction{}
		mt.UserId = memberWallet.UserId
		mt.Address = address
		mt.Type = int64(model.TransTypeStr(transactionType))
		mt.CreateTime = time * 1000
		mt.Amount = value
		mt.Symbol = symbol

		if err := d.mtRepo.Save(ctx, mt); err != nil {
			return err
		}
	}
	return nil
}
