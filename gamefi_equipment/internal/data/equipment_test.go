package data

import (
	clog "components/log/zaplogger"
	"context"
	"fmt"
	"gamefi_equipment/internal/conf"
	"gamefi_equipment/internal/data/ent"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"
)

func getData() *Data {
	c := config.New(
		config.WithSource(
			file.NewSource(fmt.Sprintf("../../configs/dev")),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	data, _, _ := NewData(bc.Data, bc.Sys, bc.Log, initLogger(bc.Log))
	return data
}
func initLogger(c *conf.Log) log.Logger {
	caller := func(depth int) log.Valuer {
		return func(context.Context) interface{} {
			_, file, line, _ := runtime.Caller(depth)
			idx := strings.LastIndexByte(file, '/')
			if idx == -1 {
				return file[idx+1:] + ":" + strconv.Itoa(line)
			}
			for i := 0; i < 2; i++ {
				idx = strings.LastIndexByte(file[:idx], '/')
			}
			return file[idx+1:] + ":" + strconv.Itoa(line)
		}
	}
	return log.With(clog.NewLoggerWithName(c.Dir, c.Name, "", c.Debug),
		"ts", log.DefaultTimestamp,
		"caller", caller(4),
		"service.id", "id",
		"service.name", c.Name,
		"service.version", "Version",
	)
}

func Test_Save(t *testing.T) {
	equipmentRepo := NewEquipmentRepo(getData(), nil)

	equipment := &ent.EquipmentEnt{
		Id:           4,
		BaseId:       11112,
		Position:     2,
		UserId:       10000,
		HeroId:       0,
		UserHeroId:   0,
		Level:        1,
		Star:         1,
		Status:       1,
		DropAttrs:    []*ent.Attr{&ent.Attr{Attr: "hp", Value: 1}, &ent.Attr{Attr: "hp", Value: 2}, &ent.Attr{Attr: "atk", Value: 3}},
		UpgradeAttrs: []*ent.Attr{&ent.Attr{Attr: "hp", Value: 1}, &ent.Attr{Attr: "hp", Value: 1}, &ent.Attr{Attr: "atk", Value: 1}},
		CreateTime:   time.Now(),
		Updatetime:   time.Now(),
	}

	err := equipmentRepo.Save(context.Background(), equipment)
	if err != nil {
		t.Error(err)
	} else {
		t.Log("save ok")
	}
}

func Test_FindByUserID(t *testing.T) {
	equipmentRepo := NewEquipmentRepo(getData(), nil)
	data, err := equipmentRepo.FindByUserID(context.Background(), 10000)
	if err != nil {
		t.Error(err)
	} else {
		for _, r := range data {
			t.Log(r)
		}
	}
}
