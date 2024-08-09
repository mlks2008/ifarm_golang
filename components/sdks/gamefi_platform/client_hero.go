package gamefi_platform

import (
	"components/sdks/gamefi_platform/pb/vo"
	"components/sdks/global"
	"context"
	"fmt"
)

func (c *Client) Hero_Add(ctx context.Context, realUserId string, itemSeq []string, itemNum []int64, recordId int64) error {
	req := &vo.AddHeroRequest{
		UserId:      realUserId,
		ReferenceId: global.GetReferenceId(c.referenceIdPrefix, recordId),
	}
	req.Data = make([]*vo.AddHeroRequest_Data, 0)
	for i, seq := range itemSeq {
		for n := 0; n < int(itemNum[i]); n++ {
			req.Data = append(req.Data, &vo.AddHeroRequest_Data{HeroId: seq, TransferCode: ""})
		}
	}

	resp, err := c.client.AddHero(ctx, req)
	if err != nil {
		return err
	}
	if resp.Data == false {
		return fmt.Errorf("%+v,data:false,req:%+v", resp, req)
	}
	return nil
}

func (c *Client) Hero_Freeze(ctx context.Context, realUserId string, userItemSeq []string, recordId int64) error {
	req := &vo.FreezeHeroRequest{
		UserId:      realUserId,
		ReferenceId: global.GetReferenceId(c.referenceIdPrefix, recordId),
	}
	req.UserHeroIds = userItemSeq

	resp, err := c.client.FreezeHero(ctx, req)
	if err != nil {
		return err
	}
	if resp.Data == false {
		return fmt.Errorf("%+v,data:false,req:%+v", resp, req)
	}
	return nil
}

func (c *Client) Hero_Sub(ctx context.Context, realUserId string, recordId int64) error {
	req := &vo.SubHeroRequest{
		UserId:      realUserId,
		ReferenceId: global.GetReferenceId(c.referenceIdPrefix, recordId),
	}

	resp, err := c.client.SubHero(ctx, req)
	if err != nil {
		return err
	}
	if resp.Data == false {
		return fmt.Errorf("%+v,data:false,req:%+v", resp, req)
	}
	return nil
}

func (c *Client) Hero_Return(ctx context.Context, realUserId string, recordId int64) error {
	req := &vo.ReturnHeroRequest{
		UserId:      realUserId,
		ReferenceId: global.GetReferenceId(c.referenceIdPrefix, recordId),
	}

	resp, err := c.client.ReturnHero(ctx, req)
	if err != nil {
		return err
	}
	if resp.Data == false {
		return fmt.Errorf("%+v,data:false,req:%+v", resp, req)
	}
	return nil
}

func (c *Client) Hero_GetId(ctx context.Context, id string) (resp *vo.GetUserHeroResponse, err error) {
	req := &vo.GetUserHeroRequest{
		Id: id,
	}
	resp, err = c.client.GetUserHero(ctx, req)

	return
}
