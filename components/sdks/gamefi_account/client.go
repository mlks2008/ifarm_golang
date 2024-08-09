package gamefi_account

import (
	"components/sdks/gamefi_account/pb"
	"components/sdks/gamefi_account/pb/vo"
	"components/sdks/global"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"
	"sync"
	"time"
)

type Client struct {
	logger            *log.Helper
	mutex             sync.Mutex
	referenceIdPrefix string
	client            pb.GamefiAccountHTTPClient
}

func NewClient(endpoint string, timeout time.Duration, referenceIdPrefix string, logger log.Logger) *Client {
	c := &Client{}
	c.referenceIdPrefix = referenceIdPrefix
	c.logger = log.NewHelper(logger)

	client, err := http.NewClient(context.Background(),
		http.WithEndpoint(endpoint),
		http.WithTimeout(timeout),
		http.WithRequestEncoder(global.EncoderRequest(c.logger)),
		http.WithResponseDecoder(global.ResponseDecoder(c.logger)))
	if err != nil {
		c.logger.Warnf("[NewClient] gamefi_account new http client error:%v", err)
		return nil
	}
	c.client = pb.NewGamefiAccountHTTPClient(client)

	return c
}

func (c *Client) SendEmail(ctx context.Context, email string, subject, content string) (bool, error) {
	c.logger.Debugf("sendemail %v %v %v", email, subject, content)
	reply, err := c.client.SendEmail(ctx, &vo.SendEmailRequest{Email: email, Subject: subject, Content: content})
	if err != nil {
		return false, err
	}
	if reply.Data == true {
		return true, nil
	} else {
		c.logger.Warnf("sendemail fail %v %v", email, subject)
		return false, nil
	}
}
