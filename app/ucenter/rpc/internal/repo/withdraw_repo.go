package repo

import (
	"context"
	"ucenter-rpc/internal/model"
	"zero-common/zerodb"
)

type WithdrawRepo interface {
	SaveRecordWithTx(ctx context.Context, conn zerodb.DbConn, record *model.WithdrawRecord) error
	UpdateTransactionRecord(ctx context.Context, record model.WithdrawRecord) error
	FindByUserId(ctx context.Context, userId int64, page int64, pageSize int64) ([]*model.WithdrawRecord, int64, error)
}
