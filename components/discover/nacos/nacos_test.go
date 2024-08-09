package nacos

import (
	"testing"
)

func TestSelectNode(t *testing.T) {
	client := NewClient(
		"127.0.0.1",
		8848,
		"nacos",
		"nacos",
		"62e5b809-a504-4185-9cb2-9d62c6dc7e70",
		5000,
		"gts-gateway",
		"gts-gateway",
		"",
		true)

	node, err := client.SelectOneNodeInfoBySvcName("gts-gateway")

	t.Logf("node:%+v, err:%v", node, err)
}
