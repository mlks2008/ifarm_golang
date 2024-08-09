package gtscenter

import (
	"components/sdks/gamefi_platform"
	"components/sdks/gts_shop"
	"github.com/go-kratos/kratos/v2/log"
	"sync"
)

type GtsCenter struct {
	logger          *log.Helper
	mutex           sync.Mutex
	gtsShopCli      *gts_shop.Client        //user与gp
	platformDustCli *gamefi_platform.Client //dust
}

func NewGtsCenter(gtsShopCli *gts_shop.Client, platformDustCli *gamefi_platform.Client, logger log.Logger) *GtsCenter {
	c := &GtsCenter{}
	c.logger = log.NewHelper(logger)
	c.gtsShopCli = gtsShopCli
	c.platformDustCli = platformDustCli

	return c
}
