package dao

import (
	"context"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"ucenter-rpc/internal/model"
)

type AssetDao struct {
	db *gorm.DB
}

func NewAssetDao(db *gorm.DB) *AssetDao {
	return &AssetDao{db}
}

func (ad *AssetDao) GetByCoinName(ctx context.Context, userId int64, coinName string) (data *model.UserWallet, err error) {
	err = ad.db.Model(&model.UserWallet{}).Where("user_id=? and coin_name=?", userId, coinName).Take(&data).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return
}

func (ad *AssetDao) Freeze(ctx context.Context, uid int64, money float64, symbol string) error {
	sql := "update user_wallet set balance=balance-?, frozen_balance=frozen_balance+? where user_id=? and coin_name=?"
	exec := ad.db.Model(&model.UserWallet{}).Exec(sql, money, money, uid, symbol)
	err := exec.Error
	if err != nil {
		return err
	}
	if exec.RowsAffected <= 0 {
		return errors.New("there is no data to update")
	}
	return nil
}
