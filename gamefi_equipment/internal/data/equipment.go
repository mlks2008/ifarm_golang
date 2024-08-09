package data

import (
	"components/common/utils/datetime"
	"context"
	"gamefi_equipment/api/in/v1/vo"
	"gamefi_equipment/internal/data/ent"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pg/pg/v10"
)

type EquipmentRepo interface {
	Save(context.Context, *ent.EquipmentEnt) error
	UpdateUpgrade(context.Context, int64, int32, []*ent.Attr) error
	UpdateHero(context.Context, int64, int64, int64, int32) error
	UpdateStatus(context.Context, int64, int32) error
	FindByID(context.Context, int64) (*ent.EquipmentEnt, error)
	FindByUserID(context.Context, int64) ([]*ent.EquipmentEnt, error)
	FindByUserHeroID(context.Context, int64) ([]*ent.EquipmentEnt, error)
}

type equipmentRepo struct {
	*Data
	log *log.Helper
}

func NewEquipmentRepo(data *Data, logger log.Logger) EquipmentRepo {
	return &equipmentRepo{
		Data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *equipmentRepo) Save(ctx context.Context, equipment *ent.EquipmentEnt) (err error) {
	_, err = r.pgCli.ModelContext(ctx, equipment).Insert()
	if err != nil {
		r.log.Warnf("[Save] create new data:%+v, error:%v", equipment, err)
	}

	return
}

func (r *equipmentRepo) UpdateUpgrade(ctx context.Context, id int64, level int32, upgrade_attrs []*ent.Attr) (err error) {
	_, err = r.pgCli.ModelContext(ctx, ent.TableEquipment).
		Set("level = ?", level).
		Set("upgrade_attrs = ?", upgrade_attrs).
		Set("update_time = ?", datetime.Now().Time).
		Where("id = ?", id).
		Update()
	if err != nil {
		r.log.Warnf("[UpdateUpgrade] id:%v update error:%v", id, err)
	}

	return
}

func (r *equipmentRepo) UpdateHero(ctx context.Context, id int64, hero_id int64, user_hero_id int64, status int32) (err error) {
	_, err = r.pgCli.ModelContext(ctx, ent.TableEquipment).
		Set("status = ?", status).
		Set("hero_id = ?", hero_id).
		Set("user_hero_id = ?", user_hero_id).
		Set("update_time = ?", datetime.Now().Time).
		Where("id = ?", id).
		Update()
	if err != nil {
		r.log.Warnf("[UpdateHero] id:%v update error:%v", id, err)
	}

	return
}

func (r *equipmentRepo) UpdateStatus(ctx context.Context, id int64, status int32) (err error) {
	_, err = r.pgCli.ModelContext(ctx, ent.TableEquipment).
		Set("status = ?", status).
		Set("update_time = ?", datetime.Now().Time).
		Where("id = ?", id).
		Update()
	if err != nil {
		r.log.Warnf("[UpdateStatus] id:%v update error:%v", id, err)
	}

	return
}

func (r *equipmentRepo) FindByID(ctx context.Context, id int64) (equipmentEnt *ent.EquipmentEnt, err error) {
	equipmentEnt = &ent.EquipmentEnt{
		Id: id,
	}
	err = r.pgCli.Model(equipmentEnt).Where("id = ?", id).Select()

	if err != nil && err != pg.ErrNoRows {
		return nil, err
	} else if err != nil && err == pg.ErrNoRows {
		return nil, nil
	}

	return
}

func (r *equipmentRepo) FindByUserID(ctx context.Context, user_id int64) (equipmentEnts []*ent.EquipmentEnt, err error) {
	err = r.pgCli.Model(ent.TableEquipment).
		Where("user_id = ?", user_id).
		WhereIn("status in (?)", []vo.EquipmentStatus{vo.EquipmentStatus_Normal, vo.EquipmentStatus_Equipped, vo.EquipmentStatus_Locked}).
		Order("create_time desc").
		Select(&equipmentEnts)

	return
}

func (r *equipmentRepo) FindByUserHeroID(ctx context.Context, user_hero_id int64) (equipmentEnts []*ent.EquipmentEnt, err error) {
	err = r.pgCli.Model(ent.TableEquipment).
		Where("user_hero_id = ?", user_hero_id).
		Select(&equipmentEnts)

	return
}
