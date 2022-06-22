/*
 * @Description:  gossip协议实现(目前采取的方式是全部轮询操作)
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
	"sync"

	"image/gossip/config"

	"github.com/sirupsen/logrus"
)

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
	var wg sync.WaitGroup
	addrs := g.Nslookup()
	wg.Add(len(addrs))
	for _, v := range addrs {
		t := v
		go func() {
			g.sendToPeer(t, filePath)
			wg.Done()
		}()
	}
	wg.Wait()
	g.clearGossipImage(filePath)
}

func (g *Gossiper) clearGossipImage(filePath string) {
	err := os.Remove(filePath)
	if err != nil {
		logrus.Warnln("Delete gossip image fail:", filePath)
	} else {
		logrus.Infoln("Delete gossip image success:", filePath)
	}
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
