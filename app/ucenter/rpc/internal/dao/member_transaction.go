package dao

import (
	"context"
	"gorm.io/gorm"
	"ucenter-rpc/internal/model"
	"ucenter-rpc/internal/repo"
	"zero-common/tools"
	"zero-common/zerodb"
	"zero-common/zerodb/gorms"
)

var _ repo.UserTransactionRepo = (*MemberTransactionDao)(nil)

type MemberTransactionDao struct {
	conn *gorms.GormConn
}

func NewMemberTransactionDao(db *zerodb.ZeroDB) *MemberTransactionDao {
	return &MemberTransactionDao{
		conn: gorms.New(db.Conn),
	}
}

// GetTransactions 找到第pageNo页的数据，一页数据量为pageSize
func (m *MemberTransactionDao) GetTransactions(ctx context.Context, pageNo, pageSize, mId int64, startTime, endTime, symbol string, transactionType int) (list []*model.UserTransaction, total int64, err error) {
	session := m.conn.Session(ctx)
	db := session.Model(&model.UserTransaction{}).Where("user_id=?", mId)
	if symbol != "" {
		db.Where("symbol=?", symbol)
	}
	db.Where("type=?", transactionType)
	if startTime != "" && endTime != "" {
		startTimeMill := tools.ToMill(startTime)
		endTimeMill := tools.ToMill(endTime)
		db.Where("create_time >= ? and create_time <= ?", startTimeMill, endTimeMill)
	}
	offset := (pageNo - 1) * pageSize
	db.Count(&total)
	db.Order("create_time DESC").Limit(int(pageSize)).Offset(int(offset))
	err = db.Find(&list).Error
	if err == gorm.ErrRecordNotFound {
		return nil, 0, err
	}
	return
}

func (m *MemberTransactionDao) FindByAddressAndAmountAndTime(ctx context.Context, address string, value float64, time int64) (mt *model.UserTransaction, err error) {
	session := m.conn.Session(ctx)
	err = session.Model(&model.UserTransaction{}).Where("address=? and amount=? and create_time=?", address, value, time).First(&mt).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return
}

func (m *MemberTransactionDao) Save(ctx context.Context, transaction *model.UserTransaction) error {
	session := m.conn.Session(ctx)
	return session.Save(&transaction).Error
}
