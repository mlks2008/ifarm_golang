package data

import (
	"components/common/utils/datetime"
	"context"
	"gamefi_equipment/internal/data/ent"
	"github.com/go-pg/pg/v10"

	"github.com/go-kratos/kratos/v2/log"
)

type UserRepo interface {
	Save(context.Context, string, string) (*ent.UserEnt, error)
	GetByTokenUserId(context.Context, string) (*ent.UserEnt, error)
	ExistWeb3Address(context.Context, string) (bool, error)
	ExistHDAddress(context.Context, string) (bool, error)
	BindWeb3Address(context.Context, string, string) error
	BindHDAddress(context.Context, string, string) error
	UnBindWeb3Address(context.Context, string) error
	UnBindHDAddress(context.Context, string) error
}

type userRepo struct {
	*Data
	log *log.Helper
}

func NewUserRepo(data *Data, logger log.Logger) UserRepo {
	return &userRepo{
		Data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *userRepo) Save(ctx context.Context, real_userid string, token_userid string) (userEnt *ent.UserEnt, err error) {
	userEnt = &ent.UserEnt{
		RealUserId:  real_userid,
		TokenUserId: token_userid,
		Updatetime:  datetime.Now().Time,
		CreateTime:  datetime.Now().Time,
	}

	_, err = r.pgCli.ModelContext(ctx, userEnt).Insert()
	if err != nil {
		r.log.Warnf("[Save] create new data:%+v, error:%v,", userEnt, err)
	}

	return
}

func (r *userRepo) GetByTokenUserId(ctx context.Context, token_userid string) (userEnt *ent.UserEnt, err error) {
	userEnt = &ent.UserEnt{
		TokenUserId: token_userid,
	}

	err = r.pgCli.ModelContext(ctx, userEnt).Where("token_userid = ?", token_userid).Select()

	if err != nil && err != pg.ErrNoRows {
		r.log.Warnf("[GetByTokenUserId] token_userid:%v error:%v", token_userid, err)
	}

	return
}

func (r *userRepo) ExistWeb3Address(ctx context.Context, web3_address string) (exists bool, err error) {
	var users ent.UserEnt
	exists, err = r.pgCli.Model(&users).Where("web3_address = ?", web3_address).Exists()
	return
}

func (r *userRepo) ExistHDAddress(ctx context.Context, hd_address string) (exists bool, err error) {
	var users ent.UserEnt
	exists, err = r.pgCli.Model(&users).Where("hd_address = ?", hd_address).Exists()
	return
}

func (r *userRepo) BindWeb3Address(ctx context.Context, token_userid string, web3_address string) (err error) {
	//绑定
	_, err = r.pgCli.ModelContext(ctx, ent.TableUser).
		Set("web3_address = ?", web3_address).
		Set("update_time = ?", datetime.Now().Time).
		Where("token_userid = ?", token_userid).
		Update()
	if err != nil {
		r.log.Warnf("[BindWeb3Address] %v update error:%v", token_userid, err)
	}
	return
}

func (r *userRepo) BindHDAddress(ctx context.Context, token_userid string, hd_address string) (err error) {
	//绑定
	_, err = r.pgCli.ModelContext(ctx, ent.TableUser).
		Set("hd_address = ?", hd_address).
		Set("update_time = ?", datetime.Now().Time).
		Where("token_userid = ?", token_userid).
		Update()
	if err != nil {
		r.log.Warnf("[BindHDAddress] %v update error:%v", token_userid, err)
	}

	return
}

func (r *userRepo) UnBindWeb3Address(ctx context.Context, token_userid string) (err error) {
	_, err = r.pgCli.ModelContext(ctx, ent.TableUser).
		Set("web3_address = ?", "").
		Set("update_time = ?", datetime.Now().Time).
		Where("token_userid = ?", token_userid).
		Update()
	if err != nil {
		r.log.Warnf("[UnBindWeb3Address] %v update error:%v", token_userid, err)
	}

	return
}

func (r *userRepo) UnBindHDAddress(ctx context.Context, token_userid string) (err error) {
	_, err = r.pgCli.ModelContext(ctx, ent.TableUser).
		Set("hd_address = ?", "").
		Set("update_time = ?", datetime.Now().Time).
		Where("token_userid = ?", token_userid).
		Update()
	if err != nil {
		r.log.Warnf("[UnBindHDAddress] %v update error:%v", token_userid, err)
	}

	return
}
