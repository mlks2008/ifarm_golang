package gts_shop

import (
	"components/sdks/global"
	"components/sdks/gts_shop/model"
	"components/sdks/gts_shop/pb/vo"
	"context"
	"fmt"
)

// 加载资产列表
func (c *Client) ListAsset(ctx context.Context, realUserId string) ([]*model.Asset, error) {
	resp, err := c.client.ListAsset(ctx, &vo.ListAssetRequest{RealUserId: realUserId})
	if err != nil {
		return nil, err
	}
	var assets = make([]*model.Asset, 0)
	assets = append(assets, &model.Asset{AssetId: global.Asset_DUST, AssetAmount: fmt.Sprintf("%v", resp.Data.UserGalaxyDust.GetDustNum()), AssetBalance: fmt.Sprintf("%v", resp.Data.UserGalaxyDust.GetDustNum())})
	assets = append(assets, &model.Asset{AssetId: global.Asset_GP, AssetAmount: fmt.Sprintf("%v", resp.Data.UserGp.GetGpValue()), AssetBalance: fmt.Sprintf("%v", resp.Data.UserGp.GetGpValue())})
	return assets, nil
}

// 加载资产详细
func (c *Client) GetAsset(ctx context.Context, realUserId string, assetId string) (*model.Asset, error) {
	assets, err := c.ListAsset(ctx, realUserId)
	if err != nil {
		return nil, err
	}
	for _, asset := range assets {
		if asset.AssetId == assetId {
			return asset, nil
		}
	}
	return nil, nil
}
