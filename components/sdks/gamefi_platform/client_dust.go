package gamefi_platform

import (
	"components/sdks/gamefi_platform/pb/vo"
	"components/sdks/global"
	"context"
	"encoding/json"
	"fmt"
)

func (c *Client) Dust_AddAsset(ctx context.Context, realUserId string, assetId string, amount string, recordId int64) error {
	req := &vo.DustAddAssetRequest{
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
		return fmt.Errorf("Dust_AddAsset resp_user_id and req_user_id is not match")
	}
	return nil
}

func (c *Client) Dust_FreezeAsset(ctx context.Context, realUserId string, assetId string, amount string, recordId int64) error {
	req := &vo.DustFreezeAssetRequest{
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
		return fmt.Errorf("Dust_FreezeAsset resp_user_id and req_user_id is not match")
	}
	return nil
}

func (c *Client) Dust_SubAsset(ctx context.Context, realUserId string, recordId int64) error {
	req := &vo.DustSubAssetRequest{
		UserId:      realUserId,
		ReferenceId: global.GetReferenceId(c.referenceIdPrefix, recordId),
	}

	resp, err := c.client.SubAsset(ctx, req)
	if err != nil {
		return err
	}
	if resp.Data == false {
		return fmt.Errorf("Dust_SubAsset %+v,data:false,req:%+v", resp, req)
	}
	return nil
}

func (c *Client) Dust_ReturnAsset(ctx context.Context, realUserId string, recordId int64) error {
	req := &vo.DustReturnAssetRequest{
		UserId:      realUserId,
		ReferenceId: global.GetReferenceId(c.referenceIdPrefix, recordId),
	}

	resp, err := c.client.ReturnAsset(ctx, req)
	if err != nil {
		return err
	}
	if fmt.Sprintf("%v", resp.Data.GetUserId()) != realUserId {
		return fmt.Errorf("Dust_Return resp_user_id and req_user_id is not match")
	}
	return nil
}
