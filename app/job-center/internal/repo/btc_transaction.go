package repo

import "job-center/internal/model"

type BtcTransactionRepo interface {
	Save(transaction *model.BtcTransaction) error
	GetByTxId(txId string) (*model.BtcTransaction, error)
}
