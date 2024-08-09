package gts_shop

import (
	"components/sdks/global"
	"components/sdks/gts_shop/pb/vo"
	"context"
	"encoding/json"
	"fmt"
)

func (c *Client) Gp_AddAsset(ctx context.Context, realUserId string, assetId string, amount string, recordId int64) error {
	req := &vo.AddAssetRequest{
		UserId:      realUserId,
		Num:         amount,
		ReferenceId: global.GetReferenceId(c.referenceIdPrefix, recordId),
	}
	b, _ := json.Marshal(req)
	req.Action = fmt.Sprintf("%s", b)

	resp, err := c.client.AddAsset(ctx, req)
	if err != nil {
		return err
	}
	if fmt.Sprintf("%v", resp.Data.GetUserId()) != realUserId {
		return fmt.Errorf("Gp_AddAsset resp_user_id and req_user_id is not match")
	}
	return nil
}

// 冻结资产
func (c *Client) Gp_FreezeAsset(ctx context.Context, realUserId string, assetId string, amount string, recordId int64) error {
	req := &vo.FreezeAssetRequest{
		UserId:      realUserId,
		Num:         amount,
		ReferenceId: global.GetReferenceId(c.referenceIdPrefix, recordId),
	}
	b, _ := json.Marshal(req)
	req.Action = fmt.Sprintf("%s", b)

	resp, err := c.client.FreezeAsset(ctx, req)
	if err != nil {
		return err
	}
	if fmt.Sprintf("%v", resp.Data.GetUserId()) != realUserId {
		return fmt.Errorf("Gp_FreezeAsset resp_user_id and req_user_id is not match")
	}
	return nil
}

// 扣资产
func (c *Client) Gp_SubAsset(ctx context.Context, realUserId string, amount string, recordId int64) error {
	req := &vo.SubAssetRequest{
		UserId:      realUserId,
		Num:         amount,
		ReferenceId: global.GetReferenceId(c.referenceIdPrefix, recordId),
	}
	b, _ := json.Marshal(req)
	req.Action = fmt.Sprintf("%s", b)

	resp, err := c.client.SubAsset(ctx, req)
	if err != nil {
		return err
	}
	if fmt.Sprintf("%v", resp.Data.GetUserId()) != realUserId {
		return fmt.Errorf("Gp_SubAsset resp_user_id and req_user_id is not match")
	}
	return nil
}

func (c *Client) Gp_ReturnAsset(ctx context.Context, realUserId string, amount string, recordId int64) error {
	req := &vo.ReturnAssetRequest{
		UserId:      realUserId,
		Num:         amount,
		ReferenceId: global.GetReferenceId(c.referenceIdPrefix, recordId),
	}
	b, _ := json.Marshal(req)
	req.Action = fmt.Sprintf("%s", b)

	resp, err := c.client.ReturnAsset(ctx, req)
	if err != nil {
		return err
	}
	if fmt.Sprintf("%v", resp.Data.GetUserId()) != realUserId {
		return fmt.Errorf("Gp_Return resp_user_id and req_user_id is not match")
	}
	return nil
}
