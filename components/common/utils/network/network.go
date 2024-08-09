package helper

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
)

var (
	_lock    sync.Mutex
	_localIp = ""
)

func GetIntranetIp() string {
	if _localIp != "" {
		return _localIp
	}

	_lock.Lock()
	defer _lock.Unlock()

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, address := range addrs {
		// 检查 ip 地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				_localIp = ipnet.IP.String()
				break
			}
		}
	}
	if _localIp == "" {
		_localIp = "127.0.0.1"
	}
	return _localIp
}

func GetIntranetIpByPrev(prev string) string {
	if prev == "" && _localIp != "" {
		return _localIp
	}

	_lock.Lock()
	defer _lock.Unlock()

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// key `10.0.50` 表示所在网段, value: `10.0.50.195` 表示 ip 地址
	ipArrMap := make(map[string]string)
	for _, address := range addrs {
		// 检查 ip 地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() == nil {
				continue
			}
			tmpIpStr := ipnet.IP.To4().String()
			ipArrMap[GetNetPrev(tmpIpStr)] = tmpIpStr
			if ipnet.IP.To4() != nil {
				_localIp = ipnet.IP.String()
			}
		}
	}
	if _localIp == "" {
		_localIp = "127.0.0.1"
	}
	if v, ok := ipArrMap[prev]; ok {
		return v
	}
	return _localIp
}

// GetUserAgent 从 http header 中获取代理名
func GetUserAgent(headers http.Header) string {
	for k, vs := range headers {
		if strings.ToLower(k) == "user-agent" {
			return strings.Join(vs, ",")
		}
	}
	return ""
}

// 获取 ip 所在的网段
func GetNetPrev(ipStr string) string {
	numArr := strings.Split(ipStr, ".")
	prevNumArr := numArr[0:3]
	return strings.Join(prevNumArr, ".")
}
