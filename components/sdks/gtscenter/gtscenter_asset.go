package gtscenter

import (
	"components/sdks/global"
	"context"
	"fmt"
)

func (c *GtsCenter) AddAsset(ctx context.Context, realUserId string, assetId string, amount string, recordId int64) (err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if assetId == global.Asset_DUST {
		err = c.platformDustCli.Dust_AddAsset(ctx, realUserId, assetId, amount, recordId)
	} else if assetId == global.Asset_GP {
		err = c.gtsShopCli.Gp_AddAsset(ctx, realUserId, assetId, amount, recordId)
	} else {
		err = fmt.Errorf("%v is not exist", assetId)
	}

	return
}

func (c *GtsCenter) FreezeAsset(ctx context.Context, realUserId string, assetId string, amount string, recordId int64) (err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if assetId == global.Asset_DUST {
		err = c.platformDustCli.Dust_FreezeAsset(ctx, realUserId, assetId, amount, recordId)
	} else if assetId == global.Asset_GP {
		err = c.gtsShopCli.Gp_FreezeAsset(ctx, realUserId, assetId, amount, recordId)
	} else {
		err = fmt.Errorf("%v is not exist", assetId)
	}

	return
}

func (c *GtsCenter) SubAsset(ctx context.Context, realUserId string, assetId string, amount string, recordId int64) (err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if assetId == global.Asset_DUST {
		err = c.platformDustCli.Dust_SubAsset(ctx, realUserId, recordId)
	} else if assetId == global.Asset_GP {
		err = c.gtsShopCli.Gp_SubAsset(ctx, realUserId, amount, recordId)
	} else {
		err = fmt.Errorf("%v is not exist", assetId)
	}

	return
}

func (c *GtsCenter) ReturnAsset(ctx context.Context, realUserId string, assetId string, amount string, recordId int64) (err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if assetId == global.Asset_DUST {
		err = c.platformDustCli.Dust_ReturnAsset(ctx, realUserId, recordId)
	} else if assetId == global.Asset_GP {
		err = c.gtsShopCli.Gp_ReturnAsset(ctx, realUserId, amount, recordId)
	} else {
		err = fmt.Errorf("%v is not exist", assetId)
	}

	return
}
