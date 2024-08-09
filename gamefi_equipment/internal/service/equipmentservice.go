package service

import (
	"components/common/global"
	"components/common/proto"
	"components/common/utils/random"
	"context"
	"fmt"
	"gamefi_equipment/api"
	"gamefi_equipment/internal/biz"
	"gamefi_equipment/internal/component/kafka"
	"gamefi_equipment/internal/component/redistool"
	"gamefi_equipment/internal/component/sdks"
	"gamefi_equipment/internal/conf"
	"gamefi_equipment/internal/data/ent"
	"gamefi_equipment/internal/resource"
	"github.com/go-kratos/kratos/v2/log"
	"strconv"
	"time"

	pb "gamefi_equipment/api/in/v1"
	vo "gamefi_equipment/api/in/v1/vo"
)

type EquipmentService struct {
	pb.UnimplementedEquipmentServiceServer

	logger *log.Helper
	*sdks.Sdks
	*kafka.Kafka
	redistool *redistool.RedisTool
	*biz.EquipmentUC
}

func NewEquipmentService(equipmentuc *biz.EquipmentUC, redistool *redistool.RedisTool, kafka *kafka.Kafka, sdks *sdks.Sdks, logger log.Logger) *EquipmentService {
	return &EquipmentService{EquipmentUC: equipmentuc, redistool: redistool, Kafka: kafka, Sdks: sdks, logger: log.NewHelper(logger)}
}

func (s *EquipmentService) AddEquipment(ctx context.Context, req *vo.AddEquipmentRequest) (*vo.AddEquipmentResponse, error) {
	resEquipment := resource.EquipmentConfigRes.GetById(req.GetBaseId())
	if resEquipment == nil {
		return nil, global.ToError(int32(api.Error_G_Resource_Not_Exist), "resource equipment not exist")
	}

	var dropAttrs = make([]*ent.Attr, 0)
	var gameAttNum = 0
	var secndAttNUm = 0

	switch resEquipment.Quality {
	case conf.Quality_white:
		gameAttNum = resource.EquipConstantsConfigRes.WhiteMainNum
		secndAttNUm = resource.EquipConstantsConfigRes.WhiteVicNum
	case conf.Quality_blue:
		gameAttNum = resource.EquipConstantsConfigRes.BlueMainNum
		secndAttNUm = resource.EquipConstantsConfigRes.BlueVicNum
	case conf.Quality_purple:
		gameAttNum = resource.EquipConstantsConfigRes.PurpleMainNum
		secndAttNUm = resource.EquipConstantsConfigRes.PurpleVicNum
	case conf.Quality_orange:
		gameAttNum = resource.EquipConstantsConfigRes.OrangeMainNum
		secndAttNUm = resource.EquipConstantsConfigRes.OrangeVicNum
	case conf.Quality_red:
		gameAttNum = resource.EquipConstantsConfigRes.RedMainNum
		secndAttNUm = resource.EquipConstantsConfigRes.RedVicNum
	default:
		return nil, global.ToError(int32(api.Error_S_Quality_Not_Exist), "resource equipment quality not exist")
	}

	if gameAttNum > 0 {
		choiceArr := make([]random.Choice[resource.AttValue, int], 0, len(resEquipment.GameAtt))
		for _, attr := range resEquipment.GameAtt {
			choiceArr = append(choiceArr, random.NewChoice(attr, attr.Weight))
		}
		chooser, _ := random.NewChooser(choiceArr...)

		for i := 0; i < gameAttNum; i++ {
			dropAttr := chooser.Pick()
			rangeVal := random.InRange(dropAttr.Value[0], dropAttr.Value[1]+1)
			dropAttrs = append(dropAttrs, &ent.Attr{Attr: dropAttr.Type, Value: int64(rangeVal)})
		}
	}

	if secndAttNUm > 0 {
		choiceArr := make([]random.Choice[resource.AttValue, int], 0, len(resEquipment.SecndAtt))
		for _, attr := range resEquipment.SecndAtt {
			choiceArr = append(choiceArr, random.NewChoice(attr, attr.Weight))
		}
		chooser, _ := random.NewChooser(choiceArr...)

		for i := 0; i < secndAttNUm; i++ {
			dropAttr := chooser.Pick()
			rangeVal := random.InRange(dropAttr.Value[0], dropAttr.Value[1]+1)
			dropAttrs = append(dropAttrs, &ent.Attr{Attr: dropAttr.Type, Value: int64(rangeVal)})
		}
	}

	equipment := &ent.EquipmentEnt{
		BaseId:     resEquipment.Id,
		Type:       resEquipment.Type,
		Position:   resEquipment.PartType,
		Quality:    resEquipment.Quality,
		UserId:     req.GetUserId(),
		Level:      1,
		Star:       1,
		Status:     int32(vo.EquipmentStatus_Normal),
		DropAttrs:  dropAttrs,
		CreateTime: time.Now(),
		Updatetime: time.Now(),
	}

	err := s.EquipmentUC.Save(context.Background(), equipment)
	if err != nil {
		return nil, global.ToError(int32(api.Error_DB_UNKNOWN), err.Error())
	}

	//掉落日志
	s.Kafka.PubAddEquipmentLog(equipment)

	return &vo.AddEquipmentResponse{
		Code:    int32(proto.Code_OK),
		Message: "",
		Data: &vo.AddEquipmentResponse_Data{
			Success: true,
		},
	}, nil
}

// todo 強化 材料消耗接口未提供
func (s *EquipmentService) UpgradeEquipment(ctx context.Context, req *vo.UpgradeEquipmentRequest) (*vo.UpgradeEquipmentResponse, error) {
	userEquipment, err := s.EquipmentUC.FindByID(ctx, req.GetId())
	if err != nil {
		return nil, global.ToError(int32(api.Error_DB_UNKNOWN), err.Error())
	}
	// 装备不存在
	if userEquipment == nil {
		return nil, global.ToError(int32(api.Error_S_Id_Not_Exist), "user equipment not exist")
	}

	//裝备配制
	resEquipment := resource.EquipmentConfigRes.GetById(userEquipment.BaseId)
	if resEquipment == nil {
		return nil, global.ToError(int32(api.Error_G_Resource_Not_Exist), "resource equipment not exist")
	}

	//装备升级配制
	upgradeLevel := userEquipment.Level + 1
	resEquipmentEnhance := resource.EquipmentEnhanceConfigRes.GetEquipmentEnhance(userEquipment.Quality, upgradeLevel)
	if resEquipmentEnhance == nil {
		return nil, global.ToError(int32(api.Error_S_Equipment_Upgrade_End), "user equipment upgrade end")
	}

	var addValues = make([]*vo.Attr, 0)
	if len(userEquipment.UpgradeAttrs) == 0 {
		for _, attr := range userEquipment.DropAttrs {
			userEquipment.UpgradeAttrs = append(userEquipment.UpgradeAttrs, &ent.Attr{Attr: attr.Attr, Value: 0})
		}
	}
	for i, attr := range userEquipment.DropAttrs {
		if userEquipment.UpgradeAttrs[i].Attr != attr.Attr {
			//升级属性与初始掉落属性不一致
			return nil, global.ToError(int32(api.Error_S_Equipment_Upgrade_Err), "user equipment upgrade err")
		} else {
			var addValue int64
			switch attr.Attr {
			//主属性
			case conf.Attr_hp:
				addValue = resEquipmentEnhance.Hp
			case conf.Attr_atk:
				addValue = resEquipmentEnhance.Atk
			case conf.Attr_def:
				addValue = resEquipmentEnhance.Def
			case conf.Attr_speed:
				addValue = resEquipmentEnhance.Speed
			case conf.Attr_pickRegion:
				addValue = resEquipmentEnhance.PickRegion
			case conf.Attr_energy:
				addValue = resEquipmentEnhance.Energy
			//附加属性
			case conf.Attr_wind:
				addValue = resEquipmentEnhance.Wind
			case conf.Attr_fire:
				addValue = resEquipmentEnhance.Fire
			case conf.Attr_water:
				addValue = resEquipmentEnhance.Water
			case conf.Attr_earth:
				addValue = resEquipmentEnhance.Earth
			case conf.Attr_thunder:
				addValue = resEquipmentEnhance.Thunder
			case conf.Attr_shadow:
				addValue = resEquipmentEnhance.Shadow
			case conf.Attr_holy:
				addValue = resEquipmentEnhance.Holy
			case conf.Attr_arcane:
				addValue = resEquipmentEnhance.Arcane
			case conf.Attr_antiWind:
				addValue = resEquipmentEnhance.AntiWind
			case conf.Attr_antiFire:
				addValue = resEquipmentEnhance.AntiFire
			case conf.Attr_antiWater:
				addValue = resEquipmentEnhance.AntiWater
			case conf.Attr_antiEarth:
				addValue = resEquipmentEnhance.AntiEarth
			case conf.Attr_antiThunder:
				addValue = resEquipmentEnhance.AntiThunder
			case conf.Attr_antiShadow:
				addValue = resEquipmentEnhance.AntiShadow
			case conf.Attr_antiHoly:
				addValue = resEquipmentEnhance.AntiHoly
			case conf.Attr_antiArcane:
				addValue = resEquipmentEnhance.AntiArcane
			default:
				return nil, global.ToError(int32(api.Error_S_Equipment_Upgrade_Attr_Not_Exist), "user equipment upgrade attr not exist")
			}
			////配制数据错误
			//if addValue <= 0 {
			//	return nil, global.ToError(int32(api.Error_G_Resource_Value_Invalid), fmt.Sprintf("equipment upgrade value invalid(resid:%v, %v:0)", enhanceRes.Id, attr.Attr))
			//}
			userEquipment.UpgradeAttrs[i].Value += addValue
			addValues = append(addValues, &vo.Attr{Attr: attr.Attr, Value: addValue})
		}
	}

	//todo 材料消耗接口
	costs := resource.EquipEnhanceCostConfigRes.GetEquipEnhanceCost(userEquipment.Quality, upgradeLevel, resEquipment.Universe)
	s.logger.Debug("UpgradeEquipment costs %+v", costs)

	// 是否升级成功
	choiceArr := make([]random.Choice[bool, int64], 0)
	choiceArr = append(choiceArr, random.NewChoice(true, resEquipmentEnhance.SucceedWeight))
	choiceArr = append(choiceArr, random.NewChoice(false, resEquipmentEnhance.FailWeight))
	chooser, _ := random.NewChooser(choiceArr...)
	if chooser.Pick() == false {
		//升级失败
		return &vo.UpgradeEquipmentResponse{
			Code:    int32(proto.Code_OK),
			Message: "",
			Data: &vo.UpgradeEquipmentResponse_Data{
				Id:      req.GetId(),
				Level:   upgradeLevel - 1,
				Success: false,
			},
		}, nil
	} else {
		//升级成功
		err = s.EquipmentUC.UpdateUpgrade(ctx, req.GetId(), upgradeLevel, userEquipment.UpgradeAttrs)
		if err != nil {
			return nil, global.ToError(int32(api.Error_DB_UNKNOWN), err.Error())
		}

		//升级日志
		userEquipment.Level = upgradeLevel
		s.Kafka.PubUpgradeEquipmentLog(userEquipment)

		return &vo.UpgradeEquipmentResponse{
			Code:    int32(proto.Code_OK),
			Message: "",
			Data: &vo.UpgradeEquipmentResponse_Data{
				Id:        req.GetId(),
				Level:     upgradeLevel,
				Success:   true,
				AddValues: addValues,
			},
		}, nil
	}
}

// todo 分解 材料增加、消耗接口未提供
func (s *EquipmentService) BreakDownEquipment(ctx context.Context, req *vo.BreakDownEquipmentRequest) (*vo.BreakDownEquipmentResponse, error) {
	userEquipment, err := s.EquipmentUC.FindByID(ctx, req.GetId())
	if err != nil {
		return nil, global.ToError(int32(api.Error_DB_UNKNOWN), err.Error())
	}
	// 装备不存在
	if userEquipment == nil {
		return nil, global.ToError(int32(api.Error_S_Id_Not_Exist), "user equipment not exist")
	}

	//裝备配制
	resEquipment := resource.EquipmentConfigRes.GetById(userEquipment.BaseId)
	if resEquipment == nil {
		return nil, global.ToError(int32(api.Error_G_Resource_Not_Exist), "resource equipment not exist")
	}

	resEquipBreakDown := resource.EquipBreakDownConfigRes.GetEquipBreakDown(userEquipment.Quality, userEquipment.Level, resEquipment.Universe)
	if resEquipBreakDown == nil {
		return nil, global.ToError(int32(api.Error_G_Resource_Not_Exist), "resource equipment breakdown not exist")
	}

	//todo 材料消耗接口
	s.logger.Debug("BreakDownEquipment costs %+v", resEquipBreakDown.BreakDownCost)
	//todo 材料掉落接口
	s.logger.Debug("BreakDownEquipment rewards %+v", resEquipBreakDown.BreakDownReward)

	err = s.EquipmentUC.UpdateStatus(ctx, req.GetId(), int32(vo.EquipmentStatus_BreakDown))
	if err != nil {
		return nil, global.ToError(int32(api.Error_DB_UNKNOWN), err.Error())
	}

	return &vo.BreakDownEquipmentResponse{
		Code:    int32(proto.Code_OK),
		Message: "",
		Data: &vo.BreakDownEquipmentResponse_Data{
			Success: false,
		},
	}, nil
}

func (s *EquipmentService) AddFightEquipment(ctx context.Context, req *vo.AddFightEquipmentRequest) (*vo.AddFightEquipmentResponse, error) {

	userEquipment, err := s.EquipmentUC.FindByID(ctx, req.GetId())
	if err != nil {
		return nil, global.ToError(int32(api.Error_DB_UNKNOWN), err.Error())
	}
	// 装备不存在
	if userEquipment == nil {
		return nil, global.ToError(int32(api.Error_S_Id_Not_Exist), "user equipment not exist")
	}

	userHeroData, err := s.GamefiPlatformCli.Hero_GetId(ctx, fmt.Sprintf("%v", req.GetUserHeroId()))
	if err != nil {
		return nil, global.ToError(int32(api.Error_SDK_Api_Failed), err.Error())
	}
	// 用户英雄不存在
	if userHeroData.GetData().GetHeroId() == 0 {
		return nil, global.ToError(int32(api.Error_SDK_Hero_Id_Not_Exist), "user hero not exist")
	}
	// 用户不匹配
	if userEquipment.UserId != userHeroData.GetData().GetUserId() {
		return nil, global.ToError(int32(api.Error_G_User_Not_Match), "userid not match")
	}

	resEquipment := resource.EquipmentConfigRes.GetById(userEquipment.BaseId)
	if resEquipment == nil {
		return nil, global.ToError(int32(api.Error_G_Resource_Not_Exist), "resource equipment not exist")
	}
	// 穿装备位置与配制装备位置不匹配
	if resEquipment.PartType != req.GetPosition() {
		return nil, global.ToError(int32(api.Error_S_Position_Not_Match), "position not match")
	}
	// 专属英雄限制
	if resEquipment.Legend > 0 && resEquipment.Legend != userHeroData.GetData().GetHeroId() {
		return nil, global.ToError(int32(api.Error_S_Legend_Not_Match), "legend not match")
	}
	// 判断英雄的品质，星级
	for _, condition := range resEquipment.UseCondition {
		switch condition.Type {
		case int(conf.ConditionType_quality):
			value, err := strconv.Atoi(fmt.Sprintf("%v", condition.Value))
			if err != nil {
				return nil, global.ToError(int32(api.Error_G_Resource_Value_Invalid), "resource value invalid")
			}
			if userHeroData.GetData().GetQuality() < int32(value) {
				return nil, global.ToError(int32(api.Error_S_Quality_Not_Meet_Condition), "quality not meet condition")
			}
		case int(conf.ConditionType_star):
			value, err := strconv.Atoi(fmt.Sprintf("%v", condition.Value))
			if err != nil {
				return nil, global.ToError(int32(api.Error_G_Resource_Value_Invalid), "resource value invalid")
			}
			if userHeroData.GetData().GetStar() < int32(value) {
				return nil, global.ToError(int32(api.Error_S_Star_Not_Meet_Condition), "star not meet condition")
			}
		}
	}

	userHeroEquips, err := s.EquipmentUC.FindByUserHeroID(ctx, req.GetUserHeroId())
	if err != nil {
		return nil, global.ToError(int32(api.Error_DB_UNKNOWN), err.Error())
	}

	var exist bool
	var exist_id int64
	for _, equip := range userHeroEquips {
		//重复装备
		if equip.Id == req.GetId() && equip.UserHeroId == req.GetUserHeroId() {
			return nil, global.ToError(int32(api.Error_S_Have_Equipped), "equipped")
		}
		//当前位置已经有装备
		if equip.Position == req.GetPosition() {
			exist = true
			exist_id = equip.Id
			break
		}
	}

	// 位置已经有装备(脱旧装备->穿新装备)，否则直接穿装备
	if exist {
		err = s.EquipmentUC.UpdateHero(ctx, exist_id, 0, 0, int32(vo.EquipmentStatus_Normal))
		if err != nil {
			return nil, global.ToError(int32(api.Error_DB_UNKNOWN), err.Error())
		}
		err = s.EquipmentUC.UpdateHero(ctx, req.GetId(), userHeroData.GetData().GetHeroId(), req.GetUserHeroId(), int32(vo.EquipmentStatus_Equipped))
		if err != nil {
			return nil, global.ToError(int32(api.Error_DB_UNKNOWN), err.Error())
		}
	} else {
		err = s.EquipmentUC.UpdateHero(ctx, req.GetId(), userHeroData.GetData().GetHeroId(), req.GetUserHeroId(), int32(vo.EquipmentStatus_Equipped))
		if err != nil {
			return nil, global.ToError(int32(api.Error_DB_UNKNOWN), err.Error())
		}
	}

	return &vo.AddFightEquipmentResponse{
		Code:    int32(proto.Code_OK),
		Message: "",
		Data: &vo.AddFightEquipmentResponse_Data{
			Success: true,
		},
	}, nil
}

func (s *EquipmentService) ClearFightEquipment(ctx context.Context, req *vo.ClearFightEquipmentRequest) (*vo.ClearFightEquipmentResponse, error) {

	userEquipment, err := s.EquipmentUC.FindByID(ctx, req.GetId())
	if err != nil {
		return nil, global.ToError(int32(api.Error_DB_UNKNOWN), err.Error())
	}
	// 装备不存在
	if userEquipment == nil {
		return nil, global.ToError(int32(api.Error_S_Id_Not_Exist), "user equipment not exist")
	}
	// 还没有装备
	if userEquipment.UserHeroId == 0 {
		return nil, global.ToError(int32(api.Error_G_Invalid_Operation), "not equipped yet")
	}

	err = s.EquipmentUC.UpdateHero(ctx, req.GetId(), 0, 0, int32(vo.EquipmentStatus_Normal))
	if err != nil {
		return nil, global.ToError(int32(api.Error_DB_UNKNOWN), err.Error())
	}

	return &vo.ClearFightEquipmentResponse{
		Code:    int32(proto.Code_OK),
		Message: "",
		Data: &vo.ClearFightEquipmentResponse_Data{
			Success: true,
		},
	}, nil
}

func (s *EquipmentService) ListEquipment(ctx context.Context, req *vo.ListEquipmentRequest) (*vo.ListEquipmentResponse, error) {

	list, err := s.EquipmentUC.FindByUserID(ctx, req.GetUserId())
	if err != nil {
		return nil, global.ToError(int32(api.Error_DB_UNKNOWN), err.Error())
	}

	var equipments = make([]*vo.EquipmentInfo, len(list))
	for i, row := range list {
		equipments[i] = s.getEquipmentInfo(row)
	}

	return &vo.ListEquipmentResponse{
		Code:    int32(proto.Code_OK),
		Message: "",
		Data: &vo.ListEquipmentResponse_Data{
			Equipments: equipments,
		},
	}, nil
}

func (s *EquipmentService) ListHeroEquipment(ctx context.Context, req *vo.ListHeroEquipmentRequest) (*vo.ListHeroEquipmentResponse, error) {
	heroEquipmentList, err := s.EquipmentUC.FindByUserHeroID(ctx, req.GetUserHeroId())
	if err != nil {
		return nil, global.ToError(int32(api.Error_DB_UNKNOWN), err.Error())
	}

	heroEquipments, totalAttr := s.getHeroEquipment(req.GetGameId(), req.GetUniverse(), heroEquipmentList)

	return &vo.ListHeroEquipmentResponse{
		Code:    int32(proto.Code_OK),
		Message: "",
		Data: &vo.ListHeroEquipmentResponse_Data{
			HeroEquipment: &vo.HeroEquipment{
				UserHeroId: req.GetUserHeroId(),
				Equipments: heroEquipments,
				TotalAttr:  totalAttr,
			},
		},
	}, nil
}

func (s *EquipmentService) BatchHeroEquipment(ctx context.Context, req *vo.BatchHeroEquipmentRequest) (*vo.BatchHeroEquipmentResponse, error) {
	batchHeroEquipments := make([]*vo.HeroEquipment, 0)
	for _, userHeroId := range req.GetUserHeroId() {
		heroEquipmentList, err := s.EquipmentUC.FindByUserHeroID(ctx, userHeroId)
		if err != nil {
			return nil, global.ToError(int32(api.Error_DB_UNKNOWN), fmt.Sprintf("%v %v", userHeroId, err.Error()))
		}
		heroEquipments, totalAttr := s.getHeroEquipment(req.GetGameId(), req.GetUniverse(), heroEquipmentList)
		batchHeroEquipments = append(batchHeroEquipments, &vo.HeroEquipment{
			UserHeroId: userHeroId,
			Equipments: heroEquipments,
			TotalAttr:  totalAttr,
		})
	}

	return &vo.BatchHeroEquipmentResponse{
		Code:    int32(proto.Code_OK),
		Message: "",
		Data: &vo.BatchHeroEquipmentResponse_Data{
			BatchHeroEquipments: batchHeroEquipments,
		},
	}, nil
}
