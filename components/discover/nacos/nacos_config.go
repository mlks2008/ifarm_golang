package nacos

type (
	RemoteNodeInfo struct {
		Ip       string `json:"ip"`
		HttpPort uint64 `json:"httpPort"`
		WsPort   uint64 `json:"wsPort"`
	}
)
