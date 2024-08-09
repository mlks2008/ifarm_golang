package sdks

import (
	"components/sdks/gamefi_account"
	"components/sdks/gamefi_platform"
	"components/sdks/gts_shop"
	"components/sdks/gtscenter"
	"gamefi_equipment/internal/conf"
	"github.com/go-kratos/kratos/v2/log"
)

type Sdks struct {
	GtsCenter         *gtscenter.GtsCenter
	GtsShopCli        *gts_shop.Client
	GamefiPlatformCli *gamefi_platform.Client
	GamefiAccountCli  *gamefi_account.Client
}

func NewSdks(config *conf.Client, logger log.Logger) (*Sdks, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the sdks")
	}

	gtsShopCli := gts_shop.NewClient(config.GetGtsShop().GetEndpoint(), config.GetGtsShop().GetTimeout().AsDuration(), conf.Prefix, logger)
	platformCli := gamefi_platform.NewClient(config.GetPlatformDust().GetEndpoint(), config.GetPlatformDust().GetTimeout().AsDuration(), conf.Prefix, logger)
	accountCli := gamefi_account.NewClient(config.GetGamefiAccount().GetEndpoint(), config.GetGamefiAccount().GetTimeout().AsDuration(), conf.Prefix, logger)

	return &Sdks{
		GtsCenter:         gtscenter.NewGtsCenter(gtsShopCli, platformCli, logger),
		GtsShopCli:        gtsShopCli,
		GamefiPlatformCli: platformCli,
		GamefiAccountCli:  accountCli,
	}, cleanup, nil
}
