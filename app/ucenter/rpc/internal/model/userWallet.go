package model

import "grpc-common/market/types/market"

type UserWallet struct {
	Id                int64   `gorm:"column:id"`
	Address           string  `gorm:"column:address"`
	Balance           float64 `gorm:"column:balance"`
	FrozenBalance     float64 `gorm:"column:frozen_balance"`
	ReleaseBalance    float64 `gorm:"column:release_balance"`
	IsLock            int     `gorm:"column:is_lock"`
	UserId            int64   `gorm:"column:user_id"`
	Version           int     `gorm:"column:version"`
	CoinId            int64   `gorm:"column:coin_id"`
	ToReleased        float64 `gorm:"column:to_released"`
	CoinName          string  `gorm:"column:coin_name"`
	AddressPrivateKey string  `gorm:"address_private_key"`
}

func (*UserWallet) TableName() string {
	return "user_wallet"
}

type UserWalletCoin struct {
	Id             int64        `json:"id" from:"id"`
	Address        string       `json:"address" from:"address"`
	Balance        float64      `json:"balance" from:"balance"`
	FrozenBalance  float64      `json:"frozenBalance" from:"frozenBalance"`
	ReleaseBalance float64      `json:"releaseBalance" from:"releaseBalance"`
	IsLock         int          `json:"isLock" from:"isLock"`
	UserId         int64        `json:"userId" from:"userId"`
	Version        int          `json:"version" from:"version"`
	Coin           *market.Coin `json:"coin" from:"coinId"`
	ToReleased     float64      `json:"toReleased" from:"toReleased"`
}

func NewWalletData(userId int64, coin *market.Coin) (*UserWallet, *UserWalletCoin) {
	mw := &UserWallet{
		UserId:   userId,
		CoinId:   int64(coin.Id),
		CoinName: coin.Unit,
	}
	mwc := &UserWalletCoin{
		UserId: userId,
		Coin:   coin,
	}
	return mw, mwc
}
