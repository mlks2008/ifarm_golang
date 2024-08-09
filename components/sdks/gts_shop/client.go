package gts_shop

import (
	"components/sdks/global"
	"components/sdks/gts_shop/pb"
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
	client            pb.GtsShopHTTPClient
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
		c.logger.Warnf("[NewClient] gts_shop new http client error:%v", err)
		return nil
	}
	c.client = pb.NewGtsShopHTTPClient(client)

	return c
}
