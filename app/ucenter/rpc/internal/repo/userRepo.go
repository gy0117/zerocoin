package repo

import (
	"context"
	"ucenter-rpc/internal/model"
)

type UserRepo interface {
	FindByPhone(ctx context.Context, phone string) (*model.User, error)
	FindUserById(ctx context.Context, userId int64) (*model.User, error)
	Save(ctx context.Context, member *model.User) error
	FindByUserName(ctx context.Context, username string) (*model.User, error)
}
