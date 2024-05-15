package dao

import (
	"context"
	"gorm.io/gorm"
	"ucenter-rpc/internal/model"
	"ucenter-rpc/internal/repo"
	"zero-common/zerodb"
	"zero-common/zerodb/gorms"
)

var _ repo.UserRepo = (*UserDao)(nil)

type UserDao struct {
	conn *gorms.GormConn
}

func NewUserMemberDao(db *zerodb.ZeroDB) repo.UserRepo {
	return &UserDao{
		conn: gorms.New(db.Conn),
	}
}

func (umd *UserDao) FindByPhone(ctx context.Context, phone string) (um *model.User, err error) {
	session := umd.conn.Session(ctx)
	err = session.Model(&model.User{}).Where("mobile_phone=?", phone).Limit(1).Take(&um).Error
	if err == gorm.ErrRecordNotFound {
		// 没有找到记录，
		// 在dao层处理干净
		return nil, nil
	}
	return
}

func (umd *UserDao) FindByUserName(ctx context.Context, username string) (um *model.User, err error) {
	session := umd.conn.Session(ctx)
	err = session.Model(&model.User{}).Where("username=?", username).Limit(1).Take(&um).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return
}

func (umd *UserDao) FindUserById(ctx context.Context, userId int64) (um *model.User, err error) {
	session := umd.conn.Session(ctx)
	err = session.Model(&model.User{}).Where("id=?", userId).Take(&um).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return
}

func (umd *UserDao) Save(ctx context.Context, member *model.User) error {
	session := umd.conn.Session(ctx)
	return session.Save(member).Error
}
