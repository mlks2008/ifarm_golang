package biz

import (
	"context"
	"gamefi_equipment/internal/data"
	"gamefi_equipment/internal/data/ent"
	"github.com/go-kratos/kratos/v2/log"
)

type EquipmentUC struct {
	repo data.EquipmentRepo
	log  *log.Helper
}

func NewEquipmentUC(repo data.EquipmentRepo, logger log.Logger) *EquipmentUC {
	return &EquipmentUC{repo: repo, log: log.NewHelper(logger)}
}

func (uc *EquipmentUC) Save(ctx context.Context, equipment *ent.EquipmentEnt) error {
	return uc.repo.Save(ctx, equipment)
}

func (uc *EquipmentUC) UpdateUpgrade(ctx context.Context, id int64, level int32, upgrade_attrs []*ent.Attr) error {
	return uc.repo.UpdateUpgrade(ctx, id, level, upgrade_attrs)
}

func (uc *EquipmentUC) UpdateHero(ctx context.Context, id int64, hero_id int64, user_hero_id int64, status int32) (err error) {
	return uc.repo.UpdateHero(ctx, id, hero_id, user_hero_id, status)
}

func (uc *EquipmentUC) UpdateStatus(ctx context.Context, id int64, status int32) error {
	return uc.repo.UpdateStatus(ctx, id, status)
}

func (uc *EquipmentUC) FindByID(ctx context.Context, id int64) (equipmentEnt *ent.EquipmentEnt, err error) {
	equipmentEnt, err = uc.repo.FindByID(ctx, id)

	return
}

func (uc *EquipmentUC) FindByUserID(ctx context.Context, user_id int64) (equipmentEnts []*ent.EquipmentEnt, err error) {
	equipmentEnts, err = uc.repo.FindByUserID(ctx, user_id)

	return
}

func (uc *EquipmentUC) FindByUserHeroID(ctx context.Context, user_hero_id int64) (equipmentEnts []*ent.EquipmentEnt, err error) {
	equipmentEnts, err = uc.repo.FindByUserHeroID(ctx, user_hero_id)

	return
}

//func (uc *EquipmentUC) FindByUserID(ctx context.Context, user_id int64) (equipmentEnts []*ent.EquipmentEnt, err error) {
//	ids, err := uc.RedisCli.SMembers(ctx, fmt.Sprintf(DB_SET_USER_EQUIPMENT, conf.Prefix, user_id)).Result()
//	if err != nil {
//		return nil, err
//	}
//
//	equipCaches, err := uc.getMultiEquipment(user_id, ids)
//	if err != nil {
//		return nil, err
//	}
//
//	if len(equipmentEnts) > 0 {
//		equipmentEnts = make([]*ent.EquipmentEnt, 0, len(equipmentEnts))
//		for _, equip := range equipCaches {
//			if equip.Status == int32(vo.EquipmentStatus_Normal) || equip.Status == int32(vo.EquipmentStatus_Equipped) || equip.Status == int32(vo.EquipmentStatus_Locked) {
//				equipmentEnts = append(equipmentEnts, equip)
//			}
//		}
//		return
//	} else {
//		equipmentEnts, err = uc.repo.FindByUserID(ctx, user_id)
//		if err != nil {
//			return nil, err
//		}
//		return
//	}
//}
//
//func (uc *EquipmentUC) FindByUserHeroID(ctx context.Context, user_id int64, user_hero_id int64) (equipmentEnts []*ent.EquipmentEnt, err error) {
//	list, err := uc.FindByUserID(ctx, user_id)
//	if err != nil {
//		return nil, err
//	}
//
//	equipmentEnts = make([]*ent.EquipmentEnt, 0, len(list))
//	for _, equip := range list {
//		if equip.UserHeroId == user_hero_id {
//			equipmentEnts = append(equipmentEnts, equip)
//		}
//	}
//	return
//}
//
//func (uc *EquipmentUC) getMultiEquipment(user_id int64, id_list []string) ([]*ent.EquipmentEnt, error) {
//	equipments := make([]*ent.EquipmentEnt, 0, len(id_list))
//	return equipments, nil
//}
