package dao

import (
	"context"
	"gorm.io/gorm"
	"ucenter-rpc/internal/model"
	"ucenter-rpc/internal/repo"
	"zero-common/zerodb"
	"zero-common/zerodb/gorms"
)

var _ repo.UserAddressRepo = (*UserAddressDao)(nil)

type UserAddressDao struct {
	conn *gorms.GormConn
}

func NewUserAddressDao(db *zerodb.ZeroDB) *UserAddressDao {
	return &UserAddressDao{
		conn: gorms.New(db.Conn),
	}
}

func (m *UserAddressDao) Save(ctx context.Context, address *model.UserAddress) error {
	session := m.conn.Session(ctx)
	return session.Save(&address).Error
}

func (m *UserAddressDao) FindByCoinId(ctx context.Context, userId int64, coinId int64) (list []*model.UserAddress, err error) {
	session := m.conn.Session(ctx)
	err = session.Model(&model.UserAddress{}).Where("user_id=? and coin_id=? and status=?", userId, coinId, model.StatusNormal).Find(&list).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return
}

func (m *UserAddressDao) FindByUserId(ctx context.Context, userId int64) (list []*model.UserAddress, err error) {
	session := m.conn.Session(ctx)
	err = session.Model(&model.UserAddress{}).Where("user_id=? and status=?", userId, model.StatusNormal).Find(&list).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return
}
