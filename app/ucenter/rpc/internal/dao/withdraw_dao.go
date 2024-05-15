package dao

import (
	"context"
	"ucenter-rpc/internal/model"
	"ucenter-rpc/internal/repo"
	"zero-common/zerodb"
	"zero-common/zerodb/gorms"
)

var _ repo.WithdrawRepo = (*WithdrawDao)(nil)

type WithdrawDao struct {
	conn *gorms.GormConn
}

func NewWithdrawDao(db *zerodb.ZeroDB) repo.WithdrawRepo {
	return &WithdrawDao{
		conn: gorms.New(db.Conn),
	}
}

func (wd *WithdrawDao) SaveRecordWithTx(ctx context.Context, conn zerodb.DbConn, record *model.WithdrawRecord) error {
	gormConn := (conn).(*gorms.GormConn)
	session := gormConn.Tx(ctx)
	return session.Save(record).Error
}

func (wd *WithdrawDao) UpdateTransactionRecord(ctx context.Context, record model.WithdrawRecord) error {
	session := wd.conn.Session(ctx)

	values := map[string]any{
		"transaction_number": record.TransactionNumber,
		"status":             record.Status,
	}
	return session.Model(&model.WithdrawRecord{}).Where("id=?", record.Id).Updates(values).Error
}

func (wd *WithdrawDao) FindByUserId(ctx context.Context, userId int64, page int64, pageSize int64) (list []*model.WithdrawRecord, total int64, err error) {
	session := wd.conn.Session(ctx)
	err = session.Model(&model.WithdrawRecord{}).Where("user_id=?", userId).Limit(int(pageSize)).Offset(int((page - 1) * pageSize)).Find(&list).Error
	err = session.Model(&model.WithdrawRecord{}).Where("user_id=?", userId).Count(&total).Error
	return
}
