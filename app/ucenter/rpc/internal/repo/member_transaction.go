package repo

import (
	"context"
	"ucenter-rpc/internal/model"
)

type UserTransactionRepo interface {
	// GetTransactions 获取交易记录
	GetTransactions(ctx context.Context, pageNo, pageSize, mId int64, startTime, endTime, symbol string, transactionType int) (list []*model.UserTransaction, total int64, err error)
	// FindByAddressAndAmountAndTime 查找是否有该记录
	FindByAddressAndAmountAndTime(ctx context.Context, address string, value float64, time int64) (*model.UserTransaction, error)
	// Save 保存
	Save(ctx context.Context, transaction *model.UserTransaction) error
}
