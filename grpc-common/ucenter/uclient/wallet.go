// Code generated by goctl. DO NOT EDIT.
// Source: wallet.proto

package uclient

import (
	"context"
	"grpc-common/ucenter/types/wallet"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	WalletReq     = wallet.WalletReq
	WalletResp    = wallet.UserWallet
	FindWalletResp = wallet.FindWalletResp
	WalletEmptyResp = wallet.WalletResp
	AddressListResp = wallet.AddressListResp

	AssetReq = wallet.AssetReq
	UserTransactionListResp = wallet.UserTransactionListResp

	FreezeUserAssetReq = wallet.FreezeUserAssetReq
	Empty = wallet.Empty
	DeductUserAssetReq = wallet.DeductUserAssetReq
	AddUserAssetReq = wallet.AddUserAssetReq

	Wallet interface {
		FindWalletBySymbol(ctx context.Context, in *WalletReq, opts ...grpc.CallOption) (*WalletResp, error)
		FindWallet(ctx context.Context, in *WalletReq, opts ...grpc.CallOption) (*FindWalletResp, error)
		ResetWalletAddress(ctx context.Context, in *WalletReq, opts ...grpc.CallOption) (*WalletEmptyResp, error)
		GetAllTransactions(ctx context.Context, in *AssetReq, opts ...grpc.CallOption) (*UserTransactionListResp, error)
		GetAddress(ctx context.Context, in *AssetReq, opts ...grpc.CallOption) (*AddressListResp, error)
		FreezeUserAsset(ctx context.Context, in *FreezeUserAssetReq, opts ...grpc.CallOption) (*Empty, error)
		UnFreezeUserAsset(ctx context.Context, in *FreezeUserAssetReq, opts ...grpc.CallOption) (*Empty, error)
		DeductUserAsset(ctx context.Context, in *DeductUserAssetReq, opts ...grpc.CallOption) (*Empty, error)
		AddUserAsset(ctx context.Context, in *AddUserAssetReq, opts ...grpc.CallOption) (*Empty, error)
	}

	defaultWallet struct {
		cli zrpc.Client
	}
)

func NewWallet(cli zrpc.Client) Wallet {
	return &defaultWallet{
		cli: cli,
	}
}



func (m *defaultWallet) FindWalletBySymbol(ctx context.Context, in *WalletReq, opts ...grpc.CallOption) (*WalletResp, error) {
	client := wallet.NewWalletClient(m.cli.Conn())
	return client.FindWalletBySymbol(ctx, in, opts...)
}


func (m *defaultWallet) FindWallet(ctx context.Context, in *WalletReq, opts ...grpc.CallOption) (*FindWalletResp, error) {
	client := wallet.NewWalletClient(m.cli.Conn())
	return client.FindWallet(ctx, in, opts...)
}

func (m *defaultWallet) ResetWalletAddress(ctx context.Context, in *WalletReq, opts ...grpc.CallOption) (*WalletEmptyResp, error) {
	client := wallet.NewWalletClient(m.cli.Conn())
	return client.ResetWalletAddress(ctx, in, opts...)
}

func (m *defaultWallet) GetAllTransactions(ctx context.Context, in *AssetReq, opts ...grpc.CallOption) (*UserTransactionListResp, error) {
	client := wallet.NewWalletClient(m.cli.Conn())
	return client.GetAllTransactions(ctx, in, opts...)
}

func (m *defaultWallet) GetAddress(ctx context.Context, in *AssetReq, opts ...grpc.CallOption) (*AddressListResp, error) {
	client := wallet.NewWalletClient(m.cli.Conn())
	return client.GetAddress(ctx, in, opts...)
}

func (m *defaultWallet) FreezeUserAsset(ctx context.Context, in *FreezeUserAssetReq, opts ...grpc.CallOption) (*Empty, error) {
	client := wallet.NewWalletClient(m.cli.Conn())
	return client.FreezeUserAsset(ctx, in, opts...)
}

func (m *defaultWallet) UnFreezeUserAsset(ctx context.Context, in *FreezeUserAssetReq, opts ...grpc.CallOption) (*Empty, error) {
	client := wallet.NewWalletClient(m.cli.Conn())
	return client.UnFreezeUserAsset(ctx, in, opts...)
}

func (m *defaultWallet) DeductUserAsset(ctx context.Context, in *DeductUserAssetReq, opts ...grpc.CallOption) (*Empty, error) {
	client := wallet.NewWalletClient(m.cli.Conn())
	return client.DeductUserAsset(ctx, in, opts...)
}

func (m *defaultWallet) AddUserAsset(ctx context.Context, in *AddUserAssetReq, opts ...grpc.CallOption) (*Empty, error) {
	client := wallet.NewWalletClient(m.cli.Conn())
	return client.AddUserAsset(ctx, in, opts...)
}