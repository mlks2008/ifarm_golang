package gamefi_platform

import (
	"components/sdks/gamefi_platform/pb"
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
	client            pb.PlatformDustHTTPClient
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
		c.logger.Warnf("[NewClient] gamefi_platform new http client error:%v", err)
		return nil
	}
	c.client = pb.NewPlatformDustHTTPClient(client)

	return c
}
