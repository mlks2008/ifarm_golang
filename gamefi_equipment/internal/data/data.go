package data

import (
	"components/database/postgres"
	"components/parser"
	"errors"
	"gamefi_equipment/internal/conf"
	"gamefi_equipment/internal/resource"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"path/filepath"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewUserRepo, NewEquipmentRepo)

type Data struct {
	pgCli *postgres.Client
}

func NewData(c *conf.Data, sys *conf.Sys, l *conf.Log, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}

	if c.Postgres == nil {
		return nil, nil, errors.New("postgres config is nil")
	}

	var pgcli = postgres.NewClient(&postgres.Config{
		ApplicationName: l.Name,
		Dsn:             c.Postgres.Source,
		PoolSize:        int(c.Postgres.PoolSize),
		ReadTimeout:     c.Postgres.ReadTimeout.AsDuration(),
		WriteTimeout:    c.Postgres.WriteTimeout.AsDuration(),
		MinIdleConns:    int(c.Postgres.Idle),
		IdleTimeout:     c.Postgres.IdleTimeout.AsDuration(),
		Debug:           c.Postgres.Debug,
	})

	log.NewHelper(logger).Debug("postgres client success")

	loadResources(sys.ConfigPath)

	return &Data{pgCli: pgcli}, cleanup, nil
}

func loadResources(configpath string) {
	res := parser.NewResource().BuildMonitorPath(filepath.Join(configpath, "equipment")).BuildExtName("txt").BuildReloadTime(1 * 60)
	ps := parser.NewBuilder(res)
	boList := []parser.IDataConfig{resource.EquipmentConfigRes, resource.EquipConstantsConfigRes, resource.EquipmentEnhanceConfigRes, resource.EquipEnhanceCostConfigRes, resource.EquipBreakDownConfigRes}
	for _, data := range boList {
		ps.Register(data)
	}
	ps.Init()
}
