package gts_shop

import (
	"components/sdks/gts_shop/pb/vo"
	"context"
)

// 登录鉴权
func (c *Client) AuthCode(ctx context.Context, code string) (accountInfo *vo.GetAccountResponse, err error) {
	accountInfo, err = c.client.AuthCode(ctx, &vo.AuthorizationCodeRequest{Code: code})
	return
}

// 发送code
func (c *Client) SendVerificationCode(ctx context.Context, realUserId string) (resp *vo.SendVerificationCodeResponse, err error) {
	resp, err = c.client.SendVerificationCode(ctx, &vo.SendVerificationCodeRequest{
		UserId: realUserId,
	})
	return
}

// 验证code
func (c *Client) VerifyCode(ctx context.Context, realUserId string, code string) (resp *vo.VerifyCodeResponse, err error) {
	//return &vo.VerifyCodeResponse{Data: true}, nil
	resp, err = c.client.VerifyCode(ctx, &vo.VerifyCodeRequest{
		UserId:     realUserId,
		GoogleCode: code,
	})
	return
}
