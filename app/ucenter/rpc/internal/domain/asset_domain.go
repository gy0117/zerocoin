package domain

import (
	"context"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"ucenter-rpc/internal/dao"
)

type AssetDomain struct {
	assetDao *dao.AssetDao
}

func NewAssetDomain(db *gorm.DB) *AssetDomain {
	return &AssetDomain{
		assetDao: dao.NewAssetDao(db),
	}
}

//  1. 冻结之前，先查看钱包的钱是否足够
//
// 2. 冻结：将对应的余额减掉，冻结余额增加
func (ad *AssetDomain) Freeze(ctx context.Context, uid int64, money float64, symbol string) error {
	wallet, err := ad.assetDao.GetByCoinName(ctx, uid, symbol)
	if err != nil {
		return err
	}
	if wallet.Balance < money {
		return errors.New("freeze user asset, but not enough balance")
	}
	return ad.assetDao.Freeze(ctx, uid, money, symbol)
}

// 冻结的钱 - money   余额 + money
func (ad *AssetDomain) UnFreeze(ctx context.Context, uid int64, money float64, symbol string) interface{} {
	wallet, err := ad.assetDao.GetByCoinName(ctx, uid, symbol)
	if err != nil {
		return err
	}
	if wallet.FrozenBalance < money {
		return errors.New("unfreeze user asset, but not enough frozenBalance")
	}
	return ad.assetDao.Unfreeze(ctx, uid, money, symbol)
}
