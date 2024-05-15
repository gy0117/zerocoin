package domain

import (
	"job-center/internal/dao"
	"job-center/internal/db"
	"job-center/internal/model"
	"job-center/internal/repo"
)

type BtcDomain struct {
	btcTransactionRepo repo.BtcTransactionRepo
}

func NewBtcDomain(cli *db.MongoClient) *BtcDomain {
	return &BtcDomain{
		btcTransactionRepo: dao.NewBtcTransactionDao(cli.Db),
	}
}

// Recharge 其实是保存到mongo中
func (dao *BtcDomain) Recharge(txId, address, blockhash string, value float64, time int64) error {
	transaction, err := dao.btcTransactionRepo.GetByTxId(txId)
	if err != nil {
		return err
	}
	if transaction == nil {
		bt := &model.BtcTransaction{}
		bt.Type = model.RECHARGE
		bt.BlockHash = blockhash
		bt.Value = value
		bt.TxId = txId
		bt.Address = address
		bt.Time = time
		err := dao.btcTransactionRepo.Save(bt)
		if err != nil {
			return err
		}
	}
	return nil
}
