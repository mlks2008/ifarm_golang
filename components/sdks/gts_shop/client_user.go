package gts_shop

import (
	"components/sdks/gts_shop/pb/vo"
	"context"
)

func (c *Client) UserProfile(ctx context.Context, tokenUserId string) (accountInfo *vo.GetAccountResponse, err error) {
	accountInfo, err = c.client.UserProfile(ctx, &vo.UserProfileRequest{TokenUserId: tokenUserId})
	return
}
