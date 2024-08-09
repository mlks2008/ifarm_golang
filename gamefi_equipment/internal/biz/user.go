package biz

import (
	"context"
	"gamefi_equipment/internal/data"
	"gamefi_equipment/internal/data/ent"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pg/pg/v10"
)

type UserUC struct {
	repo data.UserRepo
	log  *log.Helper
}

func NewUserUC(repo data.UserRepo, logger log.Logger) *UserUC {
	return &UserUC{repo: repo, log: log.NewHelper(logger)}
}

func (uc *UserUC) FindCreateUser(ctx context.Context, real_uid string, token_uid string) (user *ent.UserEnt, err error) {
	//uc.log.WithContext(ctx).Debugf("FindCreateUser: %v,%v", real_uid, token_uid)
	user, err = uc.repo.GetByTokenUserId(ctx, token_uid)
	if err != nil && err == pg.ErrNoRows {
		user, err = uc.repo.Save(ctx, real_uid, token_uid)
	}
	return
}

func (uc *UserUC) FindUser(ctx context.Context, token_uid string) (user *ent.UserEnt, err error) {
	user, err = uc.repo.GetByTokenUserId(ctx, token_uid)
	if err != nil {
		if err == pg.ErrNoRows {
			//不存在
			return nil, nil
		} else {
			//报错
			return nil, err
		}
	}
	return
}

func (uc *UserUC) ExistWeb3Address(ctx context.Context, web3_address string) (bool, error) {
	return uc.repo.ExistWeb3Address(ctx, web3_address)
}

func (uc *UserUC) ExistHDAddress(ctx context.Context, hd_address string) (bool, error) {
	return uc.repo.ExistHDAddress(ctx, hd_address)
}

func (uc *UserUC) BindWeb3Address(ctx context.Context, token_userid string, web3_address string) error {
	return uc.repo.BindWeb3Address(ctx, token_userid, web3_address)
}

func (uc *UserUC) BindHDAddress(ctx context.Context, token_userid string, hd_address string) error {
	return uc.repo.BindHDAddress(ctx, token_userid, hd_address)
}

func (uc *UserUC) UnBindWeb3Address(ctx context.Context, token_userid string) error {
	return uc.repo.UnBindWeb3Address(ctx, token_userid)
}

func (uc *UserUC) UnBindHDAddress(ctx context.Context, token_userid string) error {
	return uc.repo.UnBindHDAddress(ctx, token_userid)
}
