package gtsportal

import (
	"components/sdks/gtsportal/pb"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"
	"time"
)

type Client struct {
	logger *log.Helper
	client pb.GtsPortalHTTPClient
}

func NewClient(endpoint string, timeout time.Duration, logger log.Logger) *Client {
	c := &Client{}

	c.logger = log.NewHelper(logger)

	client, err := http.NewClient(context.Background(),
		http.WithEndpoint(endpoint),
		http.WithTimeout(timeout),
		http.WithResponseDecoder(DecoderResponse()))
	if err != nil {
		c.logger.Warnf("[NewClient] gts_portal new http client error:%v", err)
		return nil
	}

	c.client = pb.NewGtsPortalHTTPClient(client)

	return c
}
