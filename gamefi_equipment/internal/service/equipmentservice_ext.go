package service

import (
	"components/common/utils/strs"
	"fmt"
	vo "gamefi_equipment/api/in/v1/vo"
	"gamefi_equipment/internal/conf"
	"gamefi_equipment/internal/data/ent"
	"gamefi_equipment/internal/resource"
	"strings"
)

func (s *EquipmentService) getHeroEquipment(gameId int64, universe string, heroEquipmentList []*ent.EquipmentEnt) ([]*vo.EquipmentInfo, []*vo.Attr) {
	var heroEquipments = make([]*vo.EquipmentInfo, 0, len(heroEquipmentList))

	var conditionAttr = make(map[string]int64)
	for _, row := range heroEquipmentList {
		equip := s.getEquipmentInfo(row)

		resEquipment := resource.EquipmentConfigRes.GetById(row.BaseId)
		if resEquipment == nil {
			continue
		}

		//属性生效条件
		var conditionOK bool = true
		var effectiveAtt = make(map[int]string)
		for _, condition := range resEquipment.EffectiveAtt {
			effectiveAtt[condition.Type] = fmt.Sprintf("%v,%v", effectiveAtt[condition.Type], condition.Value)
		}
		for effectiveTyp, effectiveValue := range effectiveAtt {
			// 去除两端逗号
			effectiveValue = strings.Trim(effectiveValue, ",")

			if effectiveTyp == int(conf.ConditionType_game) {
				if strs.Contains(effectiveValue, "0") == false && strs.Contains(effectiveValue, fmt.Sprintf("%v", gameId)) == false {
					conditionOK = false
					break
				}
			} else if effectiveTyp == int(conf.ConditionType_universe) {
				if strs.Contains(effectiveValue, "") == false && strs.Contains(effectiveValue, universe) == false {
					conditionOK = false
					break
				}
			} else {
				conditionOK = false
				break
			}
		}

		if conditionOK {
			for _, attr := range row.DropAttrs {
				conditionAttr[attr.Attr] += attr.Value
			}
			for _, attr := range row.UpgradeAttrs {
				conditionAttr[attr.Attr] += attr.Value
			}
		}
		heroEquipments = append(heroEquipments, equip)
	}

	var totalAttr = make([]*vo.Attr, 0, len(conditionAttr))
	for key, val := range conditionAttr {
		totalAttr = append(totalAttr, &vo.Attr{Attr: key, Value: val})
	}

	return heroEquipments, totalAttr
}

func (s *EquipmentService) getEquipmentInfo(row *ent.EquipmentEnt) *vo.EquipmentInfo {
	equipmentInfo := &vo.EquipmentInfo{
		Id:           row.Id,
		BaseId:       row.BaseId,
		UserId:       row.UserId,
		HeroId:       row.HeroId,
		UserHeroId:   row.UserHeroId,
		Level:        row.Level,
		Star:         row.Star,
		Status:       vo.EquipmentStatus(row.Status),
		InitAttrs:    []*vo.Attr{},
		UpgradeAttrs: []*vo.Attr{},
		Position:     row.Position,
	}

	equipmentInfo.InitAttrs = make([]*vo.Attr, 0)
	for _, attr := range row.DropAttrs {
		equipmentInfo.InitAttrs = append(equipmentInfo.InitAttrs, &vo.Attr{Attr: attr.Attr, Value: attr.Value})
	}
	equipmentInfo.UpgradeAttrs = make([]*vo.Attr, 0)
	for _, attr := range row.UpgradeAttrs {
		equipmentInfo.UpgradeAttrs = append(equipmentInfo.UpgradeAttrs, &vo.Attr{Attr: attr.Attr, Value: attr.Value})
	}

	return equipmentInfo
}
