package nacos

import (
	unetwork "components/common/utils/network"
	clog "components/log/zaplogger"
	"errors"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

type (
	Client struct {
		NamingClient naming_client.INamingClient
		ConfigClient config_client.IConfigClient
		serviceName  string
		clusterName  string
	}
)

func NewClient(nacosIp string, nacosPort int64, userName, pwd, namespaceId string, timeoutMs int64, clusterName string, serviceName string, logDir string, debug bool) *Client {
	client := &Client{}
	client.serviceName = serviceName
	client.clusterName = clusterName

	params := nacosClientParam(nacosIp, nacosPort, userName, pwd, namespaceId, timeoutMs, serviceName, logDir, debug)

	configClient, err := clients.NewConfigClient(params)
	if err != nil {
		panic(err)
	}
	client.ConfigClient = configClient

	namingClient, err := clients.NewNamingClient(params)
	if err != nil {
		panic(err)
	}
	client.NamingClient = namingClient

	return client
}

func nacosClientParam(nacosIp string, nacosPort int64, userName, pwd, namespaceId string, timeoutMs int64, serviceName string, logDir string, debug bool) vo.NacosClientParam {
	cc := &constant.ClientConfig{
		Username:            userName,
		Password:            pwd,
		NamespaceId:         namespaceId,
		TimeoutMs:           uint64(timeoutMs),
		LogLevel:            "warn",
		LogDir:              clog.GetFileLinkPath(logDir, serviceName, "nacos"),
		CacheDir:            clog.GetFileLinkPath(logDir, serviceName, "nacos-cache"),
		NotLoadCacheAtStart: true,
	}

	sc := []constant.ServerConfig{
		{
			IpAddr: nacosIp,
			Port:   uint64(nacosPort),
		},
	}

	if debug {
		cc.AppendToStdout = true
		cc.LogLevel = "debug"
	}

	return vo.NacosClientParam{
		ClientConfig:  cc,
		ServerConfigs: sc,
	}
}

func (c *Client) Register(serverHttpPort int64) (bool, error) {
	//metadata := make(map[string]string)
	//if Ws != nil {
	//	metadata[GameServerWsPort] = cast.ToString(Ws.Port)
	//}

	return c.NamingClient.RegisterInstance(vo.RegisterInstanceParam{
		ServiceName: c.serviceName,
		Ip:          unetwork.GetIntranetIpByPrev(""),
		Port:        uint64(serverHttpPort),
		ClusterName: c.clusterName,
		Weight:      1,
		//Metadata:    metadata,
		Enable:    true,
		Healthy:   true,
		Ephemeral: true,
	})
}

func (c *Client) Stop() {
	if c.NamingClient != nil {
		c.NamingClient.CloseClient()
	}
	if c.ConfigClient != nil {
		c.ConfigClient.CloseClient()
	}
}

func (c *Client) SelectOneNodeInfoBySvcName(serviceName string) (nodeInfo *RemoteNodeInfo, err error) {
	var targetInstance *model.Instance

	clusters := make([]string, 0)
	if c.clusterName != "" {
		clusters = append(clusters, c.clusterName)
	}

	targetInstance, err = c.NamingClient.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		ServiceName: serviceName,
		Clusters:    clusters,
	})
	if err != nil {
		return nil, err
	}
	if targetInstance == nil {
		return nil, errors.New(fmt.Sprintf("[SelectOneNodeInfoBySvcName] service [%s] not found one healthy instance! ", serviceName))
	}
	nodeInfo = &RemoteNodeInfo{}
	nodeInfo.Ip = targetInstance.Ip
	nodeInfo.HttpPort = targetInstance.Port
	//if targetInstance.Metadata != nil {
	//	if wsPost, found := targetInstance.Metadata[GameServerWsPort]; found {
	//		nodeInfo.WsPort = cast.ToUint64(wsPost)
	//	}
	//}

	return
}
