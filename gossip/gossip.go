/*
 * @Description:  gossip协议实现
 */
package gossip

import (
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"path"
	"strings"

	"image/gossip/config"

	"github.com/sirupsen/logrus"
)

func LookupLocalIp() (string, error) {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		logrus.Warnln("Lookup local interface address error!")
		return "", err
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", nil
}

//
//gossip协议的实现，目前先采用轮询的方式来处理
//
type Gossiper struct {
	hostName string
	config   config.Config
}

func (g *Gossiper) Nslookup() []string {
	addrs, err := net.LookupHost(g.hostName)
	if err != nil {
		logrus.Errorln("Find hostName:", g.hostName, "err:", err)
	}
	logrus.Infoln("Fetch hostName:", g.hostName, "addres:", addrs)
	return addrs
}

func (g *Gossiper) SpreadImages(filePath string) {
	addrs := g.filterLocalIP()
	for _, v := range addrs {
		go g.sendToPeer(v, filePath)
	}
}

func (g *Gossiper) filterLocalIP() []string {
	addrs := g.Nslookup()
	ip, err := LookupLocalIp()
	if err != nil {
		return addrs
	}
	remoteAddrs := make([]string, 0)
	for _, v := range addrs {
		if strings.Compare(v, ip) != 0 {
			remoteAddrs = append(remoteAddrs, v)
		}
	}
	return remoteAddrs
}

//
// 采用Pipe的方式来处理大文件上传导致内存溢出的问题，其实就是整个文件被读进内存导致的问题
//
func (g *Gossiper) sendToPeer(addr string, filePath string) {
	url := fmt.Sprintf("http://%s:%d/v1/docker/load-image", addr, g.config.Server.Port)
	logrus.Infoln("Send image file to:", url)
	r, w := io.Pipe()
	m := multipart.NewWriter(w)

	go func() {
		defer w.Close()
		defer m.Close()
		part, err := m.CreateFormFile("file", path.Base(filePath))
		if err != nil {
			logrus.Warnln("Create fomr file part error:", filePath)
			return
		}
		file, err := os.Open(filePath)
		if err != nil {
			logrus.Warnln("Create form file error:", filePath)
			return
		}
		defer file.Close()
		if _, err = io.Copy(part, file); err != nil {
			logrus.Warnln("Copy part file error!")
			return
		}
	}()

	resp, err := http.Post(url, m.FormDataContentType(), r)
	if err != nil {
		logrus.Warnln("Send to other peer fail of:", err)
	}
	var body strings.Builder
	io.Copy(&body, resp.Body)
	logrus.Infoln("Send to peer response is:", body, resp)

}

func NewGossiper(hostName string) *Gossiper {
	return &Gossiper{
		hostName: hostName,
		config: config.Config{
			Server: config.Server{
				Port: 8080,
			},
		},
	}
}
