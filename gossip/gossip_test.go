/*
 * @Description: gossip.go的test单元测试的代码
 */
package gossip

import (
	"fmt"
	"testing"
)

func TestNslookup(t *testing.T) {
	hostName := "loki-memberlist.liuxu.svc.cluster.local"
	gossiper := NewGossiper(hostName)

	addrs := gossiper.Nslookup()

	fmt.Printf("解析结果是:%v\n", addrs)
}

func TestSendToPeer(t *testing.T) {
	hostName := "loki-memberlist.liuxu.svc.cluster.local"
	gossiper := NewGossiper(hostName)

	gossiper.sendToPeer("172.16.73.116", "/media/liuxu/data/rattrap/golang.tar")
}
